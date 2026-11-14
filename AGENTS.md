# Working on compa-experts

Expert and catalog packages for compa, each a Go package exporting an fs.FS, and compose, which merges the chosen ones into the bundle a compa build ships.

Read before changing behavior:

- [README](README.md): what compa-experts does and how to run it
- [Coding standards](docs/standards.md): how the code should read, on top of the org's
  [principles](https://github.com/lowmagic/.github/blob/main/PRINCIPLES.md) and standards (`docs/org-standards*.md`)

Docs live in `docs/` and change with behavior. Never commit hostnames, tokens, or deployment config (see
Privacy in the [org standards](docs/org-standards.md)).

Run `make check` before committing. It runs everything CI runs, except the org standards check:
`make standards` (needs the network).
