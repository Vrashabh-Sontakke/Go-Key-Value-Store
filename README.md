# kvstore
 A simple In-Memory Key Value Store HTTP API Service with Golang.

## How to Run :

#### Requirements: 
Depending on the option you choose below; Go binary, Docker or Kind must be installed on the machine.

### Option - 1 : Local (no-docker) : 

from the `app` directory, execute `go run main.go`
### Option - 2 : Using Docker :

from the `app` directory, run ...
```
docker build -t kvstore .
```
```
docker run -d -p 8080:8080 kvstore
```

### Option - 3 : Kind (Local Kubernetes Cluster)

First build a docker image using Dockerfile from the `app` directory. Remember, give the Docker Image a specific tag, for example `kvstore:0.1` in this case.

```
docker build -t kvstore:0.1 .
```
Create kind cluster
```
kind create cluster --config cluster.yml
```
Load the Docker Image into Kind
```
kind load docker-image kvstore:0.1
```

