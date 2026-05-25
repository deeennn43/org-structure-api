# API организационной структуры

REST API для управления подразделениями и сотрудниками (дерево департаментов + сотрудники).

Репозиторий: https://github.com/deeennn43/org-structure-api

## Стек

- Go `net/http`
- PostgreSQL + GORM
- Миграции [goose](https://github.com/pressly/goose)
- Docker / docker-compose

## Быстрый старт

```bash
docker compose up --build
```

API: `http://localhost:8080`

## Тестирование в GoLand

Откройте файл `api.http` в корне проекта и нажмите зелёную стрелку у нужного запроса (нужен запущенный `docker compose up`).

## Примеры запросов

```bash
# Создать корневое подразделение
curl -X POST http://localhost:8080/departments/ \
  -H "Content-Type: application/json" \
  -d '{"name":"Компания"}'

# Создать дочернее
curl -X POST http://localhost:8080/departments/ \
  -H "Content-Type: application/json" \
  -d '{"name":"IT","parent_id":1}'

# Создать сотрудника
curl -X POST http://localhost:8080/departments/2/employees/ \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Иван Иванов","position":"Backend Developer"}'

# Получить подразделение с деревом (depth=2)
curl "http://localhost:8080/departments/1?depth=2&include_employees=true"

# Переместить подразделение
curl -X PATCH http://localhost:8080/departments/2 \
  -H "Content-Type: application/json" \
  -d '{"parent_id":1}'

# Удалить каскадом
curl -X DELETE "http://localhost:8080/departments/2?mode=cascade"

# Удалить с переводом сотрудников
curl -X DELETE "http://localhost:8080/departments/2?mode=reassign&reassign_to_department_id=1"
```

## Локальная разработка (GoLand)

1. Поднять только БД: `docker compose up db`
2. Переменные окружения (Run Configuration → Environment):
   - `DB_HOST=localhost`
   - `DB_PORT=5432`
   - `DB_USER=postgres`
   - `DB_PASSWORD=postgres`
   - `DB_NAME=org_structure`
3. Запуск: `cmd/api/main.go`

## Тесты

```bash
go test ./...
```

## Структура проекта

```
cmd/api/           — точка входа (main)
internal/domain/   — сущности предметной области
internal/repository/ — интерфейсы хранилища
internal/repository/gorm/ — реализация GORM
internal/service/  — бизнес-логика
internal/handler/  — HTTP-слой
internal/validation/
migrations/        — SQL-миграции goose
```

Слои разделены по принципам SOLID: HTTP не знает про SQL, сервис не знает про `net/http`.
