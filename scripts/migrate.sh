#!/bin/sh
set -e

cd "$(dirname "$0")/.."

# Load .env if present (so DATABASE_URL can be set there)
if [ -f .env ]; then
    set -a
    . ./.env
    set +a
fi

if [ -z "$DATABASE_URL" ]; then
    echo "DATABASE_URL is not set. Example: postgres://user:password@localhost:5432/dbname?sslmode=disable"
    echo "Set DATABASE_URL in the environment or add it to a .env file in the project root."
    exit 1
fi

# Find migrate CLI (go install puts it in $HOME/go/bin)
MIGRATE_CMD=""
if command -v migrate >/dev/null 2>&1; then
    MIGRATE_CMD=migrate
elif [ -x "$HOME/go/bin/migrate" ]; then
    MIGRATE_CMD="$HOME/go/bin/migrate"
fi
if [ -z "$MIGRATE_CMD" ]; then
    echo "migrate CLI not found. Install with:"
    echo "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    echo "Then add \$HOME/go/bin to PATH, or run: PATH=\"\$HOME/go/bin:\$PATH\" $0 $1"
    exit 1
fi

case "$1" in
    up)
        "$MIGRATE_CMD" -path migrations -database "$DATABASE_URL" up
        ;;
    down)
        "$MIGRATE_CMD" -path migrations -database "$DATABASE_URL" down
        ;;
    *)
        echo "Usage: $0 up|down"
        exit 1
        ;;
esac
