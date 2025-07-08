# Развертывание

Руководство по развертыванию Tether Messenger на production серверах.

## 🚀 Обзор

### Варианты развертывания

1. **VPS/Cloud сервер** - Полный контроль, настройка с нуля
2. **Docker** - Контейнеризация (планируется)
3. **Kubernetes** - Оркестрация (планируется)
4. **PaaS** - Heroku, Railway, Render

### Рекомендуемые требования

- **CPU**: 2+ ядра
- **RAM**: 4GB+
- **Диск**: 20GB+ SSD
- **ОС**: Ubuntu 20.04+ / CentOS 8+ / Debian 11+

## 🖥 Подготовка сервера

### 1. Обновление системы

```bash
# Ubuntu/Debian
sudo apt update && sudo apt upgrade -y

# CentOS/RHEL
sudo yum update -y
```

### 2. Установка необходимого ПО

```bash
# Установка Git
sudo apt install git -y

# Установка curl
sudo apt install curl -y

# Установка build tools
sudo apt install build-essential -y
```

### 3. Установка Go

```bash
# Скачивание Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# Распаковка
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Добавление в PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Проверка установки
go version
```

### 4. Установка Node.js

```bash
# Добавление NodeSource репозитория
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -

# Установка Node.js
sudo apt-get install -y nodejs

# Проверка установки
node --version
npm --version
```

### 5. Установка PostgreSQL

```bash
# Добавление репозитория PostgreSQL
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt update

# Установка PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Запуск сервиса
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

## 🗄 Настройка базы данных

### 1. Создание пользователя и базы данных

```bash
# Переключение на пользователя postgres
sudo -u postgres psql

# Создание пользователя
CREATE USER tether_user WITH PASSWORD 'your_secure_password';

# Создание базы данных
CREATE DATABASE tether_messenger OWNER tether_user;

# Предоставление прав
GRANT ALL PRIVILEGES ON DATABASE tether_messenger TO tether_user;

# Выход из psql
\q
```

### 2. Настройка PostgreSQL

```bash
# Редактирование конфигурации
sudo nano /etc/postgresql/14/main/postgresql.conf
```

Добавьте/измените следующие параметры:

```conf
# Производительность
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 4MB
maintenance_work_mem = 64MB

# Логирование
log_destination = 'stderr'
logging_collector = on
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 100MB

# Соединения
max_connections = 100
```

### 3. Настройка аутентификации

```bash
# Редактирование pg_hba.conf
sudo nano /etc/postgresql/14/main/pg_hba.conf
```

Добавьте строку для локальных соединений:

```conf
# IPv4 local connections:
host    all             all             127.0.0.1/32            md5
```

### 4. Перезапуск PostgreSQL

```bash
sudo systemctl restart postgresql
```

## 🔧 Настройка приложения

### 1. Клонирование репозитория

```bash
# Клонирование
git clone https://github.com/M1ller13/tether-messenger.git
cd tether-messenger

# Переключение на production ветку (если есть)
git checkout production
```

### 2. Настройка Backend

```bash
cd server

# Установка зависимостей
go mod download

# Создание конфигурационного файла
nano config/config.go
```

Обновите конфигурацию:

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
    DBUser:     "tether_user",
    DBPassword: "your_secure_password",
    DBName:     "tether_messenger",
    JWTSecret:  "your_super_secure_jwt_secret_key_change_this_in_production",
    ServerPort: "8081",
}
```

### 3. Настройка Frontend

```bash
cd ../client

# Установка зависимостей
npm install

# Создание production сборки
npm run build
```

### 4. Настройка переменных окружения

```bash
# Создание .env файла
nano .env
```

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=tether_user
DB_PASSWORD=your_secure_password
DB_NAME=tether_messenger

# JWT
JWT_SECRET=your_super_secure_jwt_secret_key_change_this_in_production

# Server
SERVER_PORT=8081
NODE_ENV=production
```

## 🚀 Запуск приложения

### 1. Создание systemd сервиса для Backend

```bash
sudo nano /etc/systemd/system/tether-backend.service
```

```ini
[Unit]
Description=Tether Messenger Backend
After=network.target postgresql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/tether-messenger/server
ExecStart=/usr/local/go/bin/go run main.go
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

### 2. Создание systemd сервиса для Frontend

```bash
sudo nano /etc/systemd/system/tether-frontend.service
```

```ini
[Unit]
Description=Tether Messenger Frontend
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/tether-messenger/client
ExecStart=/usr/bin/npm run preview
Restart=always
RestartSec=5
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
```

### 3. Запуск сервисов

```bash
# Перезагрузка systemd
sudo systemctl daemon-reload

# Включение автозапуска
sudo systemctl enable tether-backend
sudo systemctl enable tether-frontend

# Запуск сервисов
sudo systemctl start tether-backend
sudo systemctl start tether-frontend

# Проверка статуса
sudo systemctl status tether-backend
sudo systemctl status tether-frontend
```

## 🌐 Настройка веб-сервера

### 1. Установка Nginx

```bash
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

### 2. Настройка Nginx

```bash
sudo nano /etc/nginx/sites-available/tether-messenger
```

```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;

    # Frontend (статичные файлы)
    location / {
        root /home/ubuntu/tether-messenger/client/dist;
        try_files $uri $uri/ /index.html;
        
        # Кэширование статических файлов
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # WebSocket (для будущих версий)
    location /ws/ {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. Активация конфигурации

```bash
# Создание символической ссылки
sudo ln -s /etc/nginx/sites-available/tether-messenger /etc/nginx/sites-enabled/

# Удаление дефолтной конфигурации
sudo rm /etc/nginx/sites-enabled/default

# Проверка конфигурации
sudo nginx -t

# Перезапуск Nginx
sudo systemctl restart nginx
```

## 🔒 Настройка SSL

### 1. Установка Certbot

```bash
sudo apt install certbot python3-certbot-nginx -y
```

### 2. Получение SSL сертификата

```bash
sudo certbot --nginx -d your-domain.com -d www.your-domain.com
```

### 3. Автоматическое обновление

```bash
# Проверка автоматического обновления
sudo certbot renew --dry-run

# Добавление в cron
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

## 🔥 Настройка файрвола

### 1. Установка UFW

```bash
sudo apt install ufw -y
```

### 2. Настройка правил

```bash
# Разрешение SSH
sudo ufw allow ssh

# Разрешение HTTP и HTTPS
sudo ufw allow 80
sudo ufw allow 443

# Включение файрвола
sudo ufw enable

# Проверка статуса
sudo ufw status
```

## 📊 Мониторинг

### 1. Установка monitoring tools

```bash
# Установка htop
sudo apt install htop -y

# Установка netdata (опционально)
bash <(curl -Ss https://my-netdata.io/kickstart.sh)
```

### 2. Настройка логирования

```bash
# Создание директории для логов
sudo mkdir -p /var/log/tether-messenger

# Настройка logrotate
sudo nano /etc/logrotate.d/tether-messenger
```

```conf
/var/log/tether-messenger/*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 ubuntu ubuntu
    postrotate
        systemctl reload tether-backend
        systemctl reload tether-frontend
    endscript
}
```

### 3. Health checks

```bash
# Создание скрипта проверки здоровья
nano /home/ubuntu/health-check.sh
```

```bash
#!/bin/bash

# Проверка Backend
if ! curl -f http://localhost:8081/api/health; then
    echo "Backend is down!"
    systemctl restart tether-backend
fi

# Проверка Frontend
if ! curl -f http://localhost:3000; then
    echo "Frontend is down!"
    systemctl restart tether-frontend
fi

# Проверка PostgreSQL
if ! sudo -u postgres pg_isready; then
    echo "PostgreSQL is down!"
    systemctl restart postgresql
fi
```

```bash
# Сделать скрипт исполняемым
chmod +x /home/ubuntu/health-check.sh

# Добавить в cron (каждые 5 минут)
echo "*/5 * * * * /home/ubuntu/health-check.sh" | crontab -
```

## 🔄 Резервное копирование

### 1. Создание скрипта бэкапа

```bash
nano /home/ubuntu/backup.sh
```

```bash
#!/bin/bash

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/home/ubuntu/backups"
DB_NAME="tether_messenger"

# Создание директории для бэкапов
mkdir -p $BACKUP_DIR

# Бэкап базы данных
pg_dump -h localhost -U tether_user -d $DB_NAME > $BACKUP_DIR/db_backup_$DATE.sql

# Бэкап конфигурационных файлов
tar -czf $BACKUP_DIR/config_backup_$DATE.tar.gz \
    /home/ubuntu/tether-messenger/server/config \
    /home/ubuntu/tether-messenger/client/.env \
    /etc/nginx/sites-available/tether-messenger

# Сжатие бэкапа БД
gzip $BACKUP_DIR/db_backup_$DATE.sql

# Удаление старых бэкапов (старше 30 дней)
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete

echo "Backup completed: $DATE"
```

### 2. Настройка автоматических бэкапов

```bash
# Сделать скрипт исполняемым
chmod +x /home/ubuntu/backup.sh

# Добавить в cron (ежедневно в 2:00)
echo "0 2 * * * /home/ubuntu/backup.sh" | crontab -
```

## 🔧 Обновления

### 1. Создание скрипта обновления

```bash
nano /home/ubuntu/update.sh
```

```bash
#!/bin/bash

cd /home/ubuntu/tether-messenger

# Остановка сервисов
sudo systemctl stop tether-backend
sudo systemctl stop tether-frontend

# Получение обновлений
git pull origin main

# Обновление Backend
cd server
go mod download
go build -o tether-server main.go

# Обновление Frontend
cd ../client
npm install
npm run build

# Запуск сервисов
sudo systemctl start tether-backend
sudo systemctl start tether-frontend

echo "Update completed successfully!"
```

### 2. Безопасное обновление

```bash
# Создание бэкапа перед обновлением
/home/ubuntu/backup.sh

# Выполнение обновления
/home/ubuntu/update.sh

# Проверка работоспособности
curl -f http://localhost:8081/api/health
```

## 🚨 Troubleshooting

### Частые проблемы

#### 1. Сервис не запускается
```bash
# Проверка логов
sudo journalctl -u tether-backend -f
sudo journalctl -u tether-frontend -f

# Проверка портов
sudo netstat -tulpn | grep :8081
sudo netstat -tulpn | grep :3000
```

#### 2. Проблемы с базой данных
```bash
# Проверка статуса PostgreSQL
sudo systemctl status postgresql

# Проверка подключения
sudo -u postgres psql -d tether_messenger -c "SELECT version();"

# Проверка логов PostgreSQL
sudo tail -f /var/log/postgresql/postgresql-14-main.log
```

#### 3. Проблемы с Nginx
```bash
# Проверка конфигурации
sudo nginx -t

# Проверка логов
sudo tail -f /var/log/nginx/error.log
sudo tail -f /var/log/nginx/access.log
```

### Полезные команды

```bash
# Перезапуск всех сервисов
sudo systemctl restart tether-backend tether-frontend nginx postgresql

# Проверка использования ресурсов
htop
df -h
free -h

# Проверка сетевых соединений
ss -tulpn
```

## 🔮 Будущие улучшения

### Docker развертывание
- Создание Dockerfile для Backend и Frontend
- Docker Compose для оркестрации
- Kubernetes для масштабирования

### CI/CD
- GitHub Actions для автоматического развертывания
- Тестирование перед деплоем
- Blue-green deployment

### Мониторинг
- Prometheus + Grafana
- ELK Stack для логов
- APM инструменты

---

Это руководство поможет вам успешно развернуть Tether Messenger на production сервере. 