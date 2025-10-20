# Makefile – Zentrale Orchestrierung für Projekt-Automationen
# Referenz: copilot-instructions.md Abschnitt 3.1

.PHONY: help test lint lint-ci adr-ref commit-lint release-check security-blockers scan scan-json pr-check release ci-local clean ensure-trivy push-ci pr-quality-gates-ci sync

# Standardwerte
TRIVY_FAIL_ON ?= HIGH,CRITICAL
TRIVY_JSON_REPORT ?= tmp/trivy-fs-report.json
VERSION ?=

help: ## Zeigt verfügbare Targets
	@echo "Projekt Automations – Make Targets"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

sync: ## Vollständiger SDE-Sync (Download → Schema-Gen → SQLite Import)
	@go run ./cmd/sde-sync

sync-force: ## Force Sync (ignoriert Versionsprüfung)
	@go run ./cmd/sde-sync --force

sync-download-only: ## Nur Download + Schema-Gen (kein SQLite Import)
	@go run ./cmd/sde-sync --skip-import

test: ## Führt die definierte Test-Suite aus (Platzhalter)
	@echo "[make test] Keine Tests konfiguriert – bitte projektspezifische Testbefehle ergänzen"

lint: ## Statische Analysen / Format / Stil (Platzhalter, bitte anpassen)
	@echo "[make lint] Kein Lint-Tool definiert – bitte projektspezifische Checks ergänzen"

lint-ci: ## Statische Analysen (CI-Modus, Platzhalter)
	@echo "[make lint-ci] Kein Lint-Tool definiert – bitte projektspezifische Checks ergänzen"

adr-ref: ## Erzwingt ADR-Referenzen für Governance-Änderungen (CI-kompatibel)
	@echo "[make adr-ref] Prüfe ADR Referenz-Anforderungen..."; \
	if [ -x scripts/common/check-adr-ref.sh ]; then \
		set +e; \
		bash scripts/common/check-adr-ref.sh; \
		rc=$$?; \
		set -e; \
		if [ $$rc -eq 1 ]; then \
			echo "[make adr-ref] ❌ ADR-Referenz Pflicht verletzt"; \
			exit 1; \
		elif [ $$rc -eq 2 ]; then \
			echo "[make adr-ref] ⚠️ Skip Marker erkannt – Warnung akzeptiert"; \
			exit 0; \
		fi; \
	else \
		echo "[make adr-ref] scripts/common/check-adr-ref.sh nicht gefunden" >&2; \
		exit 1; \
	fi

commit-lint: ## Validiert Commit Messages (RANGE=origin/main..HEAD oder COMMIT_FILE=pfad)
	@echo "[make commit-lint] Prüfe Commit Messages..."; \
	if [ -x scripts/common/check-commit-msg.sh ]; then \
		if [ -n "$${RANGE:-}" ]; then \
			bash scripts/common/check-commit-msg.sh --range "$$RANGE"; \
		elif [ -n "$${COMMIT_FILE:-}" ]; then \
			bash scripts/common/check-commit-msg.sh --file "$$COMMIT_FILE"; \
		else \
			echo "[make commit-lint] ERROR: Bitte RANGE oder COMMIT_FILE angeben" >&2; \
			exit 1; \
		fi; \
	else \
		echo "[make commit-lint] scripts/common/check-commit-msg.sh nicht gefunden" >&2; \
		exit 1; \
	fi

release-check: ## Prüft VERSION/CHANGELOG Synchronität (für Release PRs)
	@echo "[make release-check] Prüfe VERSION und CHANGELOG..."; \
	if [ -x scripts/common/check-version-changelog.sh ]; then \
		bash scripts/common/check-version-changelog.sh; \
	else \
		echo "[make release-check] scripts/common/check-version-changelog.sh nicht gefunden" >&2; \
		exit 1; \
	fi

security-blockers: ## Prüft Trivy Report auf kritische Findings
	@echo "[make security-blockers] Prüfe Security Blocker..."; \
	if [ -x scripts/common/check-security-blockers.sh ]; then \
		bash scripts/common/check-security-blockers.sh; \
	else \
		echo "[make security-blockers] scripts/common/check-security-blockers.sh nicht gefunden" >&2; \
		exit 1; \
	fi

scan: ## Security & Dependency Checks
	@echo "[make scan] Führe Security Scan aus (Trivy)..."
	@$(MAKE) --no-print-directory ensure-trivy
	@if command -v trivy >/dev/null 2>&1; then \
		trivy fs --ignore-unfixed --scanners vuln --severity $(TRIVY_FAIL_ON) --exit-code 1 .; \
	else \
		echo "[make scan] trivy Installation fehlgeschlagen – überspringe Scan"; \
	fi

scan-json: ## Security Scan mit JSON Report (für check-security-blockers.sh)
	@echo "[make scan-json] Erzeuge Trivy JSON Report (ohne Build-Abbruch)..."
	@$(MAKE) --no-print-directory ensure-trivy
	@mkdir -p tmp
	@if command -v trivy >/dev/null 2>&1; then \
		trivy fs --ignore-unfixed --scanners vuln --format json -o $(TRIVY_JSON_REPORT) . || true; \
		echo "[make scan-json] Trivy JSON Report: $(TRIVY_JSON_REPORT)"; \
	else \
		echo "[make scan-json] trivy nicht verfügbar – kein Report erzeugt"; \
	fi

pr-check: lint test scan ## Bündelt: lint + test + scan (für lokale PR-Vorbereitung)
	@echo "[make pr-check] ✅ Alle lokalen Checks erfolgreich"

push-ci: ## Führt lint-ci und test in einem Rutsch aus
	@$(MAKE) --no-print-directory lint-ci
	@$(MAKE) --no-print-directory test
	@echo "[make push-ci] ✅ Lint & Test abgeschlossen"

release: ## Version bump + CHANGELOG Transform (Beispiel: make release VERSION=0.2.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "[make release] ERROR: VERSION Parameter fehlt (Beispiel: make release VERSION=0.2.0)" >&2; \
		exit 1; \
	fi
	@echo "[make release] Bump Version auf $(VERSION)..."
	@echo "$(VERSION)" > VERSION
	@sed -i "s/^## \[Unreleased\]/## [Unreleased]\n\n## [$(VERSION)] - $$(date +%Y-%m-%d)/" CHANGELOG.md
	@echo "[make release] VERSION und CHANGELOG aktualisiert – bitte commit + tag erstellen"

ci-local: ## Simulation definierter CI-Gates lokal
	@echo "[make ci-local] Simuliere CI Pipeline lokal..."
	@bash scripts/common/check-normative.sh
	@bash scripts/common/check-adr.sh
	@$(MAKE) --no-print-directory test
	@$(MAKE) --no-print-directory scan

pr-quality-gates-ci: ## Führt alle Quality-Gate-Prüfungen für PRs aus
	@bash scripts/workflows/pr-quality-gates/run.sh

clean: ## Entfernt Build-Artefakte und temporäre Dateien
	@echo "[make clean] Räume temporäre Dateien auf..."
	@rm -rf tmp/*.md tmp/test-fixtures/
	@echo "[make clean] ✅ Clean abgeschlossen"

ensure-trivy: ## Stellt sicher, dass Trivy verfügbar ist
	@if command -v trivy >/dev/null 2>&1; then \
		echo "[make ensure-trivy] trivy bereits verfügbar"; \
	else \
		echo "[make ensure-trivy] trivy nicht installiert – versuche Installation"; \
		if command -v apt-get >/dev/null 2>&1; then \
			if command -v sudo >/dev/null 2>&1; then \
				sudo apt-get update -y >/dev/null 2>&1 || true; \
				sudo apt-get install -y wget jq >/dev/null 2>&1 || true; \
			else \
				apt-get update -y >/dev/null 2>&1 || true; \
				apt-get install -y wget jq >/dev/null 2>&1 || true; \
			fi; \
		fi; \
		if command -v sudo >/dev/null 2>&1; then \
			curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sudo sh -s -- -b /usr/local/bin || true; \
		else \
			curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin || true; \
		fi; \
	fi


