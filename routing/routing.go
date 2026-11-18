// Package routing is this repository's routing library: the routing keys
// its experts serve, compiled by compa-engine's compa expert routing from
// data/routing/keys.yaml and the experts' spec.routing.keys.
package routing

import (
	"embed"
	"io/fs"
)

//go:embed library.yaml
var content embed.FS

// FS is the package: library.yaml at its root.
var FS fs.FS = content
