# Backend for [interactive CV](https://github.com/Maxim-Ba/cv)


## Sqlc

Run in docker compose sqlc for generate go models
`docker-compose up sqlc`


## Required envs
``` 
  POSTGRES_HOST=localhost
  POSTGRES_PORT=5432
  POSTGRES_USER=postgres
  POSTGRES_PASSWORD=postgres
  POSTGRES_DB=postgres
  SERVER_ADDRESS=localhost:3333
  MIGRATION_PATH=migrations
  LOG_LEVEL=info
```


## Тестирование

### Используемые библиотеки

- **testing** — стандартная библиотека Go для написания тестов
- **testcontainers-go** — запуск PostgreSQL в Docker-контейнере для интеграционных тестов
- **stretchr/testify** — assertions и require для удобной проверки результатов

#### Интеграционные тесты (repository)

Тесты слоя repository работают с реальной базой данных PostgreSQL, запущенной в Docker-контейнере. Тестируются все CRUD-операции репозиториев: создание, получение, обновление, удаление (одиночное и списком), а также пагинация, сортировка и фильтрация.


### Запуск тестов

```bash
# Все тесты
go test ./...

# Интеграционные тесты repository (требует Docker)
go test -v ./internal/repository/...

# Юнит-тесты services
go test -v ./internal/services/...

```
