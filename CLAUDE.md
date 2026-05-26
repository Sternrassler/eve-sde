# eve-sde

SQLite-DB-Generator für EVE Static Data Export (Go 1.25). Hauptprodukt:
`eve-sde.db` (~405 MB) unter `data/`, wird von `eve-o-provit` via `pkg/evedb`
eingebunden. Package-Struktur & Schema-Typen: [→ ../docs/eve-sde.md](../docs/eve-sde.md).

## Commands

`make help` listet alle Targets. Wichtigste:

- `make sync` — vollständiger SDE-Sync (Download → Schema-Gen → SQLite-Import), läuft `go run ./cmd/sde-sync`.
- `make sync-force` — ignoriert die Versionsprüfung.
- `make sync-download-only` — nur Download + Schema-Gen, kein SQLite-Import.
- `make test` / `make lint` — Test- und Lint-Targets (z. T. Platzhalter, vor Nutzung prüfen).
- `make pr-check` — lokales PR-Gate: `lint + test + scan`.

## Gotchas

- `make scan` braucht **Trivy** (`make ensure-trivy`).
- Conventional Commits erzwungen via `.githooks` (`git config core.hooksPath .githooks`).
- Release: SemVer lebt nur in `CHANGELOG.md` + git-Tag (`vX.Y.Z`) — keine `VERSION`-Datei. `make release-check` prüft die CHANGELOG-Konsistenz; `make release VERSION=x.y.z` transformiert den `[Unreleased]`-Block.
- CI-Schedule: täglich 03:00 UTC — prüft neue SDE-Version, generiert `eve-sde.db`, veröffentlicht als Release.
