#!/usr/bin/env bash
set -euo pipefail

REPO_OWNER="$(echo "${GITHUB_REPOSITORY_OWNER}" | tr '[:upper:]' '[:lower:]')"
IMAGE_NAME="net-sentinel"

echo "🚀 Deploying ${REPO_OWNER}/${IMAGE_NAME}"

docker compose pull
docker compose up -d --build
