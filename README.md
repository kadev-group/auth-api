# auth-api

Go part of kantoo services.

## Deployment

```sh
docker build -f dockerfile -t doxanocap/auth-api:prod .
docker push doxanocap/auth-api:prod
```

### Operations

- create new migration:

```sh
migrate create -ext sql -dir api/migrations mg_name
migrate -path api/migrations -database "postgres://postgres:tdepassword@localhost:5432/auth_api?sslmode=disable" up
```

- run Postgres:
```shell
docker run --name auth-api-pq -p 5222:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=kantoo12345 -e POSTGRES_DB=auth_api -d postgres:14-alpine
```

- run Redis:
```shell
docker run --name=redis -p 6379:6379 --restart=always -d redis:latest
```

- run RabbitMQ:
```shell
docker run -d --hostname rabbit-mq --name rabbit-mq -p 15672:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password rabbitmq:3-management
```
