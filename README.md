# Tether Messenger

Современный веб-мессенджер, вдохновленный Telegram, с чистым и современным интерфейсом. Поддерживает реальное время, профили пользователей и аутентификацию по номеру телефона.

## 🚀 Функциональность

### Аутентификация
- **Регистрация**: Двухшаговый процесс (имя + телефон → подтверждение кода)
- **Вход**: По номеру телефона с подтверждением SMS-кода
- **Безопасность**: JWT токены, автоматическая очистка невалидных токенов

### Чат и сообщения
- **Создание чатов**: Поиск пользователей и создание приватных чатов
- **Отправка сообщений**: Текстовые сообщения в реальном времени
- **История сообщений**: Сохранение и загрузка истории переписки
- **Статус прочтения**: Отметки о прочтении сообщений

### Профили пользователей
- **Редактирование профиля**: Имя, био, аватар
- **Загрузка аватаров**: Поддержка изображений профиля
- **Поиск пользователей**: Поиск по имени и username
- **Статус онлайн**: Отображение времени последней активности

## 🛠 Технологии

### Backend
- **Go** - основной язык сервера
- **Fiber** - веб-фреймворк
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - основная база данных
- **JWT** - аутентификация
- **UUID** - уникальные идентификаторы

### Frontend
- **React 18** - UI библиотека
- **TypeScript** - типизированный JavaScript
- **Vite** - сборщик и dev-сервер
- **Tailwind CSS** - стилизация
- **React Router** - маршрутизация

## 📦 Установка и запуск

### Предварительные требования
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Git

### Backend

1. **Клонируйте репозиторий**
```bash
git clone https://github.com/M1ller13/tether-messenger.git
cd tether-messenger/server
```

2. **Настройте базу данных**
```bash
# Создайте базу данных PostgreSQL
createdb tether_messenger

# Или используйте существующую базу и обновите config/config.go
```

3. **Установите зависимости**
```bash
go mod download
```

4. **Запустите сервер**
```bash
go run main.go
```

Сервер запустится на `http://localhost:8081`

### Frontend

1. **Перейдите в папку клиента**
```bash
cd ../client
```

2. **Установите зависимости**
```bash
npm install
```

3. **Запустите dev-сервер**
```bash
npm run dev
```

Приложение откроется на `http://localhost:3000`

## 🔧 Конфигурация

### Backend конфигурация
Отредактируйте `server/config/config.go`:

```go
var Config = struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    JWTSecret  string
    ServerPort string
}{
    DBHost:     "localhost",
    DBPort:     "5432",
    DBUser:     "postgres",
    DBPassword: "your_password",
    DBName:     "tether_messenger",
    JWTSecret:  "your_jwt_secret",
    ServerPort: "8081",
}
```

### Frontend конфигурация
Прокси настроен в `client/vite.config.ts` для перенаправления API запросов на backend.

## 📱 Использование

### Регистрация нового пользователя
1. Откройте `http://localhost:3000/register`
2. Введите имя и номер телефона
3. Получите код в консоли сервера
4. Введите код для подтверждения

### Вход в систему
1. Откройте `http://localhost:3000/login`
2. Введите номер телефона
3. Получите код в консоли сервера
4. Введите код для входа

### Создание чата
1. В поиске введите имя пользователя
2. Выберите пользователя из результатов
3. Начните общение

## 🗄 Структура базы данных

### Таблицы
- **users** - пользователи системы
- **chats** - чаты между пользователями
- **messages** - сообщения в чатах
- **verification_codes** - коды подтверждения

### Основные поля
```sql
-- users
id (UUID), phone, username, display_name, bio, avatar_url, last_seen, created_at

-- chats  
id (UUID), user1_id (UUID), user2_id (UUID), created_at

-- messages
id (UUID), chat_id (UUID), sender_id (UUID), content, is_read, created_at

-- verification_codes
id (UUID), phone, code, display_name, expires_at, created_at
```

## 🔌 API Endpoints

### Аутентификация
- `POST /api/auth/register` - регистрация
- `POST /api/auth/request-code` - запрос кода для входа
- `POST /api/auth/verify-code` - подтверждение кода

### Профиль
- `GET /api/profile` - получить профиль
- `PUT /api/profile` - обновить профиль
- `POST /api/profile/avatar` - загрузить аватар

### Чаты
- `GET /api/chats` - список чатов
- `POST /api/chats` - создать чат
- `GET /api/chats/:chatId/messages` - сообщения чата
- `POST /api/messages` - отправить сообщение

### Пользователи
- `GET /api/users/search` - поиск пользователей

## 🚀 Развертывание

### Production
1. **Сборка frontend**
```bash
cd client
npm run build
```

2. **Сборка backend**
```bash
cd server
go build -o tether-server main.go
```

3. **Настройка reverse proxy** (nginx/apache) для статических файлов

4. **Запуск backend**
```bash
./tether-server
```

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Создайте Pull Request

## 📄 Лицензия

MIT License - см. [LICENSE](LICENSE) файл.

## 👨‍💻 Автор

**Tamerlan Akhmedov** - [GitHub](https://github.com/M1ller13)

---

⭐ Если проект вам понравился, поставьте звездочку! 