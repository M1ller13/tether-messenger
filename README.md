# Tether Messenger

Современный веб-мессенджер, вдохновленный Telegram, с чистым и современным интерфейсом. Поддерживает реальное время, профили пользователей, аутентификацию по email и **end-to-end шифрование сообщений**.

## 🚀 Функциональность

### Аутентификация
- **Регистрация**: Email + пароль с подтверждением по email
- **Вход**: По email и паролю
- **Безопасность**: JWT токены с refresh rotation, автоматическая очистка невалидных токенов

### Чат и сообщения
- **Создание чатов**: Поиск пользователей и создание приватных чатов
- **Отправка сообщений**: Текстовые сообщения в реальном времени
- **История сообщений**: Сохранение и загрузка истории переписки
- **Статус прочтения**: Отметки о прочтении сообщений

### 🔒 End-to-End Шифрование (E2EE)
- **Криптографическая защита**: P-256 ECDH + AES-GCM для шифрования сообщений
- **Аутентификация**: Подписанные prekeys для подтверждения личности
- **Forward Secrecy**: Одноразовые prekeys для защиты прошлых сообщений
- **Автоматическая генерация ключей**: Ключи создаются автоматически при регистрации
- **Визуальные индикаторы**: 🔒 иконки показывают зашифрованные сообщения
- **Совместимость**: Автоматический fallback на plaintext если E2EE недоступен

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
- **Web Crypto API** - генерация и управление ключами шифрования

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
2. Введите email, имя и пароль
3. Подтвердите email по ссылке из письма (или из логов сервера)
4. E2EE ключи автоматически сгенерируются при первом входе

### Вход в систему
1. Откройте `http://localhost:3000/login`
2. Введите email и пароль
3. E2EE ключи автоматически опубликуются на сервере

### Создание чата
1. В поиске введите имя пользователя
2. Выберите пользователя из результатов
3. Начните общение - сообщения автоматически шифруются E2EE
4. Смотрите на 🔒 индикаторы в чате и сообщениях

## 🗄 Структура базы данных

### Таблицы
- **users** - пользователи системы
- **chats** - чаты между пользователями
- **messages** - сообщения в чатах (с поддержкой E2EE полей)
- **device_keys** - публичные ключи устройств для E2EE
- **one_time_prekeys** - одноразовые prekeys для forward secrecy
- **email_verifications** - подтверждение email адресов
- **refresh_tokens** - токены обновления JWT

### Основные поля
```sql
-- users
id (UUID), email, username, display_name, bio, avatar_url, email_verified, last_seen, created_at

-- chats  
id (UUID), user1_id (UUID), user2_id (UUID), created_at

-- messages (с поддержкой E2EE)
id (UUID), chat_id (UUID), sender_id (UUID), content, ciphertext, nonce, alg, ephemeral_pub, is_read, created_at

-- device_keys (E2EE)
id (UUID), user_id (UUID), device_id, identity_key_public, signed_prekey_public, signed_prekey_signature, active, created_at

-- one_time_prekeys (E2EE)
id (UUID), device_key_id (UUID), key_id, public_key, used, created_at, used_at

-- email_verifications
id (UUID), email, token, type, expires_at, used, created_at

-- refresh_tokens
id (UUID), user_id (UUID), token, expires_at, revoked, created_at
```

## 🔌 API Endpoints

### Аутентификация
- `POST /api/auth/register` - регистрация
- `POST /api/auth/login` - вход по email/паролю
- `POST /api/auth/verify-email` - подтверждение email
- `POST /api/auth/refresh-token` - обновление access token
- `POST /api/auth/logout` - выход из системы

### Профиль
- `GET /api/profile` - получить профиль
- `PUT /api/profile` - обновить профиль
- `POST /api/profile/avatar` - загрузить аватар

### Чаты
- `GET /api/chats` - список чатов
- `POST /api/chats` - создать чат
- `GET /api/chats/:chatId` - получить детали чата
- `GET /api/chats/:chatId/messages` - сообщения чата
- `POST /api/messages` - отправить сообщение (поддерживает E2EE)

### 🔒 E2EE Endpoints
- `POST /api/e2ee/device-keys` - опубликовать ключи устройства
- `GET /api/e2ee/prekey-bundle/:userId` - получить prekey-бандл пользователя

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

## 🔐 Безопасность и E2EE

### Архитектура безопасности
- **JWT токены**: Короткоживущие access токены (15 мин) + refresh токены (7 дней)
- **Ротация ключей**: Автоматическая ротация refresh токенов при использовании
- **Хэширование паролей**: Bcrypt для безопасного хранения паролей
- **CORS защита**: Настроенная политика для защиты от CSRF

### End-to-End Шифрование
- **Алгоритмы**: ECDH P-256 для обмена ключами + AES-GCM для шифрования
- **Управление ключами**: Identity ключи (ECDSA P-256) для подписи
- **Forward Secrecy**: Одноразовые prekeys предотвращают компрометацию прошлых сообщений
- **Аутентификация**: Подписанные prekeys подтверждают личность отправителя
- **Хранение ключей**: Приватные ключи хранятся только на клиенте, сервер видит только публичные материалы

### Приватность
- **Сервер не может читать сообщения**: Только зашифрованный контент хранится на сервере
- **Метаданные защищены**: Минимальная информация о пользователях
- **Автоматическая очистка**: Неиспользуемые ключи и токены удаляются

## 🚧 Текущие ограничения

- **WebSocket не использует E2EE**: WS передает зашифрованные payload, но без дополнительной защиты
- **Нет верификации ключей**: Отсутствует система подтверждения fingerprint'ов
- **Нет бэкапов ключей**: При потере устройства ключи не восстанавливаются
- **Моки email**: В dev режиме письма выводятся в консоль

## 🔮 Планы развития

- [ ] Верификация ключей (safety numbers)
- [ ] Бэкапы ключей с паролем
- [ ] Web Push уведомления
- [ ] Файлы и медиа в E2EE
- [ ] Групповые чаты с E2EE
- [ ] Индикаторы набора текста

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Создайте Pull Request

### Для E2EE разработки
- Изучите `client/src/crypto/e2ee.ts` для понимания криптографии
- Тестируйте на разных браузерах (Web Crypto API)
- Проверяйте совместимость с существующими сообщениями

## 📄 Лицензия

MIT License - см. [LICENSE](LICENSE) файл.

## 👨‍💻 Автор

**Tamerlan Akhmedov** - [GitHub](https://github.com/M1ller13)

---

⭐ Если проект вам понравился, поставьте звездочку! 