Event Tracking Platform

Сервис для сбора, обработки и хранения пользовательских событий.

Возможности
Прием событий через REST API
Отправка событий в Kafka
Асинхронная обработка событий через Consumer
Сохранение событий в PostgreSQL
Валидация входящих данных
Swagger документация
Unit-тесты
Graceful Shutdown
Docker Compose для локального запуска
Архитектура
Client
|
v
REST API
|
v
Kafka Topic (events)
|
v
Processor
|
v
PostgreSQL
Стек технологий
Go 1.26
PostgreSQL
Apache Kafka
Docker Compose
Swagger
GitHub Actions
slog
Структура проекта
cmd/
├── api/
└── processor/

internal/
├── config/
├── handler/
├── kafka/
├── migrations/
├── models/
├── postgres/
├── repository/
└── service/

docs/
Запуск проекта
Запуск инфраструктуры
docker compose up -d
Запуск API
go run cmd/api/main.go
Запуск Processor
go run cmd/processor/main.go
Swagger

После запуска API:

http://localhost:8080/swagger/index.html

Пример создания события

POST /events

{
"user_id": "123",
"event_type": "purchase",
"page": "checkout",
"amount": 100
}
Получение событий

GET /events

Тестирование
go test ./...

кто сделал?

Андрей