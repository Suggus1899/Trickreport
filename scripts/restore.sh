#!/usr/bin/env bash
# Restore a trusted pg_dump SQL file into a NEW, explicitly named rehearsal DB.
# Compose discovers the service/network and reads COMPOSE_FILE/COMPOSE_PROJECT_NAME.
set -euo pipefail

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <dump.sql[.gz]> <trickreport_restore_name>" >&2
  exit 1
fi
BACKUP_FILE="$1"
RESTORE_DB="$2"
if [[ ! "$RESTORE_DB" =~ ^trickreport_restore_[a-zA-Z0-9_]+$ ]] || [ "${#RESTORE_DB}" -gt 63 ]; then
  echo 'Destination must be a new database named trickreport_restore_<name> (max 63 characters).' >&2
  exit 1
fi
if [ ! -r "$BACKUP_FILE" ] || [ ! -s "$BACKUP_FILE" ]; then
  echo 'The dump must be readable and nonempty.' >&2
  exit 1
fi

SQL_FILE="$(mktemp)"
trap 'rm -f "$SQL_FILE"' EXIT
case "$BACKUP_FILE" in
  *.gz) gzip -t -- "$BACKUP_FILE"; gzip -cd -- "$BACKUP_FILE" > "$SQL_FILE" ;;
  *) cat -- "$BACKUP_FILE" > "$SQL_FILE" ;;
esac
if ! grep -q '^-- PostgreSQL database dump' "$SQL_FILE" || ! grep -Eq '^(CREATE|COPY|INSERT|ALTER|SET|SELECT) ' "$SQL_FILE"; then
  echo 'Expected a nonempty, plain SQL pg_dump file without --create.' >&2
  exit 1
fi
if grep -Eiq '^[[:space:]]*((CREATE|DROP|ALTER)[[:space:]]+DATABASE|\\(connect|c|!|i|ir|include|include_relative)([[:space:]]|$))' "$SQL_FILE"; then
  echo 'Dump may not create, drop, alter or switch databases, or include external commands/files.' >&2
  exit 1
fi

# createdb fails if the destination exists; never drop or replace the source DB.
docker compose exec -T postgres sh -c '
  test "$1" != "$POSTGRES_DB" || exit 1
  export PGPASSWORD="$POSTGRES_PASSWORD"
  exec createdb --username="$POSTGRES_USER" --maintenance-db=postgres -- "$1"
' sh "$RESTORE_DB"
if ! docker compose exec -T postgres sh -c '
  export PGPASSWORD="$POSTGRES_PASSWORD"
  exec psql --username="$POSTGRES_USER" --no-psqlrc --set=ON_ERROR_STOP=1 --single-transaction --dbname="$1"
' sh "$RESTORE_DB" < "$SQL_FILE"; then
  echo "Restore failed; source database unchanged. Inspect the rehearsal database $RESTORE_DB before removing it." >&2
  exit 1
fi
echo "Restored to $RESTORE_DB. Validate schema, counts and tenant relationships before any separate cutover."
