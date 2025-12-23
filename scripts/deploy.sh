#!/usr/bin/env bash
set -euo pipefail

echo "🚀 Deploying ${GITHUB_REPOSITORY}"

docker compose pull
docker compose up -d --build

