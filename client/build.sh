#!/bin/bash

echo "Now building..."

# Remove dist directory if it exists
if [ -d "dist" ]; then
    echo "Removing existing dist directory..."
    rm -rf dist
fi

# Create dist directory
mkdir dist

# Build with anti-virus evasion flags
BUILD_TIME=$(date '+%Y%m%d_%H%M%S')
go build -ldflags="-H windowsgui -s -w -X main.buildTime=${BUILD_TIME}" -trimpath -o dist/rs-shortcuts-monitor-client.exe ./src

# Copy required files to dist directory
if [ -f "settings.ini" ]; then
    cp settings.ini dist/
    echo "Copied settings.ini to dist/"
fi

if [ -f "keys.json" ]; then
    cp keys.json dist/
    echo "Copied keys.json to dist/"
fi

if [ -f "dist/rs-shortcuts-monitor-client.exe" ]; then
    echo "Build successful!"
    echo "Files created in dist directory:"
    ls -la dist/
else
    echo "Build failed."
    exit 1
fi
