# Go Key Value Store
 A simple In-Memory Key Value Store HTTP API Service with Golang.

## How to Run :

#### Requirements: 
Depending on the option you choose below; `Go binary`, `Docker` or `Kind` must be installed on the machine.

#### Example http API test usage:

`curl -X POST -d "key=abc-1" -d "value=value1" http://localhost:8080/set`

`curl "http://localhost:8080/get/abc-1"`

`curl "http://localhost:8080/search?prefix=a"`

`curl "http://localhost:8080/get/search?suffix=1"`

Prometheus Metrics: `curl "http://localhost:8080/metrics"`

### Option - 1 : Local (no-docker) : 

##### switch to `/app` directory & run >>>
```
go run main.go
```
### Option - 2 : Using Docker :

##### switch to `/app` directory & run ...
```
docker build -t kvstore .
```
```
docker run -d -p 8080:8080 kvstore
```

### Option - 3 : Kind (Local Kubernetes Cluster)

#### Build a Docker Image : 
- Remember ! Give the Docker Image a specific tag, for example `kvstore:0.1` in this case (required to work with kind).
- alternatively, you can Push the Docker Image to a Registry (Online/Offline), just make make sure to update it in `/k8sManifests/kvstore.yml` file.
##### switch to `/app` directory & run >>>
```
docker build -t kvstore:0.1 .
```
#### Create kind cluster
```
kind create cluster
```
#### Load the Docker Image into Kind
```
kind load docker-image kvstore:0.1
```
- if Docker Image Tag is other than `kvstore:0.1`, make sure to update it in `/k8sManifests/kvstore.yml` file.
  
#### Apply Kubernetes Configuration ...


##### switch to `/k8sManifests` directory & run >>>
```
kubectl apply -f kvstore.yml
```
##### Check for pods creation : 
```
kubectl get pods
```

##### Expose Service :
```
kubectl port-forward svc/kvstore-service 8080
```

## Observability :

##### switch to `/k8sManifests` directory & run >>>
```
helm install -f ./prometheus-values.yml prometheus prometheus-community/prometheus
```
##### Expose Prometheus Server : 
```
kubectl port-forward svc/prometheus-server 9090:80
```

##### Endpoint : http://localhost:9090

## Testing :

##### switch to `/app` directory & run >>>
```
go test
```
#####Check Test Coverage : 
```
go test -cover


```
