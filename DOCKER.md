# Docker Setup для Tether Messenger

Этот документ описывает как запустить Tether Messenger с помощью Docker.

## Быстрый старт

### 1. Клонирование репозитория
```bash
git clone <repository-url>
cd tether-messenger
```

### 2. Запуск в production режиме
```bash
# Сборка и запуск всех сервисов
make build
make up

# Или напрямую через docker-compose
docker-compose up -d
```

### 3. Запуск в development режиме
```bash
# Запуск только базы данных и Redis
make dev-start

# Затем в отдельных терминалах:
cd server && go run main.go
cd client && npm run dev
```

## Доступные команды

### Production команды
- `make build` - Сборка всех Docker образов
- `make up` - Запуск всех сервисов
- `make down` - Остановка всех сервисов
- `make logs` - Просмотр логов
- `make clean` - Полная очистка (контейнеры, volumes, образы)

### Development команды
- `make dev-up` - Запуск dev окружения (только DB + Redis)
- `make dev-down` - Остановка dev окружения
- `make dev-logs` - Просмотр логов dev окружения

### Утилиты
- `make install-deps` - Установка зависимостей
- `make test` - Запуск тестов

## Архитектура

### Сервисы
1. **PostgreSQL** (порт 5432) - Основная база данных
2. **Backend** (порт 8081) - Go API сервер
3. **Frontend** (порт 3000) - React приложение с Nginx
4. **Redis** (порт 6379) - Кэширование и сессии

### Сети
- `tether-network` - Production сеть
- `tether-dev-network` - Development сеть

### Volumes
- `postgres_data` - Данные PostgreSQL
- `redis_data` - Данные Redis
- `postgres_dev_data` - Данные PostgreSQL для разработки
- `redis_dev_data` - Данные Redis для разработки

## Конфигурация

### Environment переменные

#### Backend
- `DB_HOST` - Хост базы данных (по умолчанию: postgres)
- `DB_PORT` - Порт базы данных (по умолчанию: 5432)
- `DB_USER` - Пользователь БД (по умолчанию: postgres)
- `DB_PASSWORD` - Пароль БД (по умолчанию: password)
- `DB_NAME` - Имя БД (по умолчанию: tether_messenger)
- `JWT_SECRET` - Секретный ключ для JWT
- `SERVER_PORT` - Порт сервера (по умолчанию: 8081)

#### PostgreSQL
- `POSTGRES_DB` - Имя базы данных
- `POSTGRES_USER` - Пользователь
- `POSTGRES_PASSWORD` - Пароль

## Доступ к сервисам

После запуска сервисы будут доступны по следующим адресам:

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8081
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## Разработка

### Локальная разработка с Docker
1. Запустите dev окружение: `make dev-start`
2. Запустите backend локально: `cd server && go run main.go`
3. Запустите frontend локально: `cd client && npm run dev`

### Отладка
```bash
# Просмотр логов конкретного сервиса
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres

# Подключение к контейнеру
docker-compose exec backend sh
docker-compose exec postgres psql -U postgres -d tether_messenger
```

### Пересборка после изменений
```bash
# Пересборка и перезапуск
docker-compose up -d --build

# Или только конкретного сервиса
docker-compose up -d --build backend
```

## Troubleshooting

### Проблемы с портами
Если порты заняты, измените их в `docker-compose.yml`:
```yaml
ports:
  - "3001:80"  # Frontend на порту 3001
  - "8082:8081"  # Backend на порту 8082
```

### Проблемы с базой данных
```bash
# Сброс базы данных
docker-compose down -v
docker-compose up -d

# Подключение к БД для отладки
docker-compose exec postgres psql -U postgres -d tether_messenger
```

### Очистка
```bash
# Полная очистка
make clean

# Очистка только volumes
docker-compose down -v
```

## Production развертывание

Для production развертывания:

1. Измените пароли и секретные ключи
2. Настройте SSL/TLS
3. Используйте внешнюю базу данных
4. Настройте мониторинг и логирование
5. Используйте reverse proxy (nginx/traefik)

Пример production конфигурации:
```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  backend:
    environment:
      JWT_SECRET: ${JWT_SECRET}
      DB_PASSWORD: ${DB_PASSWORD}
    # ... остальная конфигурация
```
