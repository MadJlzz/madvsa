# madvsa

`madvsa` is a vulnerability scanners orchestrator that aims to scan OCI artifacts in multiple ways. It can be periodic,
just in time, or even triggered based on supported events.

# Roadmap

- Find a way not to have to rebuild everytime docker images whenever I have to change something in the code.

- Create and use one abstraction layer to create scans between Docker and Kubernetes to avoid rewriting all
  the business logic over and over again

- Make sure we can let the user decide to mount a volume inside the pod for scans results storage
- Configure the pod we start to be able to send data over S3 like storage. (GCS implemented, need config of the pod)

- Handle cronjobs to periodically scan images
    - predefined in a list
    - fetched from an OCI registry

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
