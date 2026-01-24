#!/bin/bash

# Global variables
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

export MANIFESTS_DIR="$PROJECT_ROOT/manifests"
DEPLOYMENT_DIR="$PROJECT_ROOT/backend/deployment"

# 1. Create a k3d cluster
echo "Deleting k3d cluster..."
k3d cluster delete --config "$DEPLOYMENT_DIR/cluster/k3d-config.yaml"
