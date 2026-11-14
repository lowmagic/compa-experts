# compa-experts

Expert and catalog packages for [compa](https://github.com/lowmagic/compa), and `compose`, which merges the packages
a build chooses into the bundle it ships.

An expert is guidance for one field, and a catalog is a set of tools experts share; both are data, read by
[compa-engine](https://github.com/lowmagic/compa-engine) (its
[Expert API](https://github.com/lowmagic/compa-engine/blob/main/docs/EXPERT_API.md) is the contract). Here each is a Go
package that embeds its own directory and exports it as `FS`.

| Package | ID | What it is |
|---|---|---|
| `catalogs/files` | `generic.files` | Tools that list, search and sample files of any kind: `ls`, `find`, `tree`, `grep`, `head`, `tail`, `count` |
| `experts/artifacts` | `generic.artifacts` | How to read a set of files nobody has described yet |

The packages here are generic: they know no product.

## Running it

A build's `main` imports the packages it ships and composes them, so that file says what the build holds and its
`go.sum` records each version:

```go
import (
	"github.com/lowmagic/compa-experts/catalogs/files"
	"github.com/lowmagic/compa-experts/compose"
	"github.com/lowmagic/compa-experts/experts/artifacts"
)

bundle, err := compose.Bundle(files.FS, artifacts.FS)
// engine.Options{Bundle: bundle, …}
```

`compose.Bundle` takes packages from any module, a routing library among them, and refuses a package without exactly
one manifest at its root, two packages with one ID, a second routing library, and an expert whose catalog is not
among the packages. It leaves version ranges and tool lists to the engine's loader, so a build keeps the call in a
function a test also calls, and that test runs the engine's check over the result.

- `make check` runs the checks CI runs, offline: the linter, the tests and `go mod tidy -diff`. The tests load every
  package through the engine's own loader and run every catalog tool against its fixtures.
- `make standards` checks the org standards copies and needs the network.
