#!/usr/bin/env bash
# check-commit-msg.sh – Validiert Commit Message Konventionen (Conventional Commits)
# Referenz: copilot-instructions.md Abschnitt 2.7

set -euo pipefail

# Usage: check-commit-msg.sh [--range RANGE] [--file FILE]
MODE="stdin"
RANGE=""
FILE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --range)
            MODE="range"
            RANGE="$2"
            shift 2
            ;;
        --file)
            MODE="file"
            FILE="$2"
            shift 2
            ;;
        *)
            echo "[check-commit-msg] ERROR: Unbekannter Parameter: $1" >&2
            exit 1
            ;;
    esac
done

# Conventional Commits Regex (vereinfacht)
# Format: type(scope?): subject
COMMIT_REGEX="^(feat|fix|docs|style|refactor|test|chore|perf|ci|build|revert)(\(.+\))?: .{3,}"

validate_message() {
    local msg="$1"
    local commit_hash="${2:-}"
    
    # Erster Zeile extrahieren
    first_line=$(echo "$msg" | head -n1)
    
    # GitHub Copilot Bot Commits überspringen (automatisch generiert)
    if [ -n "$commit_hash" ]; then
        commit_author=$(git log -1 --pretty=%an "$commit_hash" 2>/dev/null || echo "")
        if echo "$commit_author" | grep -qE "copilot.*bot|Copilot.*Agent"; then
            echo "  ⏩ Übersprungen (Copilot Bot Commit): '$first_line' (Author: $commit_author)"
            return 0
        fi
    fi
    
    # GitHub Merge-Commits überspringen (automatisch generiert)
    if echo "$first_line" | grep -qE "^Merge [0-9a-f]{40} into [0-9a-f]{40}$"; then
        echo "  ⏩ Übersprungen (GitHub Merge-Commit): '$first_line'"
        return 0
    fi
    
    # Standard Merge-Commits überspringen (z.B. "Merge branch 'feat/xyz'")
    if echo "$first_line" | grep -qE "^Merge (branch|pull request|remote-tracking branch)"; then
        echo "  ⏩ Übersprungen (Merge-Commit): '$first_line'"
        return 0
    fi
    
    if ! echo "$first_line" | grep -qE "$COMMIT_REGEX"; then
        echo "  ❌ Ungültige Commit Message: '$first_line'"
        if [ -n "$commit_hash" ]; then
            echo "     Commit: $commit_hash"
        fi
        echo "     Erwartetes Format: type(scope?): subject"
        echo "     Erlaubte types: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert"
        return 1
    fi
    
    return 0
}

errors=0

case $MODE in
    stdin)
        # Git Hook Modus (pre-commit / commit-msg)
        msg=$(cat)
        echo "[check-commit-msg] Prüfe Commit Message..."
        if ! validate_message "$msg"; then
            ((errors++))
        fi
        ;;
    file)
        # Datei-Modus (für Git Hook mit .git/COMMIT_EDITMSG)
        if [ ! -f "$FILE" ]; then
            echo "[check-commit-msg] ERROR: Datei $FILE nicht gefunden" >&2
            exit 1
        fi
        msg=$(cat "$FILE")
        echo "[check-commit-msg] Prüfe Commit Message aus $FILE..."
        if ! validate_message "$msg"; then
            ((errors++))
        fi
        ;;
    range)
        # Range Modus (für CI – prüft alle Commits in Range)
        echo "[check-commit-msg] Prüfe Commit Messages in Range: $RANGE..."
        while IFS= read -r commit_hash; do
            [ -z "$commit_hash" ] && continue
            msg=$(git log -1 --pretty=%B "$commit_hash")
            if ! validate_message "$msg" "$commit_hash"; then
                ((errors++))
            fi
        done < <(git log --pretty=%H "$RANGE")
        ;;
esac

if [ "$errors" -gt 0 ]; then
    echo "[check-commit-msg] ❌ $errors ungültige Commit Message(s) gefunden"
    exit 1
fi

echo "[check-commit-msg] ✅ Commit Message(s) validiert"
exit 0
