#!/bin/sh
set -e

echo "==> Starting ragforge..."
exec ./ragforge -config /app/config/config.yaml
