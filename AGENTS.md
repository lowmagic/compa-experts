# Working on compa-experts

Expert and catalog packages for compa, each a Go package exporting an `fs.FS`, which a compa build merges into the
bundle it ships with compa-engine's `engine/compose`.

Read before changing behavior:

- [README](README.md): what compa-experts holds and how a build uses it
- [Coding standards](docs/standards.md): how packages are written and checked, the Sources rule first, on top of the
  org's [principles](https://github.com/lowmagic/.github/blob/main/PRINCIPLES.md) and standards
  (`docs/org-standards*.md`)
- compa-engine's [Expert API](https://github.com/lowmagic/compa-engine/blob/main/docs/EXPERT_API.md) and
  [Collection](https://github.com/lowmagic/compa-engine/blob/main/docs/COLLECTION.md) contracts: what a package may
  declare and how its tools run

Every package is written from public sources only; read [Sources](docs/standards.md#sources) before writing one.

Docs live in `docs/` and change with behavior. Never commit hostnames, tokens, or deployment config (see
Privacy in the [org standards](docs/org-standards.md)).

Run `make check` before committing. It runs everything CI runs, except the org standards check:
`make standards` (needs the network).
