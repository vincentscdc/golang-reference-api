[Golang Reference Guide](https://github.com/monacohq/golang-reference-guide)

## CircleCI
Steps for preparing the project in CircleCI
1. Create Project
1. Add User Key (?)
1. Create CircleCI Context to provide parameters
    | Name |
    |---|
    | AWS_ACCESS_KEY_ID |
    | AWS_ECR_ACCESS_KEY_ID |
    | AWS_ECR_ACCOUNT_URL |
    | AWS_ECR_REGISTRY_ID |
    | AWS_ECR_SECRET_ACCESS_KEY |
    | AWS_REGION |
    | AWS_SECRET_ACCESS_KEY |
    | EKS_CLUSTER_NAME |
    | EKS_NAMESPACE |
1. Modify data in .circleci/config.yml
    1. deployment-name: name of Kubernetes deployment (set with infra team)
    1. context: CircleCI context set for different phases (set with infra team)
    1. image-name: name for ECR repository (in current example: app-crypto-pay-later-api)

## Folder structure and Kong gateway setting
### Folder structure
We have 2 sets of APIs:
* internal/port/rest/internalfacing: called by internal service, such as communicating with rails main app
* internal/port/rest/userfacing: called by external clients, such as mobile Apps
    * okstylewrapper.go: User facing APIs are called by exiting clients. We have follow the current response format. This wrapper will help to convert the format.
    * useruuidmiddleware.go: For user facing APIs, we use a plug-in, crypto-auth-user-auth-token, for checking bearer token on Kong gateway. If the token is ok, it will be transfer to user's UUID. User's UUID will be put in the header, "X-CRYPTO-USER-UUID".

### Kong gateway setting
You need to prepare the following items for requesting Kong gateway setting.
* Request domain for internal APIs
    * domain:
        - {service-name}.app-{env}.local
    * paths:
        - /v1/swagger/*
        - /v1/api/internal/*
* Request domain for user facing APIs
    * domain:
        - {env}-{service-name}.3ona.co
    * paths:
        - /v1/api/{api-path}/*
    * plug-in:
        - crypto-auth-user-auth-token

## requirements

```bash
go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
```

## Seeing traces

```bash
kubectl port-forward -n observability svc/jaeger-query 16686:16686
```

then go to 127.0.0.1:16686

## Swagger

For the api generation, you will need swag:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

For all generated code, you will need openapi-generator:

```bash
brew install openapi-generator
```

### Swagger template gen

```bash
make swagger-gen
```

Unfortunately, the generated json doesn't take into account the dynamic variable update in main.go, so the only valid swagger json definition is served by the go binary.


### Protocol Buffer generate
We rely on buf help us manage the protocol buffer generation and linter.

- protoc v3.19.4 (https://github.com/protocolbuffers/protobuf/releases
- brew install protobuf (mac only)
- brew install bufbuild/buf/buf

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

```bash
make proto-gen
```

### k6 Benchmarks

Run your app and go to your [swagger json](http://localhost:3000/v1/swagger/) to gen the benchmarks with the openapi generator.
