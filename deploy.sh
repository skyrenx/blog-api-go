#!/bin/bash
# Exit immediately if a command exits with a non-zero status.
set -e

# Ensure that you have the AWS CLI configured and Go installed.
# This script assumes your Lambda code is in main.go in the current directory.

# Set environment variables for the target build
export GOOS=linux
export GOARCH=arm64

# Define output binary name and deployment package name
BINARY_NAME=bootstrap
ZIP_FILE=deployment.zip

echo "Building Go binary for Linux ARM64..."
go build -o $BINARY_NAME main.go

echo "Packaging the binary into $ZIP_FILE..."
# -j flag ensures that the zip does not include directory structure
zip -j $ZIP_FILE $BINARY_NAME

echo "Deploying the package to AWS Lambda function 'blogApi'..."
aws lambda update-function-code --function-name blogApi --zip-file fileb://$ZIP_FILE

echo "Deployment complete!"
