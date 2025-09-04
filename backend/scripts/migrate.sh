#!/bin/bash
# Database migration script

set -e

# Default values
DIRECTION="up"
STEPS=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--down)
            DIRECTION="down"
            shift
            ;;
        -s|--steps)
            STEPS="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  -d, --down     Rollback migrations"
            echo "  -s, --steps N  Number of migration steps"
            echo "  -h, --help     Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

echo "Running database migrations..."

if [ "$DIRECTION" = "down" ]; then
    if [ -n "$STEPS" ]; then
        echo "Rolling back $STEPS migration steps..."
        go run cmd/migrate/main.go down $STEPS
    else
        echo "Rolling back all migrations..."
        go run cmd/migrate/main.go down
    fi
else
    echo "Applying migrations..."
    go run cmd/migrate/main.go up
fi

echo "Migration complete!"
