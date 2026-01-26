#!/bin/bash

# Скрипт для создания тестовых данных в cinema-hall системе
# Использование: ./seed_data.sh

set -e

GATEWAY_URL="http://localhost:8085"
MOVIE_SERVICE_URL="http://localhost:8083"
CINEMA_SERVICE_URL="http://localhost:8081"

echo "🚀 Начало создания тестовых данных..."
echo ""

# 1. Регистрация/Вход пользователя
echo "=== 1. Получение токена ==="
TOKEN=$(curl -s -X POST "$GATEWAY_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@cinema.com", "password": "admin123"}' \
  | jq -r '.access_token // empty')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "Регистрация нового пользователя..."
  curl -s -X POST "$GATEWAY_URL/api/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"email": "admin@cinema.com", "password": "admin123", "name": "Admin User"}' > /dev/null
  
  TOKEN=$(curl -s -X POST "$GATEWAY_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email": "admin@cinema.com", "password": "admin123"}' \
    | jq -r '.access_token // empty')
fi

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "❌ Ошибка: не удалось получить токен"
  exit 1
fi

echo "✅ Токен получен"
echo ""

# 2. Создание жанров
echo "=== 2. Создание жанров ==="
GENRES=("Action" "Drama" "Comedy" "Thriller" "Sci-Fi" "Horror")
for genre in "${GENRES[@]}"; do
  RESPONSE=$(curl -s -X POST "$MOVIE_SERVICE_URL/genres/" \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"$genre\"}")
  GENRE_ID=$(echo "$RESPONSE" | jq -r '.id // empty')
  if [ -n "$GENRE_ID" ] && [ "$GENRE_ID" != "null" ]; then
    echo "  ✅ Создан жанр: $genre (ID: $GENRE_ID)"
  else
    echo "  ⚠️  Жанр $genre уже существует или ошибка создания"
  fi
done
echo ""

# 3. Создание фильмов
echo "=== 3. Создание фильмов ==="
MOVIES=(
  '{"title": "Inception", "description": "A mind-bending thriller about dreams and reality", "year": 2010, "duration": 148, "age_rating": "PG-13", "movie_status": "now_showing", "genres_id": [1, 4]}'
  '{"title": "The Matrix", "description": "A computer hacker learns about the true nature of reality", "year": 1999, "duration": 136, "age_rating": "R", "movie_status": "now_showing", "genres_id": [1]}'
  '{"title": "The Dark Knight", "description": "Batman faces the Joker", "year": 2008, "duration": 152, "age_rating": "PG-13", "movie_status": "now_showing", "genres_id": [1, 2]}'
)

for movie in "${MOVIES[@]}"; do
  RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/movies" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "$movie")
  MOVIE_ID=$(echo "$RESPONSE" | jq -r '.id // empty')
  MOVIE_TITLE=$(echo "$RESPONSE" | jq -r '.title // empty')
  if [ -n "$MOVIE_ID" ] && [ "$MOVIE_ID" != "null" ]; then
    echo "  ✅ Создан фильм: $MOVIE_TITLE (ID: $MOVIE_ID)"
  else
    echo "  ❌ Ошибка создания фильма: $RESPONSE"
  fi
done
echo ""

# 4. Создание залов
echo "=== 4. Создание залов ==="
HALL_IDS=()
for hall_num in {1..3}; do
  RESPONSE=$(curl -s -X POST "$CINEMA_SERVICE_URL/halls" \
    -H "Content-Type: application/json" \
    -d "{\"number\": $hall_num}")
  HALL_ID=$(echo "$RESPONSE" | jq -r '.id // empty')
  if [ -n "$HALL_ID" ] && [ "$HALL_ID" != "null" ]; then
    HALL_IDS+=($HALL_ID)
    echo "  ✅ Создан зал №$hall_num (ID: $HALL_ID)"
  else
    echo "  ⚠️  Зал №$hall_num уже существует или ошибка создания"
  fi
done
echo ""

# 5. Создание мест
echo "=== 5. Создание мест ==="
for hall_id in "${HALL_IDS[@]}"; do
  echo "  Создание мест для зала ID: $hall_id"
  SEAT_COUNT=0
  for row in {1..5}; do
    for seat_num in {1..8}; do
      RESPONSE=$(curl -s -X POST "$CINEMA_SERVICE_URL/halls/$hall_id/seats" \
        -H "Content-Type: application/json" \
        -d "{\"row\": $row, \"number\": $seat_num, \"type\": \"standard\"}")
      SEAT_ID=$(echo "$RESPONSE" | jq -r '.id // empty')
      if [ -n "$SEAT_ID" ] && [ "$SEAT_ID" != "null" ]; then
        ((SEAT_COUNT++))
      fi
    done
  done
  echo "    ✅ Создано мест: $SEAT_COUNT"
done
echo ""

echo "🎉 Все тестовые данные успешно созданы!"
echo ""
echo "📊 Итоговая сводка:"
echo "  - Жанры: ${#GENRES[@]}"
echo "  - Фильмы: ${#MOVIES[@]}"
echo "  - Залы: ${#HALL_IDS[@]}"
echo "  - Места: ~$(( ${#HALL_IDS[@]} * 5 * 8 ))"



