# invest-compass

## Install

To deploy application install Docker, k3s, kubectl and Helm.

```shell
# Install Docker using official instruction.
# https://www.docker.com/products/docker-desktop/

brew install k3d
brew install kubectl
brew install helm

# Create cluster
k3d cluster create --config backend/deployment/cluster/k3d-config.yaml

# Install custom traefik
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml
kubectl apply -f backend/deployment/cluster/traefik-config.yaml
kubectl apply -f backend/deployment/cluster/gateway.yaml

# To access the Traefik dashboard, use port-forwarding:
kubectl port-forward -n kube-system service/traefik 9001:8080
# Dashboard will be available at http://localhost:9001/dashboard/
```

## Test

```shell
helm template invest-compass . --debug
```
