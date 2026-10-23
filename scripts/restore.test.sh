#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEST_DIR="$(mktemp -d)"
trap 'rm -rf "$TEST_DIR"' EXIT
mkdir "$TEST_DIR/bin"
export RESTORE_TEST_LOG="$TEST_DIR/docker.log" POSTGRES_PASSWORD=synthetic-test-password
cat > "$TEST_DIR/bin/docker" <<'SH'
#!/usr/bin/env bash
printf '%q ' "$@" >> "$RESTORE_TEST_LOG"
printf '\n' >> "$RESTORE_TEST_LOG"
if [[ "$*" == *psql* ]]; then
  cat >/dev/null
  exit "${RESTORE_TEST_SQL_EXIT:-0}"
fi
SH
chmod +x "$TEST_DIR/bin/docker"
export PATH="$TEST_DIR/bin:$PATH"
printf '%s\n' '-- PostgreSQL database dump' 'CREATE TABLE synthetic_test (id integer);' > "$TEST_DIR/valid.sql"
gzip -c "$TEST_DIR/valid.sql" > "$TEST_DIR/valid.sql.gz"
printf 'not gzip' > "$TEST_DIR/corrupt.sql.gz"
touch "$TEST_DIR/empty.sql"
expect_rejected() {
  : > "$RESTORE_TEST_LOG"
  if bash "$ROOT/scripts/restore.sh" "$@" > "$TEST_DIR/output" 2>&1; then
    echo "FAIL: restore accepted $*" >&2
    exit 1
  fi
  test ! -s "$RESTORE_TEST_LOG" || { echo 'FAIL: invalid input reached Docker' >&2; exit 1; }
}
expect_rejected "$TEST_DIR/valid.sql"
expect_rejected "$TEST_DIR/valid.sql" trickreport
expect_rejected "$TEST_DIR/empty.sql" trickreport_restore_test
expect_rejected "$TEST_DIR/corrupt.sql.gz" trickreport_restore_test
expect_rejected "$TEST_DIR/valid.sql" 'trickreport_restore_x;drop_database'
printf '%s\n' 'not a pg_dump' > "$TEST_DIR/not-sql.txt"
expect_rejected "$TEST_DIR/not-sql.txt" trickreport_restore_test
printf '%s\n' '-- PostgreSQL database dump' '\connect trickreport' > "$TEST_DIR/switch.sql"
expect_rejected "$TEST_DIR/switch.sql" trickreport_restore_test
cp "$TEST_DIR/valid.sql.gz" "$TEST_DIR/dump with spaces.sql.gz"
for file in "$TEST_DIR/valid.sql" "$TEST_DIR/valid.sql.gz" "$TEST_DIR/dump with spaces.sql.gz"; do
  : > "$RESTORE_TEST_LOG"
  bash "$ROOT/scripts/restore.sh" "$file" trickreport_restore_test > "$TEST_DIR/output" 2>&1
  grep -q 'compose exec -T postgres.*createdb' "$RESTORE_TEST_LOG"
  grep -q 'psql.*--single-transaction.*trickreport_restore_test' "$RESTORE_TEST_LOG"
  if grep -qi 'DROP\|docker run\|--network' "$RESTORE_TEST_LOG"; then exit 1; fi
done
export RESTORE_TEST_SQL_EXIT=3
if bash "$ROOT/scripts/restore.sh" "$TEST_DIR/valid.sql" trickreport_restore_test > "$TEST_DIR/output" 2>&1; then
  echo 'FAIL: SQL error reported success' >&2
  exit 1
fi
echo 'PASS: invalid inputs rejected before Docker; isolated transaction restore; SQL errors propagate.'
