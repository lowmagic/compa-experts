# compa-experts

Expert and catalog packages for [compa](https://github.com/lowmagic/compa).

An expert is guidance for one field, and a catalog is a set of tools experts share; both are data, read by
[compa-engine](https://github.com/lowmagic/compa-engine) (its
[Expert API](https://github.com/lowmagic/compa-engine/blob/main/docs/EXPERT_API.md) is the contract). Here each is a Go
package that embeds its own directory and exports it as `FS`.

| Package | ID | What it is |
|---|---|---|
| `catalogs/files` | `generic.files` | Tools that list, search and sample files of any kind: `ls`, `find`, `tree`, `grep`, `head`, `tail`, `count` |
| `experts/artifacts` | `generic.artifacts` | How to read a set of files nobody has described yet |
| `experts/connections` | `generic.connections` | A service's failed network connections as its own logs record them; serves the routing key `connectivity` |
| `routing` | — | The routing library: the routing keys the experts serve, compiled by the engine from `data/routing/keys.yaml` and the experts' manifests |

The packages here are generic: they know no product. A build that composes `routing` has a case whose routing key is
`connectivity` start with `generic.connections` loaded; what routing keys a case has is the build's case source's to
say.

## Running it

A build's `main` imports the packages it ships and composes them with compa-engine's `engine/compose`, so that file
says what the build holds and its `go.sum` records each version:

```go
import (
	"github.com/lowmagic/compa-engine/engine/compose"
	"github.com/lowmagic/compa-experts/catalogs/files"
	"github.com/lowmagic/compa-experts/experts/artifacts"
)

bundle, err := compose.Bundle(files.FS, artifacts.FS)
// engine.Options{Bundle: bundle, …}
```

- `make check` runs the checks CI runs, offline: the linter, the tests and `go mod tidy -diff`. The tests load every
  package through the engine's own check, which runs every catalog tool's fixture cases.
- `make standards` checks the org standards copies and needs the network.
