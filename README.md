# Micro Service Tools

## Usage
Please check `example` dir.

1. start etcd service at localhost:2379
2. in `server` dir, run `go run main.go`
3. in `client` dir, run `go run main.go`

## Transfer HTTP to grpc
Please check `example/http-server` dir.

1. start etcd service at localhost:2379
2. in `server` dir, run `go run main.go`
3. in `http-server` dir, run `go run main.go`
4. in `http-client` dir, run `go run main.go`

### HTTP Call Method
```
POST /grpc/{service_name}/{method_name}
body is proto.Request for each method in json format
```

## Roadmap 
- service discover
    - [x] register
    - [x] discover
    - [x] resovler
    - [x] watcher
    - [x] auto-gen service id

- gateway
    - [x] connection container
    - [x] balancer
    - [x] inject trace id
    - [x] http to grpc transfer
    - [x] jwt generate / verify
    - [ ] auth

- middleware
    - [x] logging
    - [x] recover

- tool
    - [ ] list all service

- deploy
    - [ ] docker example