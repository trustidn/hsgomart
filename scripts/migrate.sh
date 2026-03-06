#!/bin/sh
set -e

cd "$(dirname "$0")/.."

if [ -z "$DATABASE_URL" ]; then
    echo "DATABASE_URL is not set. Example: postgres://user:password@localhost:5432/dbname?sslmode=disable"
    exit 1
fi

case "$1" in
    up)
        migrate -path migrations -database "$DATABASE_URL" up
        ;;
    down)
        migrate -path migrations -database "$DATABASE_URL" down
        ;;
    *)
        echo "Usage: $0 up|down"
        exit 1
        ;;
esac
