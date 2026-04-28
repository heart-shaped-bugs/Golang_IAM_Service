# Golang_IAM_Service

Сервис аутентификации и управления пользователями.  
Регистрирует пользователей, выполняет вход и выдаёт JWT-токены для других микросервисов.

## Возможности

- Регистрация по email и паролю
- Аутентификация с выдачей JWT (HS256, 24 часа)
- Автоматические миграции БД при старте
- Поддержка `.env` для конфигурации

## Требования

- Go 1.20+
- PostgreSQL 12+
- Docker (опционально, для запуска БД)

## Установка

1. Клонируй репозиторий:

   ```bash
   git clone https://github.com/heart-shaped-bugs/Golang_IAM_Service.git
   cd Golang_IAM_Service
   ```

2. Создай `.env` на основе шаблона:

   ```bash
   cp .env.example .env
   # Отредактируй .env под свою БД
   ```

3. Запусти сервер:

   ```bash
   go run cmd/iam/main.go
   ```

## Запуск с Docker (БД)

Запусти PostgreSQL в контейнере:

```bash
docker run --name iam-db \
  -e POSTGRES_DB=iam \
  -e POSTGRES_USER=iam \
  -e POSTGRES_PASSWORD=secret \
  -p 5432:5432 \
  -d postgres:15
```

> Сервис автоматически применит миграции при старте.

## API

### `POST /register`

Регистрация нового пользователя.

**Запрос:**

```json
{
  "email": "user@example.com",
  "password": "securePass123"
}
```

**Ответ:** `201 Created` (пустое тело)

---

### `POST /login`

Аутентификация и получение JWT.

**Запрос:**

```json
{
  "email": "user@example.com",
  "password": "securePass123"
}
```

**Успешный ответ (`200 OK`):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx"
}
```

**Ошибки:**

- `401 Unauthorized` — неверные учётные данные
- `409 Conflict` — пользователь уже существует
- `400 Bad Request` — некорректный email

## Безопасность

- Пароли хешируются с помощью **bcrypt**
- JWT подписывается с использованием **HS256**
- Секреты задаются через переменные окружения

## Тестирование

Пример регистрации:

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'
```

Пример входа:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'
```

## Структура проекта

```bash
cmd/iam/            # точка входа
internal/
├── entities/       # доменные сущности
├── repositories/   # интерфейсы и реализации (Postgres)
└── usecases/       # бизнес-логика (AuthService)
api/http/           # HTTP-обработчики и маршрутизация
migrations/         # SQL-миграции
```
