#!/bin/bash
# Simple script to test the ascii-art JSON REST API directly without the browser

set -euo pipefail

BASE_URL="${API_BASE_URL:-http://localhost:${PORT:-8080}}"
CURL_OPTS=(--silent --show-error --max-time 5)

echo "Sending API payload: 'hello' with 'standard' banner"
curl "${CURL_OPTS[@]}" -X POST "${BASE_URL}/api/ascii-art" \
  -H "Content-Type: application/json" \
  -d '{"text": "hello", "banner": "standard"}'

echo -e "\n\nTesting handling of invalid characters over JSON API:"
curl "${CURL_OPTS[@]}" -X POST "${BASE_URL}/api/ascii-art" \
  -H "Content-Type: application/json" \
  -d '{"text": "¿qué tal?", "banner": "standard"}'
