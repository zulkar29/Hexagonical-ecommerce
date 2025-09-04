#!/bin/bash
# Development setup script

set -e

echo "Setting up development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+')
echo "Go version: $GO_VERSION"

# Install dependencies
echo "Installing Go dependencies..."
go mod download
go mod tidy

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "Please update .env file with your local configuration"
fi

# Create necessary directories
mkdir -p logs
mkdir -p uploads

echo "Development environment setup complete!"
echo "Run 'make dev' to start the development server"
