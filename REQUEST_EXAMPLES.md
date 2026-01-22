# 📝 Примеры тел запросов для создания ресурсов

## 🎬 Создание фильма

**Эндпоинт:** `POST http://localhost:8085/api/movies`  
**Требует токен:** Да (Authorization: Bearer $TOKEN)

### Пример тела запроса:

```json
{
  "title": "Inception",
  "description": "A mind-bending thriller about dreams and reality",
  "year": 2010,
  "duration": 148,
  "age_rating": "PG-13",
  "movie_status": "now_showing",
  "genres_id": [1, 2]
}
```

### curl команда:

```bash
curl -X POST http://localhost:8085/api/movies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Inception",
    "description": "A mind-bending thriller about dreams and reality",
    "year": 2010,
    "duration": 148,
    "age_rating": "PG-13",
    "movie_status": "now_showing",
    "genres_id": [1, 2]
  }'
```

### Описание полей:

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `title` | string | ✅ | Название фильма |
| `description` | string | ✅ | Описание фильма |
| `year` | number | ✅ | Год выпуска |
| `duration` | number | ✅ | Длительность в минутах |
| `age_rating` | string | ✅ | Возрастной рейтинг (например: "G", "PG", "PG-13", "R", "NC-17") |
| `movie_status` | string | ✅ | Статус: `"coming_soon"`, `"now_showing"`, `"ended"` |
| `genres_id` | array | ❌ | Массив ID жанров (например: `[1, 2, 3]`) |

### Возможные значения `movie_status`:

- `"coming_soon"` - Скоро в прокате
- `"now_showing"` - Идет сейчас
- `"ended"` - Завершен

---

## 🎭 Создание сеанса

**Эндпоинт:** `POST http://localhost:8085/api/sessions`  
**Требует токен:** Да (Authorization: Bearer $TOKEN)

### Пример тела запроса:

```json
{
  "movie_id": 1,
  "hall_id": 1,
  "start_time": "2026-01-25T18:00:00Z",
  "end_time": "2026-01-25T20:30:00Z"
}
```

### curl команда:

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

### Описание полей:

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `movie_id` | number | ✅ | ID фильма (должен существовать) |
| `hall_id` | number | ✅ | ID зала (должен существовать) |
| `start_time` | string | ✅ | Время начала в формате ISO 8601 (`YYYY-MM-DDTHH:mm:ssZ`) |
| `end_time` | string | ✅ | Время окончания в формате ISO 8601 (должно быть позже `start_time`) |

### Формат даты:

Используйте формат ISO 8601 с UTC временем:
- `"2026-01-25T18:00:00Z"` - 25 января 2026, 18:00 UTC
- `"2026-01-25T20:30:00Z"` - 25 января 2026, 20:30 UTC

### Примеры для разных часовых поясов:

```bash
# Москва (UTC+3) - 21:00
"start_time": "2026-01-25T18:00:00Z"

# Нью-Йорк (UTC-5) - 13:00
"start_time": "2026-01-25T18:00:00Z"

# Токио (UTC+9) - 03:00 следующего дня
"start_time": "2026-01-25T18:00:00Z"
```

**Важно:** Всегда используйте UTC (Z в конце) для консистентности.

---

## 🏛️ Создание зала

**⚠️ ВАЖНО:** Создание залов **НЕ доступно через Gateway**!

Залы можно создавать только напрямую через `cinema-service`:

**Эндпоинт:** `POST http://localhost:8081/halls`  
**Требует токен:** Возможно (зависит от настроек cinema-service)

### Пример тела запроса:

```json
{
  "number": 1
}
```

### curl команда (напрямую к cinema-service):

```bash
curl -X POST http://localhost:8081/halls \
  -H "Content-Type: application/json" \
  -d '{
    "number": 1
  }'
```

### Описание полей:

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `number` | number | ✅ | Номер зала (уникальный) |

### Примеры:

```json
// Зал номер 1
{"number": 1}

// Зал номер 2
{"number": 2}

// Зал номер 5
{"number": 5}
```

---

## 📍 Создание мест в зале

**⚠️ ВАЖНО:** Создание мест **НЕ доступно через Gateway**!

Места создаются напрямую через `cinema-service`:

**Эндпоинт:** `POST http://localhost:8081/halls/:hall_id/seats`  
**Требует токен:** Возможно (зависит от настроек cinema-service)

### Пример тела запроса:

```json
{
  "row": 1,
  "number": 1,
  "type": "standard"
}
```

### curl команда (напрямую к cinema-service):

```bash
curl -X POST http://localhost:8081/halls/1/seats \
  -H "Content-Type: application/json" \
  -d '{
    "row": 1,
    "number": 1,
    "type": "standard"
  }'
```

### Описание полей:

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `row` | number | ✅ | Номер ряда (минимум 1) |
| `number` | number | ✅ | Номер места в ряду (минимум 1) |
| `type` | string | ❌ | Тип места: `"standard"`, `"vip"`, `"wheelchair"` (по умолчанию `"standard"`) |

### Примеры:

```json
// Стандартное место в 1 ряду, место 1
{
  "row": 1,
  "number": 1,
  "type": "standard"
}

// VIP место в 3 ряду, место 5
{
  "row": 3,
  "number": 5,
  "type": "vip"
}

// Место для инвалидной коляски в 1 ряду, место 10
{
  "row": 1,
  "number": 10,
  "type": "wheelchair"
}
```

---

## 🔄 Полный workflow создания данных

### Шаг 1: Создать жанр (если нужно)

```bash
# Напрямую к movie-service (не через gateway)
curl -X POST http://localhost:8083/genres \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Action"
  }'
```

### Шаг 2: Создать фильм

```bash
curl -X POST http://localhost:8085/api/movies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "The Matrix",
    "description": "A computer hacker learns about the true nature of reality",
    "year": 1999,
    "duration": 136,
    "age_rating": "R",
    "movie_status": "now_showing",
    "genres_id": [1]
  }'
```

### Шаг 3: Создать зал

```bash
# Напрямую к cinema-service
curl -X POST http://localhost:8081/halls \
  -H "Content-Type: application/json" \
  -d '{
    "number": 1
  }'
```

### Шаг 4: Создать места в зале

```bash
# Создать несколько мест (пример для зала с 5 рядами по 10 мест)
for row in {1..5}; do
  for seat in {1..10}; do
    curl -X POST http://localhost:8081/halls/1/seats \
      -H "Content-Type: application/json" \
      -d "{\"row\": $row, \"number\": $seat, \"type\": \"standard\"}"
  done
done
```

### Шаг 5: Создать сеанс

```bash
curl -X POST http://localhost:8085/api/sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "movie_id": 1,
    "hall_id": 1,
    "start_time": "2026-01-25T18:00:00Z",
    "end_time": "2026-01-25T20:16:00Z"
  }'
```

---

## 📋 Быстрая справка

### Через Gateway (http://localhost:8085):
- ✅ Фильмы: `POST /api/movies`
- ✅ Сеансы: `POST /api/sessions`
- ❌ Залы: недоступно
- ❌ Места: недоступно
- ❌ Жанры: недоступно

### Напрямую к сервисам:
- **Cinema Service** (http://localhost:8081):
  - Залы: `POST /halls`
  - Места: `POST /halls/:id/seats`
  - Сеансы: `POST /sessions`
  
- **Movie Service** (http://localhost:8083):
  - Фильмы: `POST /movies`
  - Жанры: `POST /genres`

---

## 💡 Полезные команды для генерации дат

### macOS/Linux:

```bash
# Текущее время + 2 часа (для start_time)
date -u -v+2H +"%Y-%m-%dT%H:%M:%SZ"

# Текущее время + 4 часа (для end_time)
date -u -v+4H +"%Y-%m-%dT%H:%M:%SZ"

# Конкретная дата и время
date -u -j -f "%Y-%m-%d %H:%M:%S" "2026-01-25 18:00:00" +"%Y-%m-%dT%H:%M:%SZ"
```

### В curl запросе:

```bash
curl -X POST http://localhost:8085/api/sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"movie_id\": 1,
    \"hall_id\": 1,
    \"start_time\": \"$(date -u -v+2H +'%Y-%m-%dT%H:%M:%SZ')\",
    \"end_time\": \"$(date -u -v+4H +'%Y-%m-%dT%H:%M:%SZ')\"
  }"
```


