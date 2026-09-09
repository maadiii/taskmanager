#!/bin/sh
set -eu

exec /migrate \
  -path=/migrations \
  -database="${PG_DSN:?PG_DSN is required}" \
  up