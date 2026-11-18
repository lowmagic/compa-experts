# Coding standards

compa-experts follows the [org standards](org-standards.md) and the standards for its topics
([go](org-standards-go.md)), unchanged copies of
[lowmagic's](https://github.com/lowmagic/.github/tree/main/code-standards). `.org-standards` lists the topics it copies
beyond the general one. This file adds compa-experts's own rules and examples. Where it differs from an org rule, it
says so, and it wins.

## Sources

The rule no change may weaken: every package here is written fresh, from public sources only.

- Public sources are upstream documentation and code, public product documentation, and compa-engine's own
  documents. A package's guidance cites only those.
- Nothing is copied, paraphrased or restructured from a private expert set, an internal playbook or a support case.
  Knowing that a private package exists is no reason to write its public twin.
- Fixtures are synthetic: written for their tool, from no real system, and they say so.

## Structure

- One Go package per expert, under `experts/`, and per catalog, under `catalogs/`. The package is one Go file that
  embeds every other file of its directory with one `//go:embed` directive and exports it as `FS`, with the manifest,
  `expert.yaml` or `catalog.yaml`, at the root. It exports nothing else: a package's ID and version are its
  manifest's.
- The Go package name is the directory name without hyphens. The ID is the manifest's, and never begins with a
  prefix a private expert set uses.
- Dependencies run one way. An expert here uses catalogs here. Nothing here names a package of another repository.
- A package's content is checked by the engine, never by a copy of its rules: the tests call
  `github.com/lowmagic/compa-engine/engine/expertcheck`.

## Versions and pins

- A package's version is `metadata.version` in its manifest, semver over what an expert or a model relies on: a
  tool's ID, its parameters and the shape of its output; an expert's ID and the tools it lists. A change to any
  file of a package takes a new version.
- `go.mod` requires compa-engine at the commit the packages were checked against.

## Testing

`internal/packages` holds the checks every package passes: its directive embeds its whole directory; and every package,
composed, passes the engine's own check with no warning, which runs each catalog tool's `fixtures/<tool>/cases.yaml`
against the catalog's `fixtures/root` and compares what each run returns with what the file says, in the format
compa-engine's [Expert API](https://github.com/lowmagic/compa-engine/blob/main/docs/EXPERT_API.md#fixture-cases)
defines. A new tool comes with its cases, and a case that shows the tool refusing what it must refuse;
`compa expert check --update` writes their results, and the diff is read before it is committed.

## Go

`go.mod`'s `go` line stays on the oldest Go release line a compa build is made with, so every build can import these
packages. This overrides the Layout rule that it names the current release.

## Known deviations

Where compa-experts doesn't follow the org standards yet. These are debt to pay down, not overrides: each one should go.

- There's no `docs/SPEC.md` or `docs/ROADMAP.md` yet (see [Docs](org-standards.md#docs)).
- Nothing checks that a changed package took a new version.
