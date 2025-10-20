#!/usr/bin/env bash
# check-security-blockers.sh – Parst Trivy JSON Report auf kritische Findings
# Referenz: copilot-instructions.md Abschnitt 2.3, 2.4

set -euo pipefail

TRIVY_REPORT="tmp/trivy-fs-report.json"

if [ ! -f "$TRIVY_REPORT" ]; then
    echo "[check-security-blockers] WARNING: Trivy Report $TRIVY_REPORT nicht gefunden – überspringe"
    exit 0
fi

echo "[check-security-blockers] Prüfe Security Findings in $TRIVY_REPORT..."

# Parse JSON und zähle HIGH/CRITICAL Vulnerabilities
if ! command -v jq >/dev/null 2>&1; then
    echo "[check-security-blockers] ERROR: jq nicht installiert – kann Report nicht parsen" >&2
    exit 1
fi

critical_count=$(jq '[.Results[]?.Vulnerabilities[]? | select(.Severity == "CRITICAL")] | length' "$TRIVY_REPORT" || echo "0")
high_count=$(jq '[.Results[]?.Vulnerabilities[]? | select(.Severity == "HIGH")] | length' "$TRIVY_REPORT" || echo "0")

echo "  - CRITICAL: $critical_count"
echo "  - HIGH: $high_count"

total_blockers=$((critical_count + high_count))

if [ "$total_blockers" -gt 0 ]; then
    echo "[check-security-blockers] ❌ $total_blockers Security Blocker gefunden"
    echo "Details:"
    jq -r '.Results[]?.Vulnerabilities[]? | select(.Severity == "CRITICAL" or .Severity == "HIGH") | "  - \(.VulnerabilityID): \(.PkgName) \(.InstalledVersion) (\(.Severity))"' "$TRIVY_REPORT" || true
    echo ""
    echo "MUST: Kritische Findings vor Merge beheben oder Ausnahme-Issue erstellen"
    exit 1
fi

echo "[check-security-blockers] ✅ Keine kritischen Security Findings"
exit 0
