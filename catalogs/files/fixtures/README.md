# Fixtures

`root/` is a small synthetic set of files, written for these tools; nothing in it comes from a real system. Each
tool's directory holds `cases.yaml`, in the format of compa-engine's Expert API (§Fixture cases): the parameters of
each run against `root/`, with every entry's modification time set to 2026-03-01T12:00:00Z since git keeps none, and
what it must return. The engine's own check (`engine/expertcheck`) runs them.
