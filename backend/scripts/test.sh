#!/bin/bash
# Test runner script

set -e

echo "Running tests..."

# Parse command line arguments
COVERAGE=false
INTEGRATION=false
VERBOSE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -i|--integration)
            INTEGRATION=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  -c, --coverage     Generate coverage report"
            echo "  -i, --integration  Run integration tests"
            echo "  -v, --verbose      Verbose output"
            echo "  -h, --help         Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Load test environment variables
if [ -f .env.test ]; then
    export $(cat .env.test | xargs)
fi

# Set test flags
TEST_FLAGS=""
if [ "$VERBOSE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -v"
fi

if [ "$COVERAGE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -cover -coverprofile=coverage.out"
fi

# Run tests
if [ "$INTEGRATION" = true ]; then
    echo "Running integration tests..."
    go test $TEST_FLAGS ./tests/...
else
    echo "Running unit tests..."
    go test $TEST_FLAGS ./internal/...
fi

# Generate coverage report if requested
if [ "$COVERAGE" = true ]; then
    echo "Generating coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated: coverage.html"
fi

echo "Tests completed!"
