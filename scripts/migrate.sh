#!/usr/bin/env sh

# We can't use command `timeout` because it has different arguments on different versions of BusyBox (alpine)
# BusyBox v1.29.3 | Usage: timeout [-t SECS] [-s SIG] PROG ARGS
# BusyBox v1.31.1 | Usage: timeout [-s SIG] SECS PROG ARGS

SERVICE=$1

echo "Migrating ${DIRECTION:-up} service ${SERVICE}"
echo "Using SVC_DB_HOST=${SVC_DB_HOST} SVC_DB_PORT=${SVC_DB_PORT} SVC_DB_DATABASE=${SVC_DB_DATABASE} SVC_DB_USER_NAME=${SVC_DB_USER_NAME}"

# check sql-migrate installed
if ! [ -x "$(command -v sql-migrate)" ]; then
  echo "FATAL: sql-migrate not installed, to install:
    go install github.com/rubenv/sql-migrate/sql-migrate"
  exit 127
fi

TIMEOUT=10
sleep $TIMEOUT &
TIMEOUT_PID=$!

cd "${SERVICE}/internal/migrations" || exit 3

while ! sql-migrate status; do
  if ! kill -0 $TIMEOUT_PID 2>/dev/null; then
    echo "Failed to connect DB during $TIMEOUT sec"
    exit 1
  fi
  echo "Sleeping 1 sec..."
  sleep 1
done

sql-migrate "${DIRECTION:-up}"
