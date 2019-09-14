# blog-api

Application that manage articles and store it in a database.

## Deploying using docker
I have opt to include both postgres and the blog-api in one docker container for easier deployment. Not a good solution,
ideally you would use docker compose but right now i dont have the resources to do it.

#### How to run
build the dockerfile
docker build -t blogapi:1.0.0

## Local Deployment

#### Requirements for Local machine
- [go 1.13](https://golang.org/)
- [go dep](https://github.com/golang/dep)
- [postgres 10](https://www.postgresql.org/download/)

#### How To Run
Environment Variables

| Key | Description | Default|
|:---|:---|:---|
|BASE_URL| Base path of blog-api | localhost|
| PORT | Port of blog-api | 8080|
|POSTGRES_HOST| Hostname of postgres | localhost|
|POSTGRES_PORT| Port of postgres|5432|
|POSTGRES_DB_NAME| Database name in postgres|blog|
|POSTGRES_USER| Postgres username|postgres|
|POSTGRES_PASSWORD| Postgress password|abc123|
Installing Go dependency
```
dep ensure
```

Run
```
go run blog.go
```

#### Running unit test
```
go run ./...
```

#####Running unit test with test coverage
install coverage tool
```
go get golang.org/x/tools/cmd/cover
```
then run
```
go test -cover ./...
```
