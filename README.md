# compa-experts

Expert and catalog packages for compa, each a Go package exporting an fs.FS, and compose, which merges the chosen ones into the bundle a compa build ships.

## Running it

- `go run .` runs it.
- `make check` runs the checks CI runs, offline. `make standards` checks the org standards copies and needs the network.
