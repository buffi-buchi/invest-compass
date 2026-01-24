#!/bin/bash

# Global variables
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

export MANIFESTS_DIR="$PROJECT_ROOT/manifests"
DEPLOYMENT_DIR="$PROJECT_ROOT/backend/deployment"

# 1. Create manifests directory
echo "Create directory for manifests: $MANIFESTS_DIR"

rm -rf "$MANIFESTS_DIR"
mkdir -p "$MANIFESTS_DIR"

# 1. Download Gateway API CRDs
GATEWAY_API_URL="https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml"

echo "Downloading Gateway API CRDs..."
curl -L "$GATEWAY_API_URL" -o "$MANIFESTS_DIR/00-gateway-crd.yaml"

# 2. Create a k3d cluster
echo "Creating k3d cluster..."
k3d cluster create --config "$DEPLOYMENT_DIR/cluster/k3d-config.yaml"
