// Package connections is the generic.connections expert: a service's failed
// network connections as its own logs record them.
package connections

import (
	"embed"
	"io/fs"
)

//go:embed *.yaml *.md all:fixtures
var content embed.FS

// FS is the package: expert.yaml at its root.
var FS fs.FS = content
