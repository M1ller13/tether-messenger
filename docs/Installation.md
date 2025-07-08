# Установка и настройка

Это руководство поможет вам установить и настроить Tether Messenger на вашей системе.

## 📋 Предварительные требования

### Системные требования
- **ОС**: Windows 10+, macOS 10.15+, Ubuntu 18.04+
- **RAM**: Минимум 2GB (рекомендуется 4GB+)
- **Дисковое пространство**: 1GB свободного места

### Необходимое ПО
- **Go**: версия 1.21 или выше
- **Node.js**: версия 18 или выше
- **PostgreSQL**: версия 14 или выше
- **Git**: для клонирования репозитория

## 🔧 Установка зависимостей

### 1. Установка Go

#### Windows
```bash
# Скачайте установщик с https://golang.org/dl/
# Запустите .msi файл и следуйте инструкциям
```

#### macOS
```bash
# Через Homebrew
brew install go

# Или скачайте с https://golang.org/dl/
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install golang-go
```

### 2. Установка Node.js

#### Windows/macOS
```bash
# Скачайте с https://nodejs.org/
# Рекомендуется LTS версия
```

#### Linux
```bash
# Через apt
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs

# Или через snap
sudo snap install node --classic
```

### 3. Установка PostgreSQL

#### Windows
```bash
# Скачайте с https://www.postgresql.org/download/windows/
# Запустите установщик
```

#### macOS
```bash
# Через Homebrew
brew install postgresql
brew services start postgresql
```

#### Linux (Ubuntu)
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

## 🚀 Установка Tether Messenger

### 1. Клонирование репозитория
```bash
git clone https://github.com/M1ller13/tether-messenger.git
cd tether-messenger
```

### 2. Настройка базы данных

#### Создание базы данных
```bash
# Подключитесь к PostgreSQL
sudo -u postgres psql

# Создайте базу данных
CREATE DATABASE tether_messenger;

# Создайте пользователя (опционально)
CREATE USER tether_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE tether_messenger TO tether_user;

# Выйдите из psql
\q
```

#### Настройка конфигурации
Отредактируйте файл `server/config/config.go`:

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
    DBUser:     "postgres",        // или tether_user
    DBPassword: "your_password",   // ваш пароль
    DBName:     "tether_messenger",
    JWTSecret:  "your_super_secret_jwt_key_change_this_in_production",
    ServerPort: "8081",
}
```

### 3. Установка зависимостей Backend
```bash
cd server
go mod download
```

### 4. Установка зависимостей Frontend
```bash
cd ../client
npm install
```

## 🏃‍♂️ Запуск приложения

### 1. Запуск Backend сервера
```bash
cd server
go run main.go
```

Вы должны увидеть:
```
Database connected successfully
Server starting on port 8081
```

### 2. Запуск Frontend (в новом терминале)
```bash
cd client
npm run dev
```

Вы должны увидеть:
```
VITE v5.4.19  ready in 1146 ms
➜  Local:   http://localhost:3000/
```

### 3. Проверка работы
Откройте браузер и перейдите на `http://localhost:3000`

## 🔧 Конфигурация для разработки

### Переменные окружения (опционально)
Создайте файл `.env` в корне проекта:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=tether_messenger

# JWT
JWT_SECRET=your_super_secret_jwt_key

# Server
SERVER_PORT=8081
```

### Настройка IDE

#### VS Code
Рекомендуемые расширения:
- Go
- TypeScript and JavaScript Language Features
- Tailwind CSS IntelliSense
- PostgreSQL

#### GoLand/IntelliJ
- Включите поддержку Go
- Настройте PostgreSQL подключение

## 🐛 Устранение неполадок

### Проблемы с базой данных
```bash
# Проверьте статус PostgreSQL
sudo systemctl status postgresql

# Перезапустите сервис
sudo systemctl restart postgresql

# Проверьте подключение
psql -h localhost -U postgres -d tether_messenger
```

### Проблемы с портами
```bash
# Проверьте занятые порты
netstat -tulpn | grep :8081
netstat -tulpn | grep :3000

# Убейте процессы если нужно
sudo kill -9 <PID>
```

### Проблемы с зависимостями
```bash
# Очистите кэш Go
go clean -modcache

# Переустановите node_modules
rm -rf node_modules package-lock.json
npm install
```

## 📝 Следующие шаги

После успешной установки:
1. [Руководство пользователя](User-Guide.md) - как использовать приложение
2. [API документация](API-Documentation.md) - для разработчиков
3. [Развертывание](Deployment.md) - для production

## 🤝 Получение помощи

Если у вас возникли проблемы:
1. Проверьте [FAQ](FAQ.md)
2. Создайте [Issue](https://github.com/M1ller13/tether-messenger/issues)
3. Обратитесь в [Discussions](https://github.com/M1ller13/tether-messenger/discussions) 