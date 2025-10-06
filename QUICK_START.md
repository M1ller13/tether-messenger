# 🚀 Быстрый запуск Tether Messenger

## Вариант 1: Локальный запуск (рекомендуется для разработки)

### Предварительные требования
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### 1. Установка зависимостей
```bash
# Backend
cd server
go mod download

# Frontend
cd ../client
npm install
```

### 2. Настройка базы данных
```sql
-- Создайте базу данных PostgreSQL
CREATE DATABASE tether_messenger;

-- Создайте пользователя (опционально)
CREATE USER postgres WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE tether_messenger TO postgres;
```

### 3. Запуск сервисов

#### Автоматический запуск (Windows)
```bash
# PowerShell
.\start-dev.ps1

# Или Batch файл
start-dev.bat
```

#### Ручной запуск
```bash
# Терминал 1: Backend
cd server
go run main.go

# Терминал 2: Frontend
cd client
npm run dev
```

### 4. Доступ к приложению
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8081
- **Health Check**: http://localhost:8081/health

---

## Вариант 2: Docker (рекомендуется для production)

### Предварительные требования
- Docker Desktop
- Docker Compose

### 1. Запуск development окружения
```bash
# Запуск только базы данных и Redis
docker-compose -f docker-compose.dev.yml up -d

# Затем запустите backend и frontend локально
cd server && go run main.go
cd client && npm run dev
```

### 2. Полный Docker запуск
```bash
# Сборка и запуск всех сервисов
docker-compose up -d --build

# Или используйте Makefile
make build
make up
```

### 3. Доступ к приложению
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8081
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

---

## 🎯 Текущий функционал

### ✅ Реализовано
1. **Аутентификация по email**
   - Регистрация с подтверждением email
   - Вход с email/паролем
   - JWT refresh tokens
   - Сброс пароля через email

2. **Система досок задач**
   - Создание личных и командных досок
   - Колонки (To Do, In Progress, Done)
   - Карточки с CRM полями
   - Права доступа

3. **CRM функциональность**
   - Lead name, contact email, company
   - Value, priority, status
   - Назначение ответственных

4. **API Endpoints**
   - `/api/auth/*` - аутентификация
   - `/api/boards/*` - управление досками
   - `/api/columns/*` - управление колонками
   - `/api/cards/*` - управление карточками

### 🔄 В процессе
- UI компоненты для досок
- Drag & Drop функциональность
- WebSocket интеграция для досок

---

## 🧪 Тестирование

### Регистрация нового пользователя
1. Откройте http://localhost:3000/register
2. Заполните форму регистрации
3. Проверьте консоль сервера для ссылки подтверждения email
4. Перейдите по ссылке для активации аккаунта

### Создание доски
1. Войдите в систему
2. Перейдите в раздел досок (пока в разработке)
3. Создайте новую доску через API:
```bash
curl -X POST http://localhost:8081/api/boards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"name": "My Board", "type": "personal"}'
```

---

## 🐛 Troubleshooting

### Проблемы с базой данных
```bash
# Проверьте подключение к PostgreSQL
psql -h localhost -U postgres -d tether_messenger

# Пересоздайте базу данных
dropdb tether_messenger
createdb tether_messenger
```

### Проблемы с Docker
```bash
# Проверьте статус Docker
docker info

# Перезапустите Docker Desktop
# Очистите контейнеры и volumes
docker-compose down -v
docker system prune -f
```

### Проблемы с портами
Если порты заняты, измените их в конфигурации:
- Backend: `server/config/config.go` (SERVER_PORT)
- Frontend: `client/vite.config.ts` (port)

---

## 📱 Интерфейс

### Страницы
- `/login` - Вход в систему
- `/register` - Регистрация
- `/forgot-password` - Сброс пароля
- `/reset-password` - Новый пароль
- `/verify-email` - Подтверждение email
- `/chats` - Чаты (базовая функциональность)

### Планируемые страницы
- `/boards` - Доски задач
- `/workspaces` - Рабочие пространства
- `/profile` - Профиль пользователя

---

## 🔧 Разработка

### Структура проекта
```
tether-messenger/
├── server/          # Go backend
│   ├── handlers/    # API handlers
│   ├── models/      # Database models
│   ├── routes/      # API routes
│   └── utils/       # Utilities
├── client/          # React frontend
│   ├── src/
│   │   ├── pages/   # React pages
│   │   ├── components/ # React components
│   │   └── hooks/   # Custom hooks
└── docs/            # Documentation
```

### Полезные команды
```bash
# Проверка линтера
cd server && go vet ./...
cd client && npm run lint

# Запуск тестов
cd server && go test ./...

# Сборка для production
cd server && go build -o tether-server main.go
cd client && npm run build
```

---

## 🎉 Готово к использованию!

Проект готов к разработке и тестированию. Основная функциональность аутентификации и API для досок задач реализована. Следующий шаг - создание UI компонентов для работы с досками.
