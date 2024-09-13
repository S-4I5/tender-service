# Tender Service

## 1. Description

Сервис для создания тендеров и работы с предложениями.

Stack: Golang / PostgreSQL / Goose (Migrations) / Swagger / Testcontainers (Интеграционные тесты)

P.S Надеюсь, что везде попал с форматами в openapi, а то пару раз приходилось всё менять =D. В некоторых местах (/my эндпоинты) сделал username required полем, ибо не смог придумать, что делать в ситуации, когда его не передают, а доке не написано.

## 2. ENVs

| Name           | Type     | Default value     | Description      |
|----------------|----------|-------------------|------------------|
| POSTGRES_CONN  | String   |                   | Psql conn string |
| SERVER_ADDRESS | String   | :8080             | Servet address   |
| MIGRATIONS_DIR | String   | ./migrations/prod | Migrations dir   |

## 3. Run
```
go run ./cmd/main.go
```

## 4. Deployment
```
docker-compose up
```

## 5. Swagger
```
http://localhost:8080/swagger/index.html#/
```

## 6. Tests

### 6.1 Integrational tests with Testcontainers
```
make run-it
```
