#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/../"

curl --location --request POST 'localhost:1323/admin-console' \
--header 'Content-Type: application/json' \
--data-raw '{
    "event": {
        "title": "customTitle",
        "environment": "production",
        "contexts": {"Research Source": {"research_source": "kacher loh"}}
    },
    "url": "customUrl",
    "culprit": "customCulprit",
    "project_slug": "customProjectSlug"
}'