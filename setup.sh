#!/bin/bash

# This script sets up all the necessary Go files for the Vaxen API

echo "Setting up Vaxen API project structure..."

# Create all necessary directories
mkdir -p internal/{config,middleware,routes,handlers,models,database,utils}
mkdir -p cmd/migrate

echo "✅ Directories created"
echo "Project structure is ready!"
echo "Next steps:"
echo "1. Run: go mod tidy"
echo "2. Build: make build or go build -o bin/api main.go"
echo "3. Run: make run or go run main.go"

