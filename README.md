# CV Backend API

[![CI](https://github.com/Maxim-Ba/cv-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/Maxim-Ba/cv-backend/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21-blue)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

REST API для [интерактивного CV-сайта](https://github.com/Maxim-Ba/cv). Предоставляет данные о тегах, технологиях, истории работы и образовании.

## Архитектура

```
Angular SSR  →  nginx  →  Go API  →  PostgreSQL
                              ↑
                         Admin Panel (templ SSR)
```

**Слоёная архитектура:**
```
router (HTTP handlers)
  └── services (бизнес-логика, валидация)
        └── repository (SQL-запросы, pgx)
              └── PostgreSQL
```

## Стек технологий

| Компонент | Технология |
|-----------|-----------|
| Язык | Go 1.21 |
| HTTP Router | go-chi/chi v5 |
| База данных | PostgreSQL 16 + pgx v5 |
| SQL генерация | sqlc |
| Миграции | golang-migrate |
| Шаблоны (Admin) | a-h/templ |
| Документация API | Swagger (swaggo) |
| Логирование | log/slog |
| Тесты | testing + testcontainers-go + testify |
| Контейнеризация | Docker + docker-compose |
| CI | GitHub Actions |

## Возможности

- **REST API**: теги, технологии, история работы, образование — полный CRUD
- **Admin-панель** (`/admin/`): SSR на templ, CSRF-защита, HMAC cookie-аутентификация
- **Swagger UI** (`/swagger/`) — живая документация
- **Health check** (`/healthz`) — статус сервера и БД
- **Graceful shutdown** — корректное завершение при SIGTERM/SIGINT
- **Пагинация, сортировка, фильтрация** во всех list-эндпоинтах
- **Интеграционные тесты** repository-слоя с реальной PostgreSQL в Docker

## Быстрый старт

```bash
# Запуск с docker-compose (app + postgres)
docker-compose up app postgres

# Только база данных
docker-compose up postgres

# Локально (требует .env файл)
cp .env.example .env
go run ./cmd/main.go
```

Swagger UI: http://localhost:3333/swagger/
Admin: http://localhost:3333/admin/

## Переменные окружения

| Переменная | Описание | По умолчанию |
|-----------|----------|-------------|
| `POSTGRES_HOST` | Хост PostgreSQL | — |
| `POSTGRES_PORT` | Порт PostgreSQL | `5432` |
| `POSTGRES_USER` | Пользователь БД | — |
| `POSTGRES_PASSWORD` | Пароль БД | — |
| `POSTGRES_DB` | Имя базы данных | — |
| `SERVER_ADDRESS` | Адрес сервера | `localhost:3333` |
| `MIGRATION_PATH` | Путь к миграциям | `migrations` |
| `LOG_LEVEL` | Уровень логирования | `error` |
| `APP_ENV` | Среда окружения | `development` |
| `ALLOWED_ORIGIN` | Разрешённый CORS origin | `http://localhost:4200` |
| `ADMIN_USER` | Логин admin-панели | `admin` |
| `ADMIN_PASSWORD` | Пароль admin-панели | — |
| `APP_SECRET` | HMAC/CSRF секрет (≥32 символа) | — |

Пример конфигурации: `.env.example`

## Тестирование

### Библиотеки

- **testing** — стандартная библиотека Go
- **testcontainers-go** — запуск PostgreSQL в Docker для интеграционных тестов
- **stretchr/testify** — assertions/require

### Запуск тестов

```bash
# Все тесты
go test ./...

# Юнит-тесты сервисов (без Docker)
go test -v ./internal/services/...

# Интеграционные тесты репозиториев (требует Docker)
go test -v ./internal/repository/...

# С покрытием
go test -cover ./...
```

## Генерация кода

```bash
# Regenerate SQL models (sqlc)
docker-compose up sqlc

# Regenerate Swagger docs
swag init -g cmd/main.go -o docs --parseDependency
```
