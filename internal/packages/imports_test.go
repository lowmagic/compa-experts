package packages_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestImportsOnlyDeclaredModules: every import is the standard library,
// this module or a module go.mod requires, and go.mod replaces none, so
// the packages are what anyone fetches. A workspace spanning other
// checkouts would let an undeclared import build.
func TestImportsOnlyDeclaredModules(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(repo, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	var allowed []string
	block := false
	for _, line := range strings.Split(string(raw), "\n") {
		fields := strings.Fields(line)
		switch {
		case len(fields) == 0:
		case fields[0] == "replace":
			t.Errorf("go.mod replaces a module: %s", line)
		case fields[0] == "module":
			allowed = append(allowed, fields[1])
		case block && fields[0] == ")":
			block = false
		case block:
			allowed = append(allowed, fields[0])
		case fields[0] == "require" && len(fields) == 2 && fields[1] == "(":
			block = true
		case fields[0] == "require" && len(fields) >= 3:
			allowed = append(allowed, fields[1])
		}
	}
	fset := token.NewFileSet()
	err = filepath.WalkDir(repo, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if name := d.Name(); path != repo && (strings.HasPrefix(name, ".") || name == "testdata" || name == "fixtures") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
	imports:
		for _, spec := range file.Imports {
			imported, _ := strconv.Unquote(spec.Path.Value)
			if first, _, _ := strings.Cut(imported, "/"); !strings.Contains(first, ".") {
				continue
			}
			for _, mod := range allowed {
				if imported == mod || strings.HasPrefix(imported, mod+"/") {
					continue imports
				}
			}
			t.Errorf("%s imports %s, which go.mod does not require", path, imported)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
