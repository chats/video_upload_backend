#!/bin/bash
# setup.sh - Setup script for the video transcoding system
set -e

echo "Setting up Video Transcoding System..."

# Create necessary directories
mkdir -p tmp/uploads
mkdir -p logs

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go 1.16 or later."
    exit 1
fi

# Check if FFmpeg is installed
if ! command -v ffmpeg &> /dev/null; then
    echo "FFmpeg is not installed. Please install FFmpeg."
    exit 1
fi

if ! command -v ffprobe &> /dev/null; then
    echo "FFprobe is not installed. Please install FFprobe (usually comes with FFmpeg)."
    exit 1
fi

# Create .env file from example if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "Please update the .env file with your configuration."
fi

# Install Go dependencies
echo "Installing Go dependencies..."
go mod tidy

# Run database migrations
echo "Setting up database..."
go run cmd/migrate/main.go

echo "Setup completed successfully!"
echo "You can now run the application with './run.sh' or 'make run'"