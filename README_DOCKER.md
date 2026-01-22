# Команды для сборки и запуска проекта

## Быстрый старт

### 1. Сборка и запуск всех сервисов (рекомендуется)

```bash
# Включить BuildKit для ускорения сборки
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

# Собрать и запустить все сервисы
docker-compose up -d --build
```

### 2. Только сборка (без запуска)

```bash
# Сборка всех сервисов
docker-compose build

# Сборка с BuildKit
DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker-compose build

# Сборка конкретного сервиса
docker-compose build booking-service
```

### 3. Только запуск (если уже собрано)

```bash
# Запустить все сервисы
docker-compose up -d

# Запустить в foreground (с логами)
docker-compose up
```

## Просмотр логов

```bash
# Логи всех сервисов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f booking-service
docker-compose logs -f gateway
docker-compose logs -f kafka

# Последние 100 строк логов
docker-compose logs --tail=100
```

## Остановка и очистка

```bash
# Остановить все сервисы (контейнеры остаются)
docker-compose stop

# Остановить и удалить контейнеры
docker-compose down

# Остановить, удалить контейнеры и volumes (удалит данные БД!)
docker-compose down -v

# Остановить и удалить контейнеры, volumes и образы
docker-compose down -v --rmi all
```

## Перезапуск сервисов

```bash
# Перезапустить все сервисы
docker-compose restart

# Перезапустить конкретный сервис
docker-compose restart booking-service

# Пересобрать и перезапустить конкретный сервис
docker-compose up -d --build booking-service
```

## Проверка статуса

```bash
# Статус всех контейнеров
docker-compose ps

# Проверить здоровье сервисов
docker-compose ps --format json | jq '.[] | {name: .Name, status: .State, health: .Health}'
```

## Полезные команды

```bash
# Выполнить команду в контейнере
docker-compose exec booking-service sh
docker-compose exec booking-postgres psql -U postgres -d booking_db

# Просмотреть использование ресурсов
docker stats

# Очистить неиспользуемые образы и кэш
docker system prune -a

# Просмотреть логи конкретного сервиса за последние 5 минут
docker-compose logs --since 5m booking-service
```

## Порядок запуска сервисов

Сервисы запускаются в правильном порядке благодаря `depends_on`:

1. **Базы данных** (PostgreSQL) - запускаются первыми
2. **Kafka** - запускается после БД
3. **Микросервисы** - запускаются после своих БД и зависимостей
4. **Gateway** - запускается последним, после всех микросервисов

## Проверка работоспособности

После запуска проверьте:

```bash
# Проверить, что все контейнеры запущены
docker-compose ps

# Проверить доступность Gateway
curl http://localhost:8085/api/movies

# Проверить доступность Kafka UI
open http://localhost:8086
```

## Решение проблем

### Если сервис не запускается:

```bash
# Посмотреть логи проблемного сервиса
docker-compose logs booking-service

# Пересобрать проблемный сервис
docker-compose build --no-cache booking-service
docker-compose up -d booking-service
```

### Если порт занят:

```bash
# Проверить, что использует порт
lsof -i :8085
lsof -i :8082

# Остановить конфликтующий процесс или изменить порт в docker-compose.yml
```

### Очистка и полная пересборка:

```bash
# Остановить все
docker-compose down -v

# Очистить build cache
docker builder prune

# Полная пересборка
DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker-compose build --no-cache
docker-compose up -d
```

