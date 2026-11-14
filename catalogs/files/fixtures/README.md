# Fixtures

`root/` is a small synthetic set of files, written for these tools; nothing in it comes from a real system. Each
tool's directory holds `cases.yaml`: the parameters of each run against `root/` and the output it must return, or
the status and the reason it must give. `internal/packages` runs them through the engine's own loader and runner
(`engine/expertcheck`), with every entry's modification time set to 2026-03-01T12:00:00Z, since git keeps none.
