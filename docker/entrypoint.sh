#!/bin/sh
set -e


DB_HOST="${DB_HOST:-db}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
MAX_RETRIES="${MAX_RETRIES:-30}"
SLEEP_SEC="${SLEEP_SEC:-1}"


try_pg_isready() {
if command -v pg_isready >/dev/null 2>&1; then
pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" >/dev/null 2>&1
return $?
fi
return 2
}


try_nc() {
if command -v nc >/dev/null 2>&1; then
nc -z "$DB_HOST" "$DB_PORT" >/dev/null 2>&1
return $?
fi
return 2
}


i=0
until { try_pg_isready || try_nc; }; do
i=$((i+1))
echo "waiting for db... attempt $i/$MAX_RETRIES"
if [ "$i" -ge "$MAX_RETRIES" ]; then
echo "db did not become ready after $MAX_RETRIES attempts" >&2
exit 1
fi
sleep "$SLEEP_SEC"
done


echo "db is ready — starting pr-service"
exec /app/pr-service