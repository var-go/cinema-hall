# 📚 Руководство по работе с Gateway API

**Base URL**: `http://localhost:8085`

Все запросы к защищенным эндпоинтам требуют JWT токен в заголовке `Authorization: Bearer <token>`

---

## 🔐 Аутентификация

### 1. Регистрация пользователя

```bash
curl -X POST http://localhost:8085/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

**Важно:** Все поля обязательны:
- `email` - валидный email адрес
- `password` - минимум 6 символов
- `name` - имя пользователя

**Ответ:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "role": "user"
}
```

### 2. Вход (получение JWT токена)

```bash
curl -X POST http://localhost:8085/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Важно:** Поле называется `access_token`, а не `token`!

**Сохраните токен для последующих запросов:**
```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 🎬 Фильмы (Movies)

### Получить список всех фильмов

```bash
curl http://localhost:8085/api/movies
```

### Получить фильм по ID

```bash
curl http://localhost:8085/api/movies/1
```

### Создать фильм (требует аутентификации)

```bash
curl -X POST http://localhost:8085/api/movies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Inception",
    "description": "A mind-bending thriller",
    "duration": 148,
    "genre_id": 1
  }'
```

---

## 🎭 Сеансы (Sessions)

### Получить список всех сеансов

```bash
curl http://localhost:8085/api/sessions
```

### Получить сеанс по ID

```bash
curl http://localhost:8085/api/sessions/1
```

### Создать сеанс (требует аутентификации)

```bash
curl -X POST http://localhost:8085/api/sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "movie_id": 1,
    "hall_id": 1,
    "start_time": "2026-01-25T18:00:00Z",
    "end_time": "2026-01-25T20:30:00Z"
  }'
```

### Получить агрегированную информацию о сеансе

Возвращает сеанс вместе с информацией о фильме и зале:

```bash
curl http://localhost:8085/api/sessions/1/aggregate
```

**Ответ:**
```json
{
  "session": {
    "id": 1,
    "movie_id": 1,
    "hall_id": 1,
    "start_time": "2026-01-25T18:00:00Z",
    "end_time": "2026-01-25T20:30:00Z"
  },
  "movie": {
    "id": 1,
    "title": "Inception",
    "duration": 148
  },
  "hall": {
    "id": 1,
    "name": "Hall 1",
    "capacity": 100
  }
}
```

---

## 🎫 Бронирования (Bookings)

**Все эндпоинты бронирований требуют JWT токен!**

### Получить список всех бронирований

```bash
curl http://localhost:8085/api/bookings \
  -H "Authorization: Bearer $TOKEN"
```

### Получить бронирование по ID

```bash
curl http://localhost:8085/api/bookings/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Создать бронирование

```bash
curl -X POST http://localhost:8085/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "session_id": 1,
    "user_id": 1,
    "seats_id": [1, 2, 3]
  }'
```

**Ответ:**
```json
{
  "id": 1,
  "session_id": 1,
  "user_id": 1,
  "booking_status": "pending",
  "payment_status": "pending",
  "expires_at": "2026-01-22T12:15:00Z",
  "booked_seats": [
    {"id": 1, "seat_id": 1},
    {"id": 2, "seat_id": 2},
    {"id": 3, "seat_id": 3}
  ]
}
```

### Обновить бронирование

```bash
curl -X PATCH http://localhost:8085/api/bookings/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "seats_id": [1, 2]
  }'
```

### Подтвердить бронирование (после оплаты)

```bash
curl -X POST http://localhost:8085/api/bookings/1/confirm \
  -H "Authorization: Bearer $TOKEN"
```

**Важно:** После подтверждения:
- Статус меняется на `confirmed`
- Статус оплаты меняется на `paid`
- Отправляется событие в Kafka

### Отменить бронирование

```bash
curl -X POST http://localhost:8085/api/bookings/1/cancel \
  -H "Authorization: Bearer $TOKEN"
```

**Важно:** После отмены:
- Статус меняется на `cancelled`
- Места освобождаются
- Отправляется событие в Kafka

### Удалить бронирование

```bash
curl -X DELETE http://localhost:8085/api/bookings/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📝 Примеры полного workflow

### Сценарий 1: Регистрация → Вход → Создание бронирования → Подтверждение

```bash
# 1. Регистрация
curl -X POST http://localhost:8085/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123", "name": "John Doe"}'

# 2. Вход и сохранение токена
TOKEN=$(curl -s -X POST http://localhost:8085/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}' \
  | jq -r '.access_token')

echo "Token: $TOKEN"

# 3. Просмотр доступных сеансов
curl http://localhost:8085/api/sessions

# 4. Просмотр деталей сеанса
curl http://localhost:8085/api/sessions/1/aggregate

# 5. Создание бронирования
curl -X POST http://localhost:8085/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "session_id": 1,
    "user_id": 1,
    "seats_id": [1, 2]
  }'

# 6. Подтверждение бронирования (после оплаты)
curl -X POST http://localhost:8085/api/bookings/1/confirm \
  -H "Authorization: Bearer $TOKEN"
```

### Сценарий 2: Создание фильма → Создание сеанса → Бронирование

```bash
# 1. Вход (получить токен)
TOKEN=$(curl -s -X POST http://localhost:8085/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}' \
  | jq -r '.access_token')

# 2. Создать фильм
curl -X POST http://localhost:8085/api/movies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "The Matrix",
    "description": "A sci-fi classic",
    "duration": 136,
    "genre_id": 1
  }'

# 3. Создать сеанс
curl -X POST http://localhost:8085/api/sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "movie_id": 1,
    "hall_id": 1,
    "start_time": "2026-01-25T20:00:00Z",
    "end_time": "2026-01-25T22:16:00Z"
  }'

# 4. Создать бронирование
curl -X POST http://localhost:8085/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "session_id": 1,
    "user_id": 1,
    "seats_id": [5, 6]
  }'
```

---

## 🔧 Полезные команды

### Сохранить токен в переменную окружения

```bash
# Linux/macOS
export TOKEN=$(curl -s -X POST http://localhost:8085/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}' \
  | jq -r '.access_token')

# Использовать в запросах
curl http://localhost:8085/api/bookings \
  -H "Authorization: Bearer $TOKEN"
```

### Форматировать JSON ответ

```bash
curl http://localhost:8085/api/movies | jq .
```

### Показать только заголовки

```bash
curl -I http://localhost:8085/api/movies
```

### Показать полный ответ (включая заголовки)

```bash
curl -v http://localhost:8085/api/movies
```

---

## ⚠️ Важные замечания

### JWT Токен

- Токен получается при входе через `/api/auth/login`
- Токен нужно передавать в заголовке: `Authorization: Bearer <token>`
- Токен действителен до истечения срока (зависит от настроек user-service)
- При истечении токена получите `401 Unauthorized`

### Статусы бронирования

- `pending` - бронирование создано, ожидает оплаты (действует 15 минут)
- `confirmed` - бронирование подтверждено и оплачено
- `cancelled` - бронирование отменено пользователем
- `expired` - бронирование истекло (не оплачено в течение 15 минут)
- `finished` - сеанс завершен, бронирование выполнено

### Таймаут бронирования

- Бронирование со статусом `pending` автоматически истекает через 15 минут
- После истечения места освобождаются автоматически
- Подтверждение бронирования (`/confirm`) должно произойти до истечения таймаута

### Kafka события

- При подтверждении бронирования (`/confirm`) отправляется событие в Kafka
- При отмене бронирования (`/cancel`) отправляется событие в Kafka
- События можно просмотреть в Kafka UI: http://localhost:8086

---

## 🐛 Обработка ошибок

### 401 Unauthorized
```json
{
  "error": "missing authorization header"
}
```
**Решение:** Добавьте заголовок `Authorization: Bearer <token>`

### 404 Not Found
```json
{
  "error": "booking not found"
}
```
**Решение:** Проверьте правильность ID

### 400 Bad Request
```json
{
  "error": "invalid JSON"
}
```
**Решение:** Проверьте формат JSON в теле запроса

### 500 Internal Server Error
**Решение:** Проверьте логи сервисов:
```bash
docker-compose logs booking-service
docker-compose logs gateway
```

---

## 📊 Коды ответов

| Код | Описание |
|-----|----------|
| 200 | Успешный запрос |
| 201 | Ресурс создан |
| 400 | Неверный запрос |
| 401 | Не авторизован |
| 404 | Ресурс не найден |
| 500 | Внутренняя ошибка сервера |
| 502 | Сервис недоступен |

---

## 🔗 Полезные ссылки

- **Gateway**: http://localhost:8085
- **Kafka UI**: http://localhost:8086
- **User Service**: http://localhost:8080
- **Cinema Service**: http://localhost:8081
- **Booking Service**: http://localhost:8082
- **Movie Service**: http://localhost:8083

---

## 📝 Примеры для Postman

### Настройка переменных

1. Создайте переменную `base_url` = `http://localhost:8085`
2. Создайте переменную `token` (будет заполняться автоматически)

### Pre-request Script для получения токена

```javascript
// В коллекции или запросе на /api/auth/login
pm.sendRequest({
    url: pm.variables.get("base_url") + "/api/auth/login",
    method: 'POST',
    header: {'Content-Type': 'application/json'},
    body: {
        mode: 'raw',
        raw: JSON.stringify({
            email: "user@example.com",
            password: "password123"
        })
    }
}, function (err, res) {
    if (res.json().access_token) {
        pm.environment.set("token", res.json().access_token);
    }
});
```

### Использование токена в запросах

В заголовках запросов добавьте:
```
Authorization: Bearer {{token}}
```

---

Готово! Теперь вы можете работать с API через Gateway. 🚀

