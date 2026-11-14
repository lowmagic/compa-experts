// Package artifacts is the generic.artifacts expert: how to read a set of
// files nobody has described yet.
package artifacts

import (
	"embed"
	"io/fs"
)

//go:embed *.yaml *.md
var content embed.FS

// FS is the package: expert.yaml at its root.
var FS fs.FS = content
