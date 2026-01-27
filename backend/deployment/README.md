# Deployment Guide

This directory contains Helm charts and scripts for deploying the Invest Compass backend to a Kubernetes cluster.

## Prerequisites

Before deploying the application, ensure you have the following tools installed and configured:

- [Docker](https://www.docker.com/products/docker-desktop/) - To build and run containers.
- [k3d](https://k3d.io/) - To create a local Kubernetes cluster.
- [kubectl](https://kubernetes.io/docs/tasks/tools/) - To interact with the Kubernetes cluster.
- [Helm](https://helm.sh/docs/intro/install/) - To deploy applications using Helm charts.

## 1. Local Cluster Setup

A helper script is provided to create a local k3d cluster with the necessary configurations (Gateway API CRDs and Traefik).

```bash
# From the project root
./backend/deployment/create-cluster.sh
```

This script will:
1. Create a `manifests` directory in the project root.
2. Download Gateway API CRDs.
3. Create a k3d cluster named `invest-compass` (as configured in `backend/deployment/cluster/k3d-config.yaml`).

## 2. Deploy Applications

The applications are deployed using Helm charts located in the `backend/deployment` directory.

### Step 2.1: Deploy PostgreSQL

Deploy the PostgreSQL database first, as the server application depends on it.

```bash
# Navigate to the postgres deployment directory
cd backend/deployment/postgres

# Install the postgres chart in the 'postgres' namespace
helm install postgres . -n postgres --create-namespace
```

Wait for the PostgreSQL pod to be ready:
```bash
kubectl get pods -l app.kubernetes.io/name=postgres -n postgres
```

### Step 2.2: Build the Server Image

Before deploying the server, you need to build the Docker image and import it into the k3d cluster (or push it to a registry).

```bash
# From the project root
docker build -t invest-compass-server:1.0.0 -f backend/build/Dockerfile --build-arg APP_PATH=./cmd/server backend

# Import the image to k3d
k3d image import invest-compass-server:1.0.0 -c invest-compass
```

### Step 2.3: Deploy the Server Application

Now, deploy the Go server application.

```bash
# Navigate to the server deployment directory
cd backend/deployment/server

# Install the server chart in the 'server' namespace
helm install server . -n server --create-namespace
```

Verify that the server is running:
```bash
kubectl get pods -l app.kubernetes.io/name=server -n server
```

## 3. Accessing the Application

The server application is exposed via a `LoadBalancer` service and a Gateway `HTTPRoute`. Since we are using k3d, the ports are mapped to your localhost as defined in the `k3d-config.yaml`.

- **Main API (via Gateway)**: [http://localhost:9000](http://localhost:9000) (requests are proxied by Traefik through the `web` entrypoint)
- **Main API (direct Service)**: [http://localhost:9000](http://localhost:9000) (mapped to cluster load balancer port 80)

### Health Checks
- Readiness: `kubectl port-forward svc/server 8084:84 -n server` then `curl http://localhost:8084/readyz`
- Liveness: `kubectl port-forward svc/server 8084:84 -n server` then `curl http://localhost:8084/livez`

## Configuration

Custom configurations can be applied by modifying the `values.yaml` files in each chart's directory or by passing `--set` flags during `helm install`.

- **PostgreSQL**: `backend/deployment/postgres/values.yaml`
- **Server**: `backend/deployment/server/values.yaml`
