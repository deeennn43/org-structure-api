# org-structure-api

REST API для подразделений и сотрудников.

## Запуск

```bash
docker compose up --build
```

Сервис: http://localhost:8080

## Тесты

```bash
go test ./...
```

## Примеры

```bash
curl -X POST http://localhost:8080/departments/ \
  -H "Content-Type: application/json" \
  -d '{"name":"Компания"}'

curl -X POST http://localhost:8080/departments/2/employees/ \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Иван Иванов","position":"Developer"}'

curl "http://localhost:8080/departments/1?depth=2&include_employees=true"

curl -X PATCH http://localhost:8080/departments/2 \
  -H "Content-Type: application/json" \
  -d '{"parent_id":1}'

curl -X DELETE "http://localhost:8080/departments/2?mode=cascade"
```

Дополнительные запросы - в `api.http`.

## Локально без Docker для API

```bash
docker compose up db
```

Переменные: `DB_HOST=localhost`, `DB_PORT=5432`, `DB_USER=postgres`, `DB_PASSWORD=postgres`, `DB_NAME=org_structure`.

Запуск: `go run ./cmd/api`

## Структура

- `cmd/api` - main
- `internal/handler` - HTTP
- `internal/service` - бизнес-логика
- `internal/repository` - БД
- `migrations` - goose
