#!/bin/bash

REGISTRY="${REGISTRY:-localhost:5000}"

echo "Logging into private registry..."
echo "$REGISTRY_PASSWORD" | docker login "$REGISTRY" -u "$REGISTRY_USERNAME" --password-stdin

echo "Building backend-go..."
docker build -t "$REGISTRY/backend-go" ./backend-go

echo "Building frontend..."
docker build -t "$REGISTRY/frontend" ./frontend

echo "Pushing images to registry..."
docker push "$REGISTRY/backend-go"
docker push "$REGISTRY/frontend"

echo "Done."
