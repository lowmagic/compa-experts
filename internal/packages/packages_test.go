// Package packages_test holds the checks every expert and catalog package
// of this repository passes: what it embeds, that the engine loads it, and
// that each catalog tool returns what its fixtures say.
package packages_test

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/lowmagic/compa-engine/engine/compose"
	"github.com/lowmagic/compa-engine/engine/expertcheck"

	"github.com/lowmagic/compa-experts/catalogs/files"
	"github.com/lowmagic/compa-experts/experts/artifacts"
)

const repo = "../.."

// all is every package of this repository. A package missing from it fails
// TestEveryPackageIsListed.
var all = []fs.FS{files.FS, artifacts.FS}

// packageDirs are the directories holding one package each.
func packageDirs(t *testing.T) []string {
	t.Helper()
	var dirs []string
	for _, kind := range []string{"catalogs", "experts"} {
		entries, err := os.ReadDir(filepath.Join(repo, kind))
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			if e.IsDir() {
				dirs = append(dirs, filepath.Join(kind, e.Name()))
			}
		}
	}
	return dirs
}

func filesUnder(t *testing.T, fsys fs.FS) []string {
	t.Helper()
	var out []string
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			out = append(out, p)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(out)
	return out
}

var (
	embedLine   = regexp.MustCompile(`(?m)^//go:embed (.+)$`)
	packageLine = regexp.MustCompile(`(?m)^package (\w+)$`)
	fsVar       = regexp.MustCompile(`(?m)^var FS fs\.FS\b`)
)

// embedded expands a //go:embed directive's patterns over dir as the embed
// package does: a directory is taken whole, without names beginning with
// "." or "_" unless the pattern begins with "all:".
func embedded(t *testing.T, dir string, patterns []string) []string {
	t.Helper()
	seen := map[string]bool{}
	for _, pattern := range patterns {
		pattern, hidden := strings.CutPrefix(pattern, "all:")
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil || len(matches) == 0 {
			t.Fatalf("%s: pattern %q matches nothing (%v)", dir, pattern, err)
		}
		for _, m := range matches {
			err := filepath.WalkDir(m, func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if name := d.Name(); p != m && !hidden && (strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")) {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				if !d.IsDir() {
					rel, _ := filepath.Rel(dir, p)
					seen[filepath.ToSlash(rel)] = true
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	out := make([]string, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// TestPackageEmbedsItsDirectory: each package is one Go file exporting FS,
// named for its directory, whose directive embeds every other file of the
// directory. A file the directive misses would be on disk, pass the checks
// that read the disk, and be absent from every build.
func TestPackageEmbedsItsDirectory(t *testing.T) {
	for _, rel := range packageDirs(t) {
		t.Run(rel, func(t *testing.T) {
			dir := filepath.Join(repo, rel)
			sources, _ := filepath.Glob(filepath.Join(dir, "*.go"))
			if len(sources) != 1 {
				t.Fatalf("%d Go files, want one", len(sources))
			}
			raw, err := os.ReadFile(sources[0])
			if err != nil {
				t.Fatal(err)
			}
			src := string(raw)
			if !fsVar.MatchString(src) {
				t.Error("no top-level var FS fs.FS")
			}
			want := strings.ReplaceAll(filepath.Base(rel), "-", "")
			if m := packageLine.FindStringSubmatch(src); m == nil || m[1] != want {
				t.Errorf("package clause %v, want package %s", m, want)
			}
			directives := embedLine.FindAllStringSubmatch(src, -1)
			if len(directives) != 1 {
				t.Fatalf("%d //go:embed directives, want one", len(directives))
			}
			var onDisk []string
			for _, p := range filesUnder(t, os.DirFS(dir)) {
				if !strings.HasSuffix(p, ".go") {
					onDisk = append(onDisk, p)
				}
			}
			got := embedded(t, dir, strings.Fields(directives[0][1]))
			if strings.Join(got, "\n") != strings.Join(onDisk, "\n") {
				t.Errorf("the directive embeds:\n%s\nthe directory holds:\n%s", strings.Join(got, "\n"), strings.Join(onDisk, "\n"))
			}
		})
	}
}

// TestEveryPackageIsListed: the packages composed below are the
// directories on disk, file for file.
func TestEveryPackageIsListed(t *testing.T) {
	b, err := compose.Bundle(all...)
	if err != nil {
		t.Fatal(err)
	}
	var onDisk int
	for _, rel := range packageDirs(t) {
		for _, p := range filesUnder(t, os.DirFS(filepath.Join(repo, rel))) {
			if !strings.HasSuffix(p, ".go") {
				onDisk++
			}
		}
	}
	if got := len(filesUnder(t, b)); got != onDisk {
		t.Errorf("the composed packages hold %d files, the package directories %d: list every package in all", got, onDisk)
	}
}

// TestComposedPackagesLoad: every package of the repository, composed into
// one bundle and written out as the engine writes a bundle, loads through
// the engine's own check with no error and no warning, and every catalog
// tool's fixture cases return what they say (compa-engine's
// docs/EXPERT_API.md §Fixture cases).
func TestComposedPackagesLoad(t *testing.T) {
	b, err := compose.Bundle(all...)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.CopyFS(root, b); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := expertcheck.Command(context.Background(), []string{"--fail-on-warning", filepath.Join(root, "experts")}, &out); err != nil {
		t.Errorf("%v\n%s", err, out.String())
	}
	if strings.Contains(out.String(), "; 0 cases run") {
		t.Errorf("no fixture case ran:\n%s", out.String())
	}
}
