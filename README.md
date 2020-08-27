# k8s-utils-backend-golang
Golang Backend for Kubernetes ui utils

## Using a private registry

```shell
docker build -t golang-backend .
docker tag golang-backend docker.registry.pr.ch/golang-backend:1.0.0
docker push docker.registry.pr.ch/golang-backend:1.0.0
```

## Deploy on docker with private registry

```shell
docker pull docker.registry.pr.ch/golang-backend:1.0.0
docker run -d --name go-backend -p 80:80 -v /root/.kube/config:/root/.kube/config docker.registry.pr.ch/golang-backend:1.0.0
```
