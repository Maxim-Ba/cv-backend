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
