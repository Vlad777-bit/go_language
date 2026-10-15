# Finance Project

Проект состоит из двух сервисов:

- `ledger` — бизнес-логика работы с транзакциями и бюджетами;
- `gateway` — REST API для доступа к функциям Ledger.

## Структура проекта

```text
project/
├── README.md
├── ledger/
│   ├── go.mod
│   ├── ledger.go
│   ├── ledger_test.go
│   └── cmd/
│       └── ledger/
│           └── main.go
└── gateway/
    ├── go.mod
    ├── cmd/
    │   └── gateway/
    │       └── main.go
    └── internal/
        └── api/
            ├── dto.go
            ├── handlers.go
            ├── middleware.go
            ├── server.go
            └── server_test.go
```

## Запуск Ledger

```bash
cd ledger
go run ./cmd/ledger
```

## Запуск Gateway

```bash
cd gateway
go run ./cmd/gateway
```

Gateway запускается на:

```text
http://localhost:8080
```

## API

### Создать транзакцию

```http
POST /api/transactions
```

Пример:

```bash
curl -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-d '{
  "amount": 450,
  "category": "еда",
  "description": "ланч",
  "date": "2026-10-15"
}'
```

Успешный ответ:

```json
{
  "amount": 450,
  "category": "еда",
  "description": "ланч",
  "date": "2026-10-15"
}
```

Статус:

```text
201 Created
```

### Получить транзакции

```http
GET /api/transactions
```

Пример:

```bash
curl http://localhost:8080/api/transactions
```

Если транзакций нет:

```json
[]
```

### Создать или обновить бюджет

```http
POST /api/budgets
```

Пример:

```bash
curl -X POST http://localhost:8080/api/budgets \
-H "Content-Type: application/json" \
-d '{
  "category": "еда",
  "limit": 5000
}'
```

Ответ:

```json
{
  "category": "еда",
  "limit": 5000
}
```

### Получить бюджеты

```http
GET /api/budgets
```

Пример:

```bash
curl http://localhost:8080/api/budgets
```

## Ошибки

Все ошибки возвращаются в едином формате:

```json
{
  "error": "описание ошибки"
}
```

Используемые HTTP-коды:

- `400 Bad Request` — ошибка валидации или некорректный JSON;
- `409 Conflict` — превышен бюджет;
- `500 Internal Server Error` — внутренняя ошибка.

Все ответы имеют заголовок:

```text
Content-Type: application/json; charset=utf-8
```

## Тесты

Ledger:

```bash
cd ledger
go test -v ./...
```

Gateway:

```bash
cd gateway
go test -v ./...
```
