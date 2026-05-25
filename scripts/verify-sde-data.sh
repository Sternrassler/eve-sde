#!/usr/bin/env bash
# Verifiziert die SDE-Daten-Korrektheit nach Regenerierung + Import.
# Nutzung: scripts/verify-sde-data.sh [db-pfad]
set -euo pipefail

DB="${1:-data/sqlite/eve-sde.db}"
JSONL_DIR="data/jsonl"
fail=0

echo "=== (a) Jede JSONL-Datei hat eine Tabelle (minus Ausschlüsse) ==="
for f in "$JSONL_DIR"/*.jsonl; do
  n="$(basename "$f" .jsonl)"
  cnt="$(sqlite3 "$DB" "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='$n';")"
  if [ "$cnt" = "0" ]; then
    echo "FEHLT: $n"
    fail=1
  fi
done
[ "$fail" = "0" ] && echo "OK: alle JSONL-Dateien haben eine Tabelle"

echo "=== (b) capacity-Spalte ist REAL ==="
captype="$(sqlite3 "$DB" "SELECT type FROM pragma_table_info('types') WHERE name='capacity';")"
echo "capacity type = $captype"
[ "$captype" = "REAL" ] || { echo "FEHLER: capacity ist nicht REAL"; fail=1; }

echo "=== (c) _key=484 behält capacity=0.5 ==="
cap484="$(sqlite3 "$DB" "SELECT capacity FROM types WHERE _key=484;")"
echo "capacity(_key=484) = $cap484"
[ "$cap484" = "0.5" ] || { echo "FEHLER: capacity(_key=484) != 0.5"; fail=1; }

echo "=== (d) Keine Tabelle ist leer ==="
while read -r tbl; do
  rows="$(sqlite3 "$DB" "SELECT COUNT(*) FROM \"$tbl\";")"
  if [ "$rows" = "0" ]; then
    echo "LEER: $tbl"
    fail=1
  fi
done < <(sqlite3 "$DB" "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
[ "$fail" = "0" ] && echo "OK: keine leere Tabelle"

if [ "$fail" = "0" ]; then
  echo "=== ALLE CHECKS BESTANDEN ==="
else
  echo "=== CHECKS FEHLGESCHLAGEN ==="; exit 1
fi
