#!/usr/bin/env bash
#
# Trickreport — restore a PostgreSQL backup.
#
# Usage:
#   ./scripts/restore.sh <backup-file>
#
# Example:
#   ./scripts/restore.sh backups/trickreport_2024-01-01T00:00:00Z.sql.gz
#
# Requirements:
#   - docker compose with the main docker-compose.yml running a healthy postgres.
#   - The backup file must be readable by the postgres container (mounted under /backups).
#
set -euo pipefail

if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <backup-file>" >&2
  echo "Example: $0 backups/trickreport_2024-01-01T00:00:00Z.sql.gz" >&2
  exit 1
fi

BACKUP_FILE="$1"

if [ ! -f "$BACKUP_FILE" ]; then
  echo "Error: backup file not found: $BACKUP_FILE" >&2
  exit 1
fi

# Resolve to an absolute path so we can mount it reliably.
case "$BACKUP_FILE" in
  /*) ABS_PATH="$BACKUP_FILE" ;;
  *)  ABS_PATH="$(cd "$(dirname "$BACKUP_FILE")" && pwd)/$(basename "$BACKUP_FILE")" ;;
esac

FILE_NAME="$(basename "$ABS_PATH")"

# Load env (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB) if a .env file exists.
if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env
  set +a
fi

POSTGRES_USER="${POSTGRES_USER:-trickreport}"
POSTGRES_DB="${POSTGRES_DB:-trickreport}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:?POSTGRES_PASSWORD is required}"

echo "==> Restoring backup: $FILE_NAME"
echo "    Database: $POSTGRES_DB (user: $POSTGRES_USER)"
echo

# Detect compression by extension.
case "$FILE_NAME" in
  *.gz)
    CMD="gunzip -c /backups/$FILE_NAME | psql -v ON_ERROR_STOP=1"
    ;;
  *)
    CMD="psql -v ON_ERROR_STOP=1 -f /backups/$FILE_NAME"
    ;;
esac

# Run the restore inside a throwaway postgres container on the same network.
docker run --rm \
  -v "$ABS_PATH:/backups/$FILE_NAME:ro" \
  -e PGPASSWORD="$POSTGRES_PASSWORD" \
  --network trickreport_trickreport_net \
  --entrypoint /bin/sh \
  postgres:16-alpine \
  -c "psql -h postgres -U $POSTGRES_USER -d $POSTGRES_DB -v ON_ERROR_STOP=1 -c 'DROP SCHEMA public CASCADE; CREATE SCHEMA public;' && $CMD -h postgres -U $POSTGRES_USER -d $POSTGRES_DB"

echo
echo "==> Restore complete."
