#!/usr/bin/env sh

HOST="$1"
PORT="$2"
COMMAND="while ! nc -vz $HOST $PORT; do echo Sleeping... && sleep 1; done"
timeout 10s sh -c "$COMMAND"
