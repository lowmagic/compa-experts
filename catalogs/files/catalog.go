// Package files is the generic.files catalog: tools that list, search and
// sample files of any kind, composed from the engine's primitives.
package files

import (
	"embed"
	"io/fs"
)

//go:embed *.yaml all:fixtures
var content embed.FS

// FS is the package: catalog.yaml at its root.
var FS fs.FS = content
