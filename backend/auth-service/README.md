# Auth Service

Сервис аутентификации ИС «Коворкинг».  
Отвечает за регистрацию, логин, управление сессиями и выпуск JWT-токенов.

Swagger-спецификация эндпоинтов:  
`../docs/swagger/yaml`

---

## Архитектура

Используется модель:

- **Stateless Access Token**
- **Stateful Refresh Token (с ротацией)**

### Access Token
- Формат: JWT
- Алгоритм: **RS256**
- Короткий TTL (`AUTH_ACCESS_TOKEN_TTL`)
- Проверяется другими сервисами локально по публичному ключу
- Не хранится в БД

### Refresh Token
- JWT (RS256)
- Долгий TTL (`AUTH_REFRESH_TOKEN_TTL`)
- Привязан к серверной сессии
- В БД хранится только `hash(refreshToken)`
- При refresh происходит **ротация токена**

---

## Сессии

При каждом login/register создаётся новая сессия.

Хранится:
- `sessionID`
- `userID`
- `expiresAt`
- `revoked`
- `tokenHash`

Refresh токен повторно использовать нельзя — старая сессия инвалидируется.

---

## Основные операции

### Register
Создание пользователя → создание сессии → выдача access + refresh.

### Login
Проверка учётных данных → новая сессия → выдача токенов.

### Refresh
- Проверка подписи и claims
- Проверка сессии
- Проверка hash refresh токена
- Инвалидация старой сессии
- Выдача новой пары токенов

### Logout
Инвалидация текущей сессии.

### Logout All
Инвалидация всех активных сессий пользователя.

---

## Интеграция

Другие сервисы (Gateway, Booking и др.):

- получают публичный RSA-ключ
- валидируют access token локально
- не обращаются к auth-service для проверки токена

---

## Конфигурация

Через `.env`:

- DB параметры
- `AUTH_ACCESS_TOKEN_TTL`
- `AUTH_REFRESH_TOKEN_TTL`
- `AUTH_PRIVATE_KEY`
- `AUTH_PUBLIC_KEY`

---

## Безопасность

- Асимметричная подпись (RS256)
- Stateless access
- Refresh rotation
- Session revoke
- Хранение только hash refresh токена