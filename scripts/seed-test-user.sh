#!/bin/sh
# Register a test user via POST /auth/register.
# Run with backend up: ./scripts/seed-test-user.sh
# Then login with email + password below.

cd "$(dirname "$0")/.."

if [ -f .env ]; then
  set -a
  . ./.env
  set +a
fi

PORT="${APP_PORT:-8080}"
URL="http://localhost:${PORT}/auth/register"

EMAIL="${TEST_USER_EMAIL:-test@example.com}"
PASSWORD="${TEST_USER_PASSWORD:-password123}"
NAME="${TEST_USER_NAME:-Test Owner}"

echo "Registering test user: $EMAIL"
echo "URL: $URL"
echo ""

# Check if server is up
if ! curl -s -o /dev/null -w "%{http_code}" "http://localhost:${PORT}/health" | grep -q 200; then
  echo "Backend not responding at http://localhost:${PORT}/health"
  echo "Start the server first: go run cmd/server/main.go"
  exit 1
fi

RESP=$(curl -s -w "\n%{http_code}" -X POST "$URL" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"$NAME\",\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

HTTP_CODE=$(echo "$RESP" | tail -n1)
BODY=$(echo "$RESP" | sed '$d')

if [ "$HTTP_CODE" = "201" ]; then
  echo "OK (201). Test user created."
  echo "Login with: email=$EMAIL password=$PASSWORD"
else
  echo "HTTP $HTTP_CODE"
  echo "$BODY"
  exit 1
fi
