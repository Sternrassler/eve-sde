# eve-sde

SQLite-DB-Generator für EVE Static Data Export (Go 1.25). Hauptprodukt:
`eve-sde.db` (~474 MB) unter `data/sqlite/`, wird von `eve-o-provit` via
`backend/pkg/evedb` eingebunden. Package-Struktur & Schema-Typen: [→ ../docs/eve-sde.md](../docs/eve-sde.md).

## Aufbau

- Binaries: `cmd/{sde-schema-gen,sde-sync,sde-to-sqlite,sde-version-check}`; Kernlogik in `internal/{schema,registry,sde,sqlite}`.
- `internal/registry/mappings_gen.go` ist **GENERIERT** (durch `sde-schema-gen`) — nicht von Hand editieren. `registry.go`/`overrides.go` werden manuell gepflegt.

## Commands

`make help` listet alle Targets. Wichtigste:

- `make sync` — vollständiger SDE-Sync (Download → Schema-Gen → SQLite-Import), läuft `go run ./cmd/sde-sync`.
- `make sync-force` — ignoriert die Versionsprüfung.
- `make sync-download-only` — nur Download + Schema-Gen, kein SQLite-Import.
- `make pr-check` — lokales PR-Gate: `lint + test + scan`.
- `make ci-local` — spiegelt die CI-Gates lokal (test + scan); `make pr-quality-gates-ci` ist das volle PR-Gate (`scripts/workflows/pr-quality-gates/run.sh`: commit-lint + `make scan-json` + `check-security-blockers.sh`; Release-Label-Check nur bei `CHECK_RELEASE_LABEL=1`).

## Lint & Test

- `make test`, `make lint` und `make lint-ci` sind aktuell **echo-only Platzhalter** ("Keine Tests/kein Lint-Tool konfiguriert") — kein echter Test-/Lint-Lauf. Folge: CI (`lint-test.yml` → `make push-ci`) fährt **keine** echten Tests.
- Es existieren dennoch reale Go-Tests (7 `_test.go` in `internal/sqlite/importer`, `internal/sqlite/schema`, `internal/sde/version`) — manuell via `go test ./...` ausführen, nicht über `make test`.
- Das **echte Linting** läuft über golangci-lint v2 (`.golangci.yml`, 7 Linter: govet, errcheck, staticcheck, unused, ineffassign, gocritic, misspell) — ausgelöst durch den `.githooks/pre-commit`-Hook (`golangci-lint run --new-from-rev=HEAD~1`), nicht durch einen Make-Target.

## Gotchas

- `make scan` braucht **Trivy** (`make ensure-trivy`).
- Conventional Commits erzwungen via `.githooks` (`git config core.hooksPath .githooks`).
- Release: SemVer lebt nur in `CHANGELOG.md` + git-Tag (`vX.Y.Z`) — keine `VERSION`-Datei. `make release-check` prüft die CHANGELOG-Konsistenz; `make release VERSION=x.y.z` transformiert den `[Unreleased]`-Block.
- CI-Schedule: täglich 03:00 UTC — prüft neue SDE-Version, generiert `eve-sde.db`, veröffentlicht als Release.
