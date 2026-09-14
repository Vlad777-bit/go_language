# Разработка на Go

Репозиторий предназначен для хранения домашних заданий по дисциплине «Разработка на Go».

## Информация

**Преподаватель:** Андрей Пешков  
**Форма проверки:** самопроверка  

Репозиторий содержит два независимых Go-сервиса:

- `gateway` — HTTP-шлюз;
- `ledger` — сервис бизнес-логики для транзакций.

## Требования

- Go 1.22 или новее.

## Запуск Gateway

```bash
cd gateway
go run .
```

По умолчанию Gateway запускается на порту `8080`.

Проверка:

```bash
curl -i http://localhost:8080/ping
```

Ожидаемый ответ:

```text
HTTP/1.1 200 OK
...

pong
```

## Запуск Ledger

```bash
cd ledger
go run .
```

Ledger создаёт несколько тестовых транзакций в памяти и выводит их в консоль.

Пример вывода:

```text
Ledger service started
Transactions:
ID: 1, Amount: 25.50, Category: Food, Description: Lunch, Date: 2026-09-14
ID: 2, Amount: 60.00, Category: Transport, Description: Train ticket, Date: 2026-09-14
ID: 3, Amount: 120.99, Category: Shopping, Description: Headphones, Date: 2026-09-14
```

