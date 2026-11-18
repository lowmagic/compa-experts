// Command version-check is compa-engine's compa expert versions over this
// repository's package roots (docs/standards.md §Versions and pins: a
// change to any file of a package takes a new version):
//
//	go run ./cmd/version-check                # against the merge base with origin/main
//	go run ./cmd/version-check -base <rev>    # against a revision
//
// Exit status 1 on a package changed with its version kept; 2 on a usage
// error.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/lowmagic/compa-engine/engine/expertcheck"
)

func main() {
	err := expertcheck.Versions(context.Background(), append(os.Args[1:], "experts", "catalogs"), os.Stdout)
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "version-check: %v\n", err)
	if errors.Is(err, expertcheck.ErrUsage) {
		os.Exit(2)
	}
	os.Exit(1)
}
