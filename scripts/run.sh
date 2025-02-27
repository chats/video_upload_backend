#!/bin/bash
# run.sh - Run script for the video transcoding system
set -e

echo "Starting Video Transcoding System..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found. Please run './setup.sh' first."
    exit 1
fi

# Source environment variables
export $(grep -v '^#' .env | xargs)

# Run the application
go run cmd/api/main.go

# The application should catch SIGTERM and exit gracefully,
# but just in case, handle Ctrl+C
trap "echo 'Stopping application...'; exit 0" INT TERM