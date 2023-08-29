#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/../"

docker compose exec -T workspace go build -v -o bin/sms github.com/FscoreLab/sentry-mattermost-sidecar
