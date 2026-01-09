# Deployment

## Build Container Image

```bash
docker build -t ghcr.io/guidonguido/exposing-guidongui:latest -f ./deploy/blog/Dockerfile .
docker push ghcr.io/guidonguido/exposing-guidongui:latest
```

## Deploy to Kubernetes

```bash
# Deploy the blog
kubectl apply -f ./deploy/blog

# Deploy asciinema server
kubectl apply -f ./deploy/asciinema
```
