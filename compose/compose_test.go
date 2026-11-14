package compose_test

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/lowmagic/compa-experts/compose"
)

func file(s string) *fstest.MapFile { return &fstest.MapFile{Data: []byte(s)} }

func catalog(id string) fstest.MapFS {
	return fstest.MapFS{
		"catalog.yaml":         file("kind: Catalog\nmetadata: {id: " + id + "}\n"),
		"tests.yaml":           file("kind: Tests\n"),
		"fixtures/ls/expected": file("a\n"),
	}
}

func expert(id string, catalogs ...string) fstest.MapFS {
	spec := ""
	for _, c := range catalogs {
		spec += "    - {id: " + c + "}\n"
	}
	if spec != "" {
		spec = "spec:\n  catalogs:\n" + spec
	}
	return fstest.MapFS{
		"expert.yaml": file("kind: Expert\nmetadata: {id: " + id + "}\n" + spec),
		"analyst.md":  file("# " + id + "\n"),
	}
}

func library() fstest.MapFS {
	return fstest.MapFS{
		"library.yaml": file("kind: RoutingLibrary\nkeys: []\n"),
		"terms.yaml":   file("kind: TermDictionary\n"),
	}
}

// TestBundleLayout: the bundle is laid out as the engine reads one, and
// reads as a file system should.
func TestBundleLayout(t *testing.T) {
	b, err := compose.Bundle(library(), catalog("generic.files"), expert("generic.artifacts", "generic.files"))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"experts/generic.artifacts/analyst.md",
		"experts/generic.artifacts/expert.yaml",
		"experts/generic.files/catalog.yaml",
		"experts/generic.files/fixtures/ls/expected",
		"experts/generic.files/tests.yaml",
		"routing/library.yaml",
		"routing/terms.yaml",
	}
	if err := fstest.TestFS(b, want...); err != nil {
		t.Fatal(err)
	}
	var got []string
	err = fs.WalkDir(b, ".", func(p string, d fs.DirEntry, err error) error {
		if err == nil && d.Type().IsRegular() {
			got = append(got, p)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("files:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	raw, err := fs.ReadFile(b, "routing/library.yaml")
	if err != nil || string(raw) != "kind: RoutingLibrary\nkeys: []\n" {
		t.Errorf("the library was not passed through unchanged: %q, %v", raw, err)
	}
}

// TestBundleWithoutLibrary: a build with no routing library has no
// routing directory, which is how the engine knows to route nothing.
func TestBundleWithoutLibrary(t *testing.T) {
	b, err := compose.Bundle(catalog("generic.files"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fs.Stat(b, "routing/library.yaml"); err == nil {
		t.Error("a bundle composed without a library holds routing/library.yaml")
	}
	if _, err := fs.Stat(b, "experts"); err != nil {
		t.Errorf("experts: %v", err)
	}
}

// TestBundleOfNothing: no packages is an empty bundle, not an error.
func TestBundleOfNothing(t *testing.T) {
	b, err := compose.Bundle()
	if err != nil {
		t.Fatal(err)
	}
	entries, err := fs.ReadDir(b, ".")
	if err != nil || len(entries) != 0 {
		t.Errorf("entries %v, err %v; want none", entries, err)
	}
}

func TestBundleRefuses(t *testing.T) {
	link := fstest.MapFS{
		"expert.yaml": file("kind: Expert\nmetadata: {id: generic.linked}\n"),
		"alias":       &fstest.MapFile{Mode: fs.ModeSymlink, Data: []byte("expert.yaml")},
	}
	for _, tc := range []struct {
		name     string
		packages []fs.FS
		want     string
	}{
		{"an expert without its catalog", []fs.FS{expert("generic.artifacts", "generic.files")}, "expert generic.artifacts declares catalog generic.files, which no package supplies"},
		{"a declared catalog that is an expert", []fs.FS{expert("a", "b"), expert("b")}, "expert a declares catalog b"},
		{"one ID twice", []fs.FS{catalog("generic.files"), catalog("generic.files")}, "catalog generic.files is given twice"},
		{"one ID in both kinds", []fs.FS{catalog("generic.files"), expert("generic.files")}, "expert generic.files is given twice, the first time as a catalog"},
		{"no manifest", []fs.FS{fstest.MapFS{"analyst.md": file("x")}}, "package 1: no expert.yaml, catalog.yaml or library.yaml at its root"},
		{"two manifests", []fs.FS{fstest.MapFS{"expert.yaml": file(""), "catalog.yaml": file("")}}, "package 1: expert.yaml and catalog.yaml at its root"},
		{"two libraries", []fs.FS{library(), library()}, "package 2 is a second routing library"},
		{"a manifest of the other kind", []fs.FS{fstest.MapFS{"expert.yaml": file("kind: Catalog\nmetadata: {id: x}\n")}}, `kind is "Catalog", want Expert`},
		{"an ID that is a path", []fs.FS{fstest.MapFS{"expert.yaml": file("kind: Expert\nmetadata: {id: ../x}\n")}}, "cannot name a directory"},
		{"no ID", []fs.FS{fstest.MapFS{"catalog.yaml": file("kind: Catalog\n")}}, "cannot name a directory"},
		{"a symlink", []fs.FS{link}, "alias: a package holds only files and directories"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := compose.Bundle(tc.packages...)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("err %v, want one holding %q", err, tc.want)
			}
		})
	}
}

// TestBundleKeepsUnchosenLibraryEntries: a library naming experts the build
// leaves out composes; the engine ignores those entries.
func TestBundleKeepsUnchosenLibraryEntries(t *testing.T) {
	lib := fstest.MapFS{"library.yaml": file("kind: RoutingLibrary\nkeys: [{id: k, experts: [{id: not.chosen}]}]\n")}
	if _, err := compose.Bundle(lib, catalog("generic.files")); err != nil {
		t.Errorf("a library naming an expert that was not chosen: %v", err)
	}
}
