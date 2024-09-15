# Tender Service

## 1. Description

Сервис для создания тендеров и работы с предложениями.

Stack: Golang / PostgreSQL / Goose (Migrations) / Swagger / Testcontainers (Интеграционные тесты)

P.S Надеюсь, что везде попал с форматами в openapi, а то пару раз приходилось много менять =D. В некоторых местах (/my эндпоинты) сделал username required полем, ибо не смог придумать, что делать в ситуации, когда его не передают, а доке не написано.

## 2. ENVs

| Name           | Type     | Default value     | Description      |
|----------------|----------|-------------------|------------------|
| POSTGRES_CONN  | String   |                   | Psql conn string |
| SERVER_ADDRESS | String   | :8080             | Servet address   |
| MIGRATIONS_DIR | String   | ./migrations/prod | Migrations dir   |

## 3. How to run

Хотелось бы тут чуть подробнее рассказать про разные "профили" миграций. Все они находятся в ./migrations и лежат в разных директориях:
* prod : не содержит запросов для создания таблиц юзеров, организаций и таблицы для их связи. Была сделана для раскататки на стенд с (Вариант по умолчанию вариант)
* test/clear : содержит все необходимые таблицы для работы с "пустой" дб. Используется в it тестах и для локального запсука. (Используется в тестах)
* test/users : содержит всё то же, что и clear, но и несколько юзеров и групп для локальной отладки. (Используется в приложенном docker-compose)

### 3.1 Via go
```
go run ./cmd/main.go
```

### 3.2 Via docker
```
docker build -t tender-service .
docker run -p 8080:8080 --name tender-service_container -e {some envs} tender-service
```

### 3.3 Via docker compose (With PostgreSQL container)
```
docker-compose up
```
или
```
make run-compose
```
или (с пересборкой контейнера)
```
run-compose-b
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
