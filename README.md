# reproduce-grpc-error

## Build grpc server

```
docker build -t rge-grpc-server -f ./deployments/Dockerfile.grpc .
docker tag rge-grpc-server zeihanaulia/rge-grpc-server:1.0.0
```

## Build http server

```
docker build -t rge-http-server -f ./deployments/Dockerfile.http .
docker tag rge-http-server zeihanaulia/rge-http-server:1.0.0
```

## Deploy service

```
 kubectl apply -f .\deployments\grpc-server.yml
 kubectl apply -f .\deployments\http-server.yml

 -- RESTART
 kubectl rollout restart deployments/rge-grpc-server
 kubectl rollout restart deployments/rge-http-server

 -- Check logs
 kubectl logs -f deployments/rge-grpc-server
 kubectl logs -f deployments/rge-http-server

 -- delete service
kubectl delete deployment rge-grpc-server
kubectl delete deployment rge-http-server

-- check
kubectl get pods
kubectl get service
kubectl get deployments
```

## Hit service

Using [hey](https://github.com/rakyll/hey) for hit concurrent request.

hey -z 10m http://localhost:31533/connect

## Learn

Still, I cannot reproduce err `cannot assign requested address`.
If I restart the grpc service, the http still reconnect after service up.
and get error `could not greet: rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp 10.98.96.166:50051: connect: connection refused"`