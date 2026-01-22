
# Сборка всех сервисов
docker-compose build

# Сборка конкретного сервиса
docker-compose build booking-service

# Запустить все сервисы
docker-compose up -d



## Просмотр логов

# Логи всех сервисов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f booking-service
docker-compose logs -f gateway
docker-compose logs -f kafka

# Последние 100 строк логов
docker-compose logs --tail=100

# Остановить все сервисы (контейнеры остаются)
docker-compose stop

# Остановить и удалить контейнеры
docker-compose down

# Остановить, удалить контейнеры и volumes (удалит данные БД!)
docker-compose down -v

# Остановить и удалить контейнеры, volumes и образы
docker-compose down -v --rmi all

# Перезапустить все сервисы
docker-compose restart

# Перезапустить конкретный сервис
docker-compose restart booking-service

# Пересобрать и перезапустить конкретный сервис
docker-compose up -d --build booking-service

# Статус всех контейнеров
docker-compose ps

## Полезные команды

# Проверить доступность Gateway
curl http://localhost:8085/api/movies

# Проверить доступность Kafka UI
open http://localhost:8086

### Если сервис не запускается:

# Посмотреть логи проблемного сервиса
docker-compose logs booking-service

# Пересобрать проблемный сервис
docker-compose build --no-cache booking-service
docker-compose up -d booking-service
```

### Если порт занят:

# Проверить, что использует порт
lsof -i :8085
lsof -i :8082

# Cиды 

./seed_data.sh