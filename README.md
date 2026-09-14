# Users API

REST API для управления пользователями (Go + Gin + SQLite).

## Требования

Docker и Docker Compose, либо Go 1.26+.

## Запуск через Docker

```
docker compose up --build
```

## Запуск вручную

Создать `.env` в корне проекта (переменные: `DB_PATH`, `MEDIA_DIR`, `BASE_URL`, `PORT`, `MAX_AVATAR_SIZE_MB`, `MEDIA_PUBLIC_PREFIX`).

`src/docs` — сгенерированный пакет (Swagger-доки), в git не хранится. Перед первым запуском и после любых правок `@Failure`/`@Success`/`@Param` аннотаций в `src/handlers` его нужно сгенерировать заново, иначе `go run .` упадёт с ошибкой отсутствующего пакета `users-api-v1/src/docs`:

```
go install github.com/swaggo/swag/cmd/swag@v1.16.6
swag init -g main.go -o src/docs
```

Затем:

```
go run .
```

## Документация API

`http://localhost:<PORT>/swagger/index.html`

## Маршруты

| Метод  | Путь                 | Описание                                   |
|--------|----------------------|---------------------------------------------|
| POST   | /users               | Создать пользователя                        |
| GET    | /users               | Список пользователей (фильтры, сортировка)  |
| GET    | /users/{id}          | Получить пользователя по id                 |
| PATCH  | /users/{id}          | Обновить пользователя                       |
| DELETE | /users/{id}          | Удалить пользователя                        |
| POST   | /auth/registration   | Регистрация                                 |
| POST   | /auth/login          | Вход, возвращает токен                      |
| GET    | /auth/me             | Текущий пользователь по токену              |
| GET    | /media/{file}        | Файл аватарки                               |
# users-api-v1
