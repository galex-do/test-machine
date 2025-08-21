#!/bin/bash

# Build script for Test Management Platform

echo "Building Test Management Platform..."

# Build the Go application
echo "Building Go binary..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main simple_main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful! Binary created: ./main"
    echo ""
    echo "To run the application:"
    echo "  ./main"
    echo ""
    echo "Or with Docker:"
    echo "  docker-compose up --build"
else
    echo "❌ Build failed!"
    exit 1
fi