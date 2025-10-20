# Copilot Instruction – Generische Engineering Richtlinien

Diese Regeldatei gilt für technische Entwicklungs- und Automatisierungsprojekte (z. B. Anwendungs-/Service-Entwicklung, Plattform-/Konfigurationsautomation, Build- und Bereitstellungsprozesse). Sie beschreibt gewünschte Verhaltensweisen, Qualitätsprinzipien, Kollaborations- und Governance-Aspekte für Copilot-Vorschläge – unabhängig von konkreten Technologien oder einzelnen Tools. Konkrete produkt-/domänenspezifische Vorgaben oder Tool-Policies können ergänzend in separaten Dokumenten (z. B. ADRs, SECURITY, Betriebs-Runbooks) definiert werden.

> Normative Schlüsselwörter (MUST / MUST NOT / SHOULD / SHOULD NOT / MAY) lehnen sich an RFC 2119 an.

## AI Execution Contract (Zusatz für KI-gestützte Assistenz)

MUST: 
- Niemals Code vorschlagen, der ein neues Feature ohne zugehöriges Issue + Test einführt.
- Patch-Ausgaben minimal halten (nur betroffene Zeilen / Dateien, keine kosmetische Reformatierung).
- Vor Änderungen an Architektur-/Sicherheits-relevanten Dateien prüfen, ob eine passende ADR existiert und referenziert werden muss.
- Sicherheits- und Geheimnis-Policy beachten: Keine Secrets, Tokens, personenbezogenen Daten generieren, spiegeln oder im Klartext einfügen.
- Bei Unklarheit (fehlender Pfad, widersprüchliche Anforderungen) Rückfrage einleiten statt zu raten.

SHOULD:
- Tests zuerst ergänzen/ändern (Red) bevor Implementierung (Green) erfolgt.
- Jede Empfehlung mit kurzer Begründung (Warum / Effekt / Risiko) versehen.
- Fehlerausgaben interpretieren und konkrete Fix-Vorschläge liefern.

MAY:
- Kleinere, klar vorteilhafte Refactorings (Duplicate Removal, offensichtliche Naming-Verbesserungen) – sofern kein Scope Creep.

MUST NOT:
- Globale Formatierung (mass Reflow) ohne ausdrückliche Anweisung durchführen.
- Sicherheitsprüfungen „stumm“ entschärfen (z. B. `|| true` entfernen/ersetzen ohne Hinweis) oder bewusst umgehen.

Escalation / Rückfrage Kriterien (Rückfrage statt Änderung):
- Fehlende ADR bei architekturrelevantem Eingriff.
- Ungelöster Merge-Konflikt.
- Nicht reproduzierbarer Testfehler (unklarer deterministischer Zustand).

Validierungs-Erwartung:
- Nach bedeutender Änderung: Lint / Tests (mind. relevante Teil-Suite) gedanklich oder real ausführen und Status melden.
- Klare Kennzeichnung normativer Stufen in Befehlssequenzen (MUST / SHOULD / MAY).

Output Format Empfehlung (Antwortstruktur):
1. Kontext / Ziel (max 2 Sätze)
2. Delta-Plan (Stichpunkte)
3. Patch / Snippet
4. Validierung / Hinweise
5. Optionale Next Steps

---

## 1. Prinzipien & Grundlagen

1.1 **Zweck & Geltungsbereich**  
Diese Richtlinien definieren universelle Prinzipien für nachhaltige, sichere, nachvollziehbare Software- und Automatisierungs-Entwicklung. Sie adressieren Produkt-Features, Plattform-/Konfigurationsänderungen und operative Anpassungen gleichermaßen.

1.2 **Kernprinzipien**
- Transparenter, nachvollziehbarer Änderungsfluss (Issue → Branch → Review → Merge → Release → Betrieb → Feedback)
- Test- & Qualitätsorientierung (Fehlerprävention vor Nachbesserung)
- Least Privilege & Minimierung von Angriffs-/Fehleroberflächen
- Reproduzierbarkeit & Determinismus
- Evolvierbare Architektur (ADR-gesteuert statt implizit)
- Dokumentation fokussiert auf Entscheidungen & Betriebsrelevanz

1.3 **Engineering-Lifecycle (abstrakt)**
1. Planung  
2. Test First / TDD  
3. Implementierung (kleine, fokussierte Schritte)  
4. Review / Self-Audit  
5. Automatisiertes Testen (alle Ebenen)  
6. Sicherheits- & Qualitäts-Checks  
7. Build & Bereitstellung (pipeline-gesteuert)  
8. Betriebsbeobachtung & Drift-Erkennung  
9. Kontinuierliche Verbesserung / Lessons Learned  

1.4 **Sprache**  
Alle Kommentare, Commits, Issues, PR-Bodies und Richtlinien in Deutsch (Ausnahme: externe API-/Lib-Namen, Standardbegriffe). 

1.5 **Modularität & Struktur**
- Schichten / Domänen klar getrennt (kein „God“-Modul)
- Öffentliche Schnittstellen minimal & eindeutig dokumentiert
- Keine impliziten globalen Zustände; explizite Abhängigkeits-Injektion bevorzugt

1.6 **Versionierung & Nachverfolgung**
- `VERSION` als Single Source of Truth (SemVer) (MUST)
- `CHANGELOG.md` nach Keep a Changelog (MUST)
- Release-Prozess: Unreleased → Freeze → Bump → Tag → Veröffentlichung (MUST)

1.7 **Konfiguration & Laufzeit**
- Keine Hardcodierung sensibler / variabler Parameter (MUST)
- Health/Readiness Indikatoren für zentrale Funktionen (SHOULD)
- Isolierung interner Komponenten von externen Schnittstellen (SHOULD)
- Nur notwendige Exposition (Least Exposure) (MUST)
- Persistente vs. flüchtige Daten klar getrennt; Minimalzugriff (SHOULD)

1.8 **Dokumentation**
- Fokus auf Warum + Konsequenzen, nicht redundante Code-Nacherzählung (SHOULD)
- ADRs für jede signifikante Architektur-/Governance-Entscheidung (MUST)
- Runbooks für wiederkehrende Betriebsaufgaben (SHOULD)

## 2. Qualitäts- & Governance-Gates

2.1 **Tests & Qualitätssäulen**
- Unit → Integration → System / E2E: Pyramidenansatz
- Policy-/Konfigurationsprüfungen als Code (z. B. Validierungsregeln, statische Analysen)
- Keine neuen Features bei roten Tests

2.2 **Testanforderungen**
- Aussagekräftige, deterministische Tests (kein Sleep-basiertes Timing)
- Fehlerfälle & Edge Cases werden explizit adressiert
- Testcode gleichwertig gepflegt (Lesbarkeit, Refactor bei Geruch)

2.3 **Security & Geheimnisse**
- Keine Klartext-Secrets im Repo / Commit-Historie
- Geheimnisse nur über kontrollierte Kanäle (Secret Manager, verschlüsselte Konfiguration, o. ä.)
- Prinzip Geringster Rechte (Daten, Rollen, Services, Pipelines)
- Minimierung der Angriffsoberfläche (kein unnötiger Code / Endpunkt / Port)
- Härtung sensibler Pfade (AuthN, AuthZ, Rate Limits, Logging, Header/Policy Hardening)
- Supply Chain Schutz (Abhängigkeits-Scans, Signaturen, SBOM wo möglich)

2.4 **Automatisierte Qualitätsschranken**
Verbindliche Checks vor Merge / Deployment:
- Formatierung & Linting (MUST)
- Statische / Semantische Analyse (MUST)
- Dependency / Vulnerability Scan (keine offenen kritischen ohne Ausnahme) (MUST)
- Test-Suites grün (Unit + definierte höhere Ebenen) (MUST)
- (MAY) Dry-Run / Konfig-Validierung bei deklarativen Artefakten

2.5 **Architektur & ADR Disziplin**
- Jede Änderung prüft bestehende ADRs (kein stilles Override)
- Supersession über neue ADR mit Referenz, niemals Direktedit akzeptierter Historie
- Temporäre Abweichungen: Issue + Risiko + Ablaufdatum + Dokumentation im ADR Abschnitt "Known Deviations"

2.6 **Beobachtbarkeit & Betrieb**
- Metriken: Latenz, Durchsatz, Fehlerrate, Ressourcenverbrauch
- Strukturiertes Logging mit Korrelation / Trace IDs (keine sensiblen Inhalte)
- Drift-Erkennung (Soll vs. Ist) – Abweichung erzeugt Issue
- Definierte Wiederherstellungsziele & dokumentierte Backup-/Restore-Pfade

2.7 **Sicherer Änderungsfluss**
- Kein Direkt-Push auf `main` (MUST)
- Branch Protection + Status Checks obligatorisch (MUST)
- Commit Hygiene: kleine, thematisch fokussierte Schritte (SHOULD)

2.8 **Rollback & Reproduzierbarkeit**
- Jede Release-Version rekonstruierbar (Tag + Artefakt + Konfigurationsstand) (MUST)
- Keine „floating“ Produktionsabhängigkeiten ohne Versionsbindung (MUST)

2.9 **Qualitätsmetriken & Kontinuierliche Verbesserung**
- Regelmäßige Auswertung von Fehlerraten, MTTR, Teststabilität, Sicherheitsfunden
- Erkenntnisse → neue Issues / ADR Anpassungen

## 3. Operativer Workflow (konkret)

### 3.0 GitHub Integration (MCP Tools)

**Primäre Methode:** GitHub MCP Tools (direkte API-Integration, seit 2025-09-30)  
**Fallback:** `gh` CLI (nur wenn MCP Tools nicht verfügbar)

**Verfügbare GitHub MCP Tools:**
- `mcp_github_github_create_gist` / `mcp_github_github_update_gist` / `mcp_github_github_list_gists` - Gist Management
- `mcp_github_github_assign_copilot_to_issue` - Copilot Coding Agent Issue Assignment
- `mcp_github_github_create_pull_request_with_copilot` - Copilot PR Delegation
- `mcp_github_github_request_copilot_review` - Copilot Code Review Request
- `mcp_github_github_get_me` - Authentifizierter User Info
- `mcp_github_github_get_team_members` / `mcp_github_github_get_teams` - Team Management
- `mcp_github_github_list_starred_repositories` - Repository Discovery
- `mcp_github_github_get_discussion` / `mcp_github_github_list_discussions` - Discussion APIs
- `mcp_github_github_get_project` / `mcp_github_github_list_projects` - Project Management

**Zusätzliche Tool-Kategorien (via activate_* functions):**
- `activate_github_issue_management` - Issue CRUD, Comments, Sub-Issues
- `activate_github_pull_request_management` - PR CRUD, Reviews, Merge, Diffs
- `activate_github_repository_management` - Repo Creation, Files, Branches, Tags
- `activate_github_workflow_management` - GitHub Actions Workflows
- `activate_github_notification_management` - Notifications & Subscriptions
- `activate_github_search_tools` - Code/Issue/PR/Repo/User Search
- `activate_github_security_management` - Security Alerts & Advisories

**AI Execution Guidance:**
- PREFER: GitHub MCP Tools für Issue/PR/Gist Operations (direkter, typsicher)
- FALLBACK: `gh` CLI nur wenn MCP Tool fehlt oder Fehler auftritt
- VALIDATE: Nach MCP Tool-Call immer Ergebnis prüfen (Success/Error Handling)
- DOCUMENT: Bei neuem MCP Tool-Einsatz kurz im Commit erwähnen

**Beispiel (Issue Erstellung):**
```text
# PREFER (MCP Tools)
activate_github_issue_management → create_issue(title, body, labels)

# FALLBACK (gh CLI)
gh issue create --title "<Titel>" --body-file tmp/issue-body.md --label "feat"
```

### Quick Start (normativer Kurzablauf)
1. (MUST) Issue anlegen (Ziel + Akzeptanzkriterien + ADR Referenzen) → **MCP Tools bevorzugt**
2. (MUST) Branch vom aktuellen `main` erstellen (konformes Naming) → `git` CLI
3. (MUST) Failing Test hinzufügen (Red) – kein Produktionscode davor
4. (MUST) Minimalen Code schreiben bis Tests grün (Green) – kein Scope Creep
5. (SHOULD) Lint / Security lokal prüfen; (MAY) zusätzliche Scans
6. (MUST) PR erstellen mit Issue-Referenz + ADR IDs (falls relevant) → **MCP Tools bevorzugt**
7. (MUST) Alle Gates grün (Tests, Lint, Security) → Merge via PR → **MCP Tools bevorzugt**
8. (MUST) Versionierung / Follow-ups / Drift prüfen & Issues nachziehen
9. (MUST) Vor Issue-Abschluss Pläne, `CHANGELOG.md` und betroffene Dokumentation aktualisieren

Hinweis (ADR Referenz Enforcement): Bei Änderungen an Governance-/Architektur-relevanten Pfaden (`docs/adr/`, `scripts/`, `.github/`, `Makefile`, allgemeine `docs/`) MUSS der PR Body eine ADR Referenz (`ADR-XYZ`) enthalten oder explizit den Skip Marker `ADR-NOT-REQ` mit kurzer Begründung. Fehlt beides, schlägt das ADR Reference Gate fehl. Skip Marker führt nur zu einer Warnung (nicht-blockierend).

3.1 **Einordnung & Mapping**  
Der folgende Ablauf operationalisiert die Lifecycle-Schritte aus Abschnitt 1.3.

Implementierte Make Targets (seit v0.1.0, vollständig verfügbar):
- `make test` (führt Unit/Integration Tests aus)
- `make lint` (statische Analysen / Format / Stil)
- `make scan` (Security & Dependency Checks)
- `make pr-check` (bündelt: lint + test + scan für lokale PR-Vorbereitung)
- `make release VERSION=X.Y.Z` (Version bump + Changelog Transformation + Tag-Vorbereitung)
- `make ci-local` (Simulation definierter CI-Gates lokal)

Details / Vollständige Spezifikation: siehe `docs/make-targets-plan.md`.

Diese Targets können ab sofort anstelle der expliziten Befehlsblöcke in 3.2–3.9 verwendet werden (bevorzugte Methode).

**CI-Pipeline Integration:** 
- `make pr-check` eignet sich zur Verwendung in GitHub Actions Workflows
- `make ci-local` kann als lokaler CI-Simulator eingesetzt werden  
- Empfohlene Pipeline-Nutzung: `make ci-local` in CI-Environment für vollständige Gate-Simulation

| Lifecycle (1.3) | Operative Umsetzung (Abschnitt 3.x / Querverweis) |
|-----------------|---------------------------------------------------|
| 1 Planung | Issue-Erstellung (3.2) |
| 2 Test First / TDD | Erster Commit mit rotem Test vor Implementierung (3.3) |
| 3 Implementierung | Iterative Branch-Commits (3.3) |
| 4 Review / Audit | PR Body (3.4), PR Review & Workflow (3.5) |
| 5 Automatisiertes Testen | CI-Ausführung der Suites (3.6) |
| 6 Sicherheits-/Qualitäts-Checks | CI Gates / Scans (3.6 / Abschnitt 2.4) |
| 7 Build & Bereitstellung | Release / Tagging (3.7) |
| 8 Betriebsbeobachtung & Drift | Monitoring & Drift Tickets (2.6) |
| 9 Verbesserung | Lessons Learned → Issues / ADRs (2.9) |

3.2 **Issue → Branch**
- Jede Änderung startet mit einem Issue (Ziel + Akzeptanzkriterien + Referenzen zu ADRs)
- Branch-Naming: `feat/<kurz>`, `fix/<kurz>`, `chore/<kurz>`, `refactor/<kurz>` oder Issue-ID basiert
- Keine parallele Vermischung mehrerer unzusammenhängender Ziele
 - Nutzung der bereitgestellten Issue Templates (Feature / Bug) wird erwartet (GitHub UI Auswahl). Blank Issues sind deaktiviert.
 - Issue-Bodies MÜSSEN gültiges Markdown sein (Headings für Kontext / Akzeptanzkriterien, ungefüllte Sektionen als Platzhalter klar markiert). Reiner Fließtext ohne Struktur gilt als unvollständig.
 - Sicherheits-Hinweis: In `tmp/` abgelegte Issue-/PR-Bodies dürfen keine sensitiven Inhalte (Secrets, personenbezogene Daten, Zugangstokens) enthalten. Dateien gelten als temporär und können rotierend gelöscht werden.
 - Befehle (Erstellung & Vorbereitung) (MUST):
	 ```bash
	 # Issue erstellen (interaktiv)
	 gh issue create --title "<Titel>" --body-file tmp/issue-<nr>-body.md --label feat

	 # Alle offenen Issues anzeigen
	 gh issue list --limit 30

	 # Branch aus Issue heraus anlegen (gh extension workflow, falls vorhanden) – sonst manuell:
	 # (Falls 'gh issue develop' nicht verfügbar: Branch manuell erstellen)
	 git fetch origin
	 git checkout -b feat/<slug> origin/main

	 # Issue im Browser öffnen (MAY)
	 gh issue view <ISSUE_NR> --web
	 ```

3.3 **Implementierung & Tests**
- Erster Commit setzt (mind.) einen fehlenden Test (rot) für neues Verhalten
- Kleinste sinnvolle Schritte; früh grüner Zustand wiederherstellen
- Kein Scope Creep: Zusätzliche Ideen → neue Issues
 - Befehle (Beispielablauf): Tests (MUST) / Lint (MAY):
	 ```bash
	 # Status prüfen
	 git status

	 # Neuen (roten) Test hinzufügen
	 git add tests/unit/<neuer_test>.go
	 git commit -m "test: spezifiziert neues Verhalten <kurz>"  # failing expected

	 # Implementierung anpassen
	 git add internal/<pfad>/logic.go
	 git commit -m "feat: implementiert Basislogik für <kurz>" 

	 # Lokale Tests ausführen (Beispiel Make Target)
	 make test

	 # Lint (MAY)
	 make lint || true

	 # Änderungen pushen
	 git push -u origin feat/<slug>
	 ```

3.4 **Pull Request Body (Copilot Workspace)**

**Closing-Keywords (MUST):**
- Immer das Schlüsselwort **Closes Issue <NR>** verwenden (konsistente Formulierung bevorzugt)
- Position: **Erste Zeile** des PR Body (vor `## Overview` oder anderen Headings)
- Syntax: **Closes Issue 123** (ohne Repository-Prefix bei gleichem Repo) oder **Closes <owner>/<repository> Issue 123** (mit Prefix)
- GitHub Action `.github/workflows/pr-closing-keyword-fix.yml` korrigiert automatisch falsche Platzierungen
- Auch **Fixes Issue <NR>** oder **Resolves Issue <NR>** werden erkannt und zu Closes normalisiert

**Beispiel (korrekte Struktur):**
```markdown
Closes Issue 123

## Overview
Implements feature XYZ...

## Changes
- Added new functionality
- Updated tests
```

**Begründung:** 
- Zuverlässigste Auto-Close Methode für Issues
- Keine Template-Abhängigkeit (Copilot-kompatibel)
- Automatische Korrektur durch GitHub Action verhindert fehlerhafte Platzierung
- Konsistenz über alle PRs (manuell + Copilot-generiert)

**Vermeidung:**
- ❌ Keyword nach `</details>` oder HTML-Tags
- ❌ Keyword in Code-Blöcken (````markdown ... ```)
- ❌ Keyword nach Zeile 200 (zu weit unten)
- ❌ Eingerücktes Keyword (Listen, Quotes)

3.5 **Pull Request Review & Workflow**
- PR referenziert Issue & relevante ADR IDs
- PR Body: Was / Warum / Risiken / Testnachweis (kurz) – per Datei (`--body-file`)
- Alle Diskussionen geklärt vor Merge
 - Befehle: PR Erstellung (MUST) / Labels & Diff Web (MAY):
	 ```bash
	 # Preview Diff
	 gh pr diff --web || true   # Falls bereits ein PR existiert

	 # PR erstellen (automatische Befüllung: Titel aus erstem Commit, Body interaktiv)
	 gh pr create --base main --head feat/<slug> --title "feat: <kurztitel>" \
		 --body-file tmp/pr-<slug>-body.md

	 # PR anzeigen
	 gh pr view --web

	 # Labels setzen (falls nötig)
	 gh pr edit <PR_NR> --add-label "feature"

	 # ADR Referenzen im Body validieren (manuell / MAY: zusätzlicher Check-Script)
	 ```

3.6 **Pipeline & Gates**
- Vollständige Ausführung definierter Checks (Tests, Lint, Scans, Validierungen)
- Hard Fail bei kritischen Sicherheitsfunden ohne Ausnahme-Issue
- Keine manuelle Umgehung gesperrter Gates
- **Auto-Approve für Bots:** Copilot/Dependabot PRs benötigen keine manuelle Workflow-Approval (siehe `docs/ci-cd/bot-workflow-approval.md`)
 - Befehle / Checks (Beispiele): Tests, Lint, Security (MUST) / Format separat, zusätzliche Scans (MAY):
	 ```bash
	 # Tests
	 make test

	 # Lint (Beispiel)
	 make lint-ci

	 # Security Scan (z. B. mit Trivy oder vergleichbaren Tools)
	 make scan

	 # Format (optional)
	 make lint

	 # PR Checks ansehen
	 gh pr checks <PR_NR>
	 ```

3.7 **Versionierung & Release**
- Unreleased Abschnitt bereinigen → neue Version vergeben → `VERSION` aktualisieren
- Commit Message: `chore: Version auf X.Y.Z erhöht`
- Tag erstellen & pushen (`vX.Y.Z`)
 - Befehle: Version + Tag (MUST):
	 ```bash
	 # Version setzen (Beispiel)
	 echo "X.Y.Z" > VERSION
	 sed -i "s/^## \[Unreleased\]/## [Unreleased]\n\n## [X.Y.Z] - $(date +%Y-%m-%d)/" CHANGELOG.md

	 git add VERSION CHANGELOG.md
	 git commit -m "chore: Version auf X.Y.Z erhöht"
	 git push origin feat/<slug>

	 # (Nach Merge auf main) Tag erstellen
	 git fetch origin
	 git checkout main
	 git pull --ff-only
	 git tag -a vX.Y.Z -m "Version X.Y.Z"
	 git push origin vX.Y.Z
	 ```

3.8 **Merge & Clean-up**
- Merge nur via PR (kein Rebase-Rewrite der Historie)
- Branch löschen nach erfolgreichem Merge
- Keine nachträglichen Commits auf gemergte Feature Branches
- Merge-Vorgang nach Freigabe ausschließlich mit `git` CLI durchführen (kein Merge-Button, kein `gh pr merge`)
 - Befehle: Merge (MUST) / Manuelle Prüfung entfernte Branches (MAY):
	 ```bash
	 # Merge-Ablauf (nach Approval)
	 git checkout main
	 git pull --ff-only
	 git merge --ff-only origin/<branch>
	 git push origin main

	 # Remote-Branch entfernen
	 git push origin --delete <branch>
	 git fetch --prune
	 git branch -a | grep <slug> || echo "Branch entfernt"
	 ```

Hinweis: Der finale Push auf `main` im Zuge des Merge-Vorgangs gilt als einzige Ausnahme zur "kein Direkt-Push"-Regel und setzt eine freigegebene PR voraus.

3.9 **Post-Merge & Betrieb**
- Monitoring prüfen (Anomalien?)
- Offene Folgeaufgaben als Issues erfassen (keine stillen TODOs)
- ADR Supersessions bei Architekturfolgen zeitnah anstoßen
- Dokumentationsartefakte (Pläne, `CHANGELOG.md`, Nutzer-/Ops-Doku) auf aktuellen Stand bringen
 - Befehle / Nachbereitung: TODO Scan & Folge-Issue (MUST) / ADR Draft sofort (MAY):
	 ```bash
	 # Letzten Merge-Commit anzeigen
	 git log -1 --oneline

	 # Offene TODO Marker scannen (Beispiel)
	 grep -R "TODO(" -n . || true

	 # Neues Folge-Issue (Beispiel)
	 gh issue create --title "Follow-up: <konkret>" --body-file tmp/issue-followup-<slug>.md --label chore

	 # ADR Vorschlag Template kopieren (falls neue Entscheidung nötig)
	 cp docs/adr/000-template.md docs/adr/ADR-<nr>-<slug>.md
	 git add docs/adr/ADR-<nr>-<slug>.md
	 git commit -m "docs(adr): propose <slug>"
	 git push origin <branch>  # falls noch nicht merged oder neuer Branch
	 ```

3.10 **Kontinuierliche Verbesserung**
- Retro / Lessons Learned bündeln → Prozess-/Test-/Policy-Anpassungen
- Metriken und Vorfälle fließen in priorisierte Verbesserungs-Issues
 - Hilfsbefehle (MAY) (Metriken / Übersicht):
	 ```bash
	 # Offene Issues nach Label sortieren
	 gh issue list --label improvement --limit 50

	 # Statistische einfache Übersicht (Commits letzte 7 Tage)
	 git log --since="7 days ago" --oneline | wc -l

	 # PRs der letzten Woche
	 gh pr list --state merged --limit 30 --search "merged:>=$(date -I -d '7 days ago')"
	 ```

## 4. Hinweise

- Änderungen an diesen Richtlinien folgen selbst dem definierten Workflow.
- Tool-/Technologie-spezifische Leitlinien werden ausgelagert (z. B. separate Infrastruktur-, API- oder UI-Guides).
- Erweiterbare Referenzen: Security-Guides, Architekturübersichten, Betriebs-Runbooks, Qualitäts- & Metrik-Standards, Lieferketten-/Supply-Chain-Richtlinien.
- Konsistenzprüfung der Norm-Labels lokal: `scripts/common/check-normative.sh` ausführen (CI-Integration empfohlen).
 - Pre-Commit Aktivierung (lokal): `git config core.hooksPath .githooks` (führt u. a. normative Check vor jedem Commit aus).
 - Issue Templates: Siehe `.github/ISSUE_TEMPLATE/` – Issues sollten strukturiert mit Markdown-Headings erstellt werden.
 - PR Body Guidelines: Siehe Abschnitt 3.4 – Closing-Keywords werden automatisch via GitHub Action korrigiert (`.github/workflows/pr-closing-keyword-fix.yml`).
 - ADR Erstellung & Prüfung: Neue ADR via `make adr-new SLUG=<slug>` (Template: `docs/adr/000-template.md`), Validierung mit `scripts/common/check-adr.sh` (läuft auch im CI Workflow `quality-gates`).
 - Issue Body Qualität: CI / zukünftige Hooks KÖNNEN fehlende Pflichtsektionen (Kontext, Ziel, Akzeptanzkriterien) beanstanden.

## 5. Glossar (Platzhalter / Konventionen)
### Normative Legende
| Begriff | RFC Ebene |
|---------|-----------|
| MUST | Verbindlich, keine Abweichung ohne dokumentierte Ausnahme |
| SHOULD | Starke Empfehlung – Abweichung nur mit Begründung |
| MAY | Situativ, kein Gate |
| MUST NOT | Verbot – nicht implementieren |
| SHOULD NOT | Vermeiden – nur in begründeten Ausnahmefällen |


### ADR Kurzreferenz (Zusammenfassung)

| Status | Bedeutung | Aktion bei Änderung |
|--------|-----------|---------------------|
| Proposed | In Diskussion | Review, ggf. Anpassungen einpflegen |
| Accepted | Verbindlich | Bei Widerspruch neue ADR (Supersede) |
| Superseded | Historisch, ersetzt | Nicht mehr anwenden, nur Referenz |
| Deprecated | Nutzung auslaufend | Migration planen / Issue anlegen |
| Rejected | Abgelehnt | Nicht implementieren |

Hinweis: Architekturänderungen ohne passende Accepted ADR sind zu vermeiden (MUST). Supersessions immer bidirektional referenzieren (MUST).

| Platzhalter | Bedeutung | Hinweis |
|-------------|-----------|---------|
| `<slug>` | Kurzlesbarer Bezeichner (kebab-case) für Branch/Feature | Keine Leer- oder Sonderzeichen, eindeutig im Kontext |
| `<nr>` | Laufende Nummer / Issue-ID | Entspricht exakter Plattform-ID (z. B. GitHub Issue Nummer) |
| `<ISSUE_NR>` | Explizite Referenz auf eine Issue-ID | Großschreibung zur Hervorhebung im Befehlsbeispiel |
| `<PR_NR>` | PR Nummer | Wird nach PR-Erstellung sichtbar (gh pr view) |
| `<branch>` | Aktueller Arbeits- oder Ziel-Branch | Vermeide direkte Arbeit auf `main` |
| `<kurz>` | Kurzbeschreibung eines Features / Tests / Logikbestandteils | Max. ~3 Worte, semantisch aussagekräftig |
| `X.Y.Z` | SemVer Version | MAJOR.MINOR.PATCH |
| `<pfad>` | Relativer Pfad zu Quellcode / Modul | Projektkonventionen beachten |
| `<konkret>` | Frei ausformulierter spezifischer Folgeschritt | Sollte test-/issue-fähig sein |
| `<Titel>` | Volltext-Titel eines Issues / PR | Prägnant, Ergebnis-orientiert |
| `<kurztitel>` | Kompakte PR Titel-Variante | Für Übersicht / Listenansichten |

Konventionen:
- Platzhalter nicht wörtlich übernehmen – immer konkret ersetzen.
- Neue Platzhalter erst nach Dokumentation im Glossar verwenden.
- Sensible Inhalte niemals in Platzhalter-Beispielen zeigen.
 - Für implementierte Automatisierungs-Targets siehe `docs/make-targets-plan.md` (Vollständige Spezifikation).

---

### Maschinenlesbarer Norm-Export (JSON)
Dieser Block KANN von Automatisierung / Policies geparst werden (Single Source: dieses Dokument).

```json
{
	"normative": {
		"no_direct_push_main": "MUST",
		"tests_green_before_merge": "MUST",
		"security_scans_blockers": "MUST",
		"version_single_source": "MUST",
		"adr_reference_on_arch_changes": "MUST",
		"lint_before_pr": "SHOULD",
		"format_consistency": "SHOULD",
		"metrics_review": "MAY"
	}
}
```
