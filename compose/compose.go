// Package compose merges expert, catalog and routing-library packages into
// the one file system a compa build passes to the engine as its bundle.
//
// A package is an fs.FS with its manifest at the root: expert.yaml,
// catalog.yaml or library.yaml. A build's main imports the packages it
// ships and hands their file systems to Bundle, so that file says what the
// build holds.
package compose

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"slices"
	"sort"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

// The manifest at a package's root says which kind it is.
const (
	expertManifest  = "expert.yaml"
	catalogManifest = "catalog.yaml"
	libraryManifest = "library.yaml"
)

// Where the engine reads a bundle's packages and its routing library.
const (
	expertsDir = "experts"
	routingDir = "routing"
)

// manifest is the part of a package's manifest Bundle reads. Everything
// else is the engine's loader's to validate.
type manifest struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		ID string `yaml:"id"`
	} `yaml:"metadata"`
	Spec struct {
		Catalogs []struct {
			ID string `yaml:"id"`
		} `yaml:"catalogs"`
	} `yaml:"spec"`
}

// Bundle merges the packages into the file system engine.Options.Bundle
// takes: each expert and catalog under experts/<its ID>/, and the routing
// library's files under routing/.
//
// It refuses a package without exactly one manifest at its root, two
// packages with one ID, a second routing library, an expert that declares a
// catalog no package supplies, and anything in a package that is not a
// regular file or a directory. It does not check version ranges or tool
// lists, which the engine's loader evaluates, and it passes a routing
// library through unchanged, entries for experts that were not chosen
// included: the engine ignores those and says so. With no routing library
// the bundle has no routing/, and the engine routes nothing.
func Bundle(packages ...fs.FS) (fs.FS, error) {
	b := &bundle{files: map[string][]byte{}}
	kinds := map[string]string{} // package ID to "expert" or "catalog"
	needs := map[string][]string{}
	library := false
	for i, pkg := range packages {
		name, err := manifestName(pkg)
		if err != nil {
			return nil, fmt.Errorf("compose: package %d: %w", i+1, err)
		}
		if name == libraryManifest {
			if library {
				return nil, fmt.Errorf("compose: package %d is a second routing library; a bundle holds one", i+1)
			}
			library = true
			if err := b.add(pkg, routingDir); err != nil {
				return nil, fmt.Errorf("compose: the routing library: %w", err)
			}
			continue
		}
		m, err := readManifest(pkg, name)
		if err != nil {
			return nil, fmt.Errorf("compose: package %d: %w", i+1, err)
		}
		kind := "expert"
		if name == catalogManifest {
			kind = "catalog"
		}
		id := m.Metadata.ID
		if prev, ok := kinds[id]; ok {
			return nil, fmt.Errorf("compose: %s %s is given twice, the first time as a %s; package IDs are unique across experts and catalogs", kind, id, prev)
		}
		kinds[id] = kind
		for _, c := range m.Spec.Catalogs {
			needs[id] = append(needs[id], c.ID)
		}
		if err := b.add(pkg, path.Join(expertsDir, id)); err != nil {
			return nil, fmt.Errorf("compose: %s %s: %w", kind, id, err)
		}
	}
	var missing []string
	for id, catalogs := range needs {
		for _, c := range catalogs {
			if kinds[c] != "catalog" {
				missing = append(missing, fmt.Sprintf("expert %s declares catalog %s, which no package supplies", id, c))
			}
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return nil, fmt.Errorf("compose: %s", strings.Join(missing, "; "))
	}
	b.index()
	return b, nil
}

// manifestName is the one manifest at the package's root.
func manifestName(pkg fs.FS) (string, error) {
	var found []string
	for _, name := range []string{expertManifest, catalogManifest, libraryManifest} {
		info, err := fs.Stat(pkg, name)
		switch {
		case errors.Is(err, fs.ErrNotExist):
		case err != nil:
			return "", err
		case info.Mode().IsRegular():
			found = append(found, name)
		}
	}
	switch len(found) {
	case 0:
		return "", fmt.Errorf("no %s, %s or %s at its root", expertManifest, catalogManifest, libraryManifest)
	case 1:
		return found[0], nil
	}
	return "", fmt.Errorf("%s at its root; a package has one manifest", strings.Join(found, " and "))
}

func readManifest(pkg fs.FS, name string) (manifest, error) {
	var m manifest
	raw, err := fs.ReadFile(pkg, name)
	if err != nil {
		return m, err
	}
	if err := yaml.Unmarshal(raw, &m); err != nil {
		return m, fmt.Errorf("%s: %w", name, err)
	}
	want := "Expert"
	if name == catalogManifest {
		want = "Catalog"
	}
	if m.Kind != want {
		return m, fmt.Errorf("%s: kind is %q, want %s", name, m.Kind, want)
	}
	// The ID names a directory of the bundle, so it must be one path
	// element; the loader checks its form.
	if id := m.Metadata.ID; id == "" || id == "." || id == ".." || strings.ContainsAny(id, `/\`) {
		return m, fmt.Errorf("%s: metadata.id %q cannot name a directory", name, id)
	}
	return m, nil
}

// bundle is a read-only file system held in memory: a few megabytes of
// manifests, guidance and fixtures, copied once when a build starts.
type bundle struct {
	files map[string][]byte
	dirs  map[string][]string // directory to its entries' base names, sorted
}

// add copies every file of pkg under dir.
func (b *bundle) add(pkg fs.FS, dir string) error {
	return fs.WalkDir(pkg, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("%s: a package holds only files and directories", p)
		}
		raw, err := fs.ReadFile(pkg, p)
		if err != nil {
			return err
		}
		b.files[path.Join(dir, p)] = raw
		return nil
	})
}

// index derives the directories from the files' paths.
func (b *bundle) index() {
	b.dirs = map[string][]string{".": nil}
	for p := range b.files {
		for child := p; child != "."; child = path.Dir(child) {
			parent := path.Dir(child)
			if !slices.Contains(b.dirs[parent], path.Base(child)) {
				b.dirs[parent] = append(b.dirs[parent], path.Base(child))
			}
			if _, ok := b.dirs[parent]; !ok {
				b.dirs[parent] = nil
			}
		}
	}
	for dir := range b.dirs {
		sort.Strings(b.dirs[dir])
	}
}

// Open implements fs.FS.
func (b *bundle) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	if raw, ok := b.files[name]; ok {
		return &file{entry: entry{name: path.Base(name), size: int64(len(raw))}, r: strings.NewReader(string(raw))}, nil
	}
	names, ok := b.dirs[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	entries := make([]fs.DirEntry, 0, len(names))
	for _, base := range names {
		child := path.Join(name, base)
		if raw, ok := b.files[child]; ok {
			entries = append(entries, entry{name: base, size: int64(len(raw))})
		} else {
			entries = append(entries, entry{name: base, dir: true})
		}
	}
	return &dir{entry: entry{name: path.Base(name), dir: true}, entries: entries}, nil
}

// entry is a file's or a directory's fs.FileInfo and fs.DirEntry.
type entry struct {
	name string
	size int64
	dir  bool
}

func (e entry) Name() string { return e.name }
func (e entry) Size() int64  { return e.size }
func (e entry) IsDir() bool  { return e.dir }
func (entry) Sys() any       { return nil }

// ModTime is the zero time: a bundle is content, with no times of its own.
func (entry) ModTime() time.Time { return time.Time{} }

func (e entry) Mode() fs.FileMode {
	if e.dir {
		return fs.ModeDir | 0o555
	}
	return 0o444
}
func (e entry) Type() fs.FileMode          { return e.Mode().Type() }
func (e entry) Info() (fs.FileInfo, error) { return e, nil }

type file struct {
	entry
	r *strings.Reader
}

func (f *file) Stat() (fs.FileInfo, error) { return f.entry, nil }
func (f *file) Read(p []byte) (int, error) { return f.r.Read(p) }
func (*file) Close() error                 { return nil }

type dir struct {
	entry
	entries []fs.DirEntry
	next    int
}

func (d *dir) Stat() (fs.FileInfo, error) { return d.entry, nil }
func (*dir) Close() error                 { return nil }

func (d *dir) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.name, Err: errors.New("is a directory")}
}

// ReadDir implements fs.ReadDirFile.
func (d *dir) ReadDir(n int) ([]fs.DirEntry, error) {
	rest := d.entries[d.next:]
	if n <= 0 {
		d.next = len(d.entries)
		return rest, nil
	}
	if len(rest) == 0 {
		return nil, io.EOF
	}
	if n > len(rest) {
		n = len(rest)
	}
	d.next += n
	return rest[:n], nil
}
