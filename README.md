# madvsa

`madvsa` is a vulnerability scanners orchestrator that aims to scan OCI artifacts in multiple ways. It can be periodic,
just in time, or even triggered based on supported events.

# Developing

I use Golang to develop all components, whether it's the controlplane or the wrapper used to run the scanners.

To get started, start the controlplane
```go
cd controlplane && go run *.go
```

Build the scanner's image
```shell
docker build -t madvsa/grype:latest -f scanner/Dockerfile.grype scanner
docker build -t madvsa/trivy:latest -f scanner/Dockerfile.trivy scanner
```

Trigger scans using the REST API
```shell
# This will start a new image scan of alpine:3.17 using Trivy and Docker/Podman as orchestrator of work.
curl -XPOST 'http://localhost:3000/api/v1/scanner/trivy/trigger?img=alpine:3.17'
```

If you prefer to use Kubernetes like I do, you can always spin up a new cluster with `kind`.
```shell
kind create cluster --name madvsa
```

Now deploy everything inside that local cluster
```shell
kubectl apply -f deploy/manifests
```
