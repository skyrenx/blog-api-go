#!/bin/bash

# Enable strict error handling to catch any issues early
set -e

# Function to print an error message and exit
error_exit() {
  echo "❌ Error: $1"
  exit 1
}

# Check if main.go exists in the current directory
if [ ! -f "main.go" ]; then
  error_exit "main.go not found. Ensure your Go entry file is named main.go and is located in this directory."
fi

# Check if template.yaml exists in the current directory
if [ ! -f "template.yaml" ]; then
  error_exit "template.yaml not found. This file is required for AWS SAM to simulate the Lambda environment."
fi

echo "✅ main.go and template.yaml found."

# Clean previous SAM build to avoid caching issues
echo "🧹 Cleaning previous SAM build..."
rm -rf .aws-sam/build

# Build the Go binary with SAM CLI
echo "🏗️ Building project with AWS SAM CLI..."
sam build

if [ $? -ne 0 ]; then
  error_exit "❌ SAM build failed."
fi
echo "✅ SAM build completed successfully."
echo "CLUSTER_ENDPOINT before SAM: $CLUSTER_ENDPOINT"
# Start AWS SAM local API server
echo "🚀 Starting AWS SAM local API gateway on http://localhost:3000..."
sam local start-api --env-vars env.json
