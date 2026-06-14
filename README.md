# Event Tracking Platform

Сервис для сбора и обработки пользовательских событий.

## Стек

* Go
* PostgreSQL
* Kafka
* Docker
* Prometheus
* Grafana
* Swagger

## Архитектура

Client
↓
REST API
↓
Kafka Producer
↓
Kafka
↓
Processor Consumer
↓
PostgreSQL

Prometheus ← Processor
↓
Grafana

API принимает события от клиентов и публикует их в Kafka.

Processor читает сообщения из Kafka и сохраняет их в PostgreSQL.

Prometheus собирает метрики приложения.

Grafana визуализирует метрики.

## Запуск

```bash
docker compose up -d
```

Запуск API:

```bash
go run cmd/api/main.go
```

Запуск Processor:

```bash
go run cmd/processor/main.go
```

## Swagger

http://localhost:8080/swagger/index.html

![img.png](img.png)

## Метрики

http://localhost:2112/metrics

![img_1.png](img_1.png)

## Grafana

http://localhost:3000

![img_2.png](img_2.png)

## Основные возможности

* Создание событий через REST API
* Асинхронная обработка через Kafka
* Сохранение событий в PostgreSQL
* Мониторинг через Prometheus
* Дашборды Grafana
* Swagger документация
* Graceful Shutdown
* Unit Tests

```
```
