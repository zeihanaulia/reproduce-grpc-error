# reproduce-grpc-error

## Build grpc server

```
cd grpc-server
docker build -t rge-grpc-server -f ./deployments/Dockerfile.grpc .
docker tag rge-grpc-server zeihanaulia/rge-grpc-server:1.0.0
```

## Build http server

```
cd http-server
docker build -t rge-http-server -f ./deployments/Dockerfile.http .
docker tag rge-http-server zeihanaulia/rge-http-server:1.0.0
```