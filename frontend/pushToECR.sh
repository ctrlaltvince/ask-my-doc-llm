#!/bin/bash

set -e

# Set once
REGION=us-west-1
ACCOUNT_ID=242650469816
REPO_NAME=ask-my-doc-frontend

# Authenticate
aws ecr get-login-password --region "$REGION" \
  | docker login --username AWS \
  --password-stdin "${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

# Create ECR repo (if not already done)
aws ecr create-repository \
  --repository-name "$REPO_NAME" \
  --region "$REGION" 2>/dev/null || echo "ECR repo $REPO_NAME may already exist."

# Delete all images in the repo if any exist
IMAGES=$(aws ecr list-images --repository-name "$REPO_NAME" --query 'imageIds' --output json)
if [[ "$IMAGES" != "[]" ]]; then
  aws ecr batch-delete-image \
    --repository-name "$REPO_NAME" \
    --image-ids "$IMAGES"
else
  echo "No images to delete in $REPO_NAME."
fi

# Build and push image
docker build -t "$REPO_NAME" .
docker tag "$REPO_NAME:latest" "${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/$REPO_NAME:latest"
docker push "${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/$REPO_NAME:latest"
