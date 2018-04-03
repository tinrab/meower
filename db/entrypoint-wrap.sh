#!/bin/bash

until cat "./create.cql" | cqlsh; do
  echo "scylla is unavailable - retry later"
  sleep 2
done &

exec /docker-entrypoint.py "$@"
