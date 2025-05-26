#!/bin/sh

HOST=$1
PORT=$2
TIMEOUT=${3:-15}

shift 3
CMD="$@"

echo "Waiting for $HOST:$PORT for up to $TIMEOUT seconds..."

for i in $(seq 1 $TIMEOUT); do
  nc -z "$HOST" "$PORT" > /dev/null 2>&1
  if [ $? -eq 0 ]; then
    echo "$HOST:$PORT is available!"
    if [ -n "$CMD" ]; then
      exec $CMD
    fi
    exit 0
  fi
  sleep 1
done

echo "Timeout after $TIMEOUT seconds waiting for $HOST:$PORT"
exit 1
