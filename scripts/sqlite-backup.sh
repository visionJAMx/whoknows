#!/usr/bin/env bash
# SQLite laver et konsistent snapshot; kopier ikke en aktiv .db-fil med cp.
set -euo pipefail
umask 077

if [ "$#" -ne 2 ]; then
    echo "Brug: bash scripts/sqlite-backup.sh DATABASE OUTPUTMAPPE" >&2
    exit 1
fi
command -v sqlite3 >/dev/null || { echo "sqlite3 skal være installeret." >&2; exit 1; }
[ -f "$1" ] && [ -s "$1" ] || { echo "Databasen mangler eller er tom." >&2; exit 1; }

# Absolut sti gør kommandoen uafhængig af den efterfølgende mappeændring.
source_db="$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
mkdir -p "$2"
output_root="$(cd "$2" && pwd)"
snapshot_dir="$(mktemp -d "$output_root/sqlite-$(date -u +%Y%m%dT%H%M%SZ)-XXXXXX")"

# Unik, privat mappe: eksisterende databaser og backups overskrives aldrig.
# Ved fejl bevares mappen til fejlsøgning, men uden SUCCESS-markering.
cd "$snapshot_dir"
sqlite3 -batch -bail -readonly "$source_db" '.timeout 10000' '.backup snapshot.db'
check="$(sqlite3 -batch -bail -readonly snapshot.db 'PRAGMA integrity_check;')"
if [ "$check" != 'ok' ]; then
    echo "Backup fejlede integritetskontrol: $snapshot_dir" >&2
    exit 1
fi
touch SUCCESS
printf '%s\n' "$snapshot_dir/snapshot.db"
