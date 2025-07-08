# База данных

Tether Messenger использует PostgreSQL как основную систему управления базами данных. Это обеспечивает надежность, производительность и поддержку сложных запросов.

## 🗄 Обзор

### Технологии
- **СУБД**: PostgreSQL 14+
- **ORM**: GORM (Go)
- **Типы данных**: UUID, JSON, TIMESTAMP
- **Индексы**: B-tree, Hash

### Преимущества PostgreSQL
- ACID транзакции
- Поддержка JSON
- UUID типы данных
- Продвинутые индексы
- Триггеры и процедуры
- Репликация

## 📊 Схема базы данных

### Таблица `users`

Основная таблица для хранения пользователей.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    bio TEXT,
    avatar_url TEXT,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Поля:**
- `id` - Уникальный идентификатор пользователя (UUID)
- `phone` - Номер телефона (уникальный)
- `username` - Имя пользователя (уникальное)
- `display_name` - Отображаемое имя
- `bio` - Биография пользователя
- `avatar_url` - URL аватара
- `last_seen` - Время последней активности
- `created_at` - Время создания аккаунта
- `updated_at` - Время последнего обновления

### Таблица `chats`

Таблица для хранения чатов между пользователями.

```sql
CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user1_id, user2_id)
);
```

**Поля:**
- `id` - Уникальный идентификатор чата
- `user1_id` - ID первого пользователя
- `user2_id` - ID второго пользователя
- `created_at` - Время создания чата

**Ограничения:**
- Уникальная пара пользователей (нельзя создать два чата между одними пользователями)
- Каскадное удаление при удалении пользователя

### Таблица `messages`

Таблица для хранения сообщений в чатах.

```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Поля:**
- `id` - Уникальный идентификатор сообщения
- `chat_id` - ID чата
- `sender_id` - ID отправителя
- `content` - Текст сообщения
- `is_read` - Статус прочтения
- `created_at` - Время отправки

### Таблица `verification_codes`

Таблица для хранения временных кодов верификации.

```sql
CREATE TABLE verification_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) NOT NULL,
    code VARCHAR(6) NOT NULL,
    display_name VARCHAR(100),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Поля:**
- `id` - Уникальный идентификатор кода
- `phone` - Номер телефона
- `code` - 6-значный код верификации
- `display_name` - Временное имя для регистрации
- `expires_at` - Время истечения кода
- `created_at` - Время создания кода

## 🔍 Индексы

### Основные индексы

```sql
-- Индексы для таблицы users
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_last_seen ON users(last_seen);

-- Индексы для таблицы chats
CREATE INDEX idx_chats_user1 ON chats(user1_id);
CREATE INDEX idx_chats_user2 ON chats(user2_id);
CREATE INDEX idx_chats_created_at ON chats(created_at);

-- Индексы для таблицы messages
CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_messages_is_read ON messages(is_read);

-- Составной индекс для быстрого поиска сообщений в чате
CREATE INDEX idx_messages_chat_created ON messages(chat_id, created_at);

-- Индексы для таблицы verification_codes
CREATE INDEX idx_verification_codes_phone ON verification_codes(phone);
CREATE INDEX idx_verification_codes_expires_at ON verification_codes(expires_at);
```

### Специальные индексы

```sql
-- Частичный индекс для непрочитанных сообщений
CREATE INDEX idx_messages_unread ON messages(chat_id, created_at) 
WHERE is_read = FALSE;

-- Индекс для поиска пользователей по имени
CREATE INDEX idx_users_display_name_gin ON users USING gin(to_tsvector('russian', display_name));
```

## 🔄 Миграции

### Создание миграций

```go
// database/migrations/001_initial_schema.go
package migrations

import (
    "gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
    // Создание таблиц
    if err := db.AutoMigrate(&User{}, &Chat{}, &Message{}, &VerificationCode{}); err != nil {
        return err
    }
    
    // Создание индексов
    if err := createIndexes(db); err != nil {
        return err
    }
    
    return nil
}

func createIndexes(db *gorm.DB) error {
    indexes := []string{
        "CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone)",
        "CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
        "CREATE INDEX IF NOT EXISTS idx_chats_user1 ON chats(user1_id)",
        "CREATE INDEX IF NOT EXISTS idx_chats_user2 ON chats(user2_id)",
        "CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id)",
        "CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at)",
    }
    
    for _, index := range indexes {
        if err := db.Exec(index).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

### Применение миграций

```go
// database/database.go
func InitDatabase() (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Применение миграций
    if err := migrations.Migrate(db); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

## 📈 Оптимизация запросов

### Примеры оптимизированных запросов

#### Получение чатов с последними сообщениями

```sql
SELECT 
    c.id,
    c.created_at,
    u1.id as user1_id,
    u1.display_name as user1_name,
    u1.avatar_url as user1_avatar,
    u1.last_seen as user1_last_seen,
    u2.id as user2_id,
    u2.display_name as user2_name,
    u2.avatar_url as user2_avatar,
    u2.last_seen as user2_last_seen,
    last_message.content as last_message_content,
    last_message.created_at as last_message_time
FROM chats c
JOIN users u1 ON c.user1_id = u1.id
JOIN users u2 ON c.user2_id = u2.id
LEFT JOIN LATERAL (
    SELECT content, created_at
    FROM messages m
    WHERE m.chat_id = c.id
    ORDER BY created_at DESC
    LIMIT 1
) last_message ON true
WHERE c.user1_id = $1 OR c.user2_id = $1
ORDER BY last_message.created_at DESC NULLS LAST;
```

#### Получение сообщений с пагинацией

```sql
SELECT 
    m.id,
    m.content,
    m.is_read,
    m.created_at,
    u.id as sender_id,
    u.display_name as sender_name,
    u.avatar_url as sender_avatar
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.chat_id = $1
ORDER BY m.created_at DESC
LIMIT $2 OFFSET $3;
```

#### Поиск пользователей

```sql
SELECT 
    id,
    display_name,
    username,
    bio,
    avatar_url,
    last_seen
FROM users
WHERE 
    display_name ILIKE $1 OR 
    username ILIKE $1
ORDER BY 
    CASE WHEN display_name ILIKE $1 THEN 1 ELSE 2 END,
    last_seen DESC
LIMIT 20;
```

### GORM оптимизации

```go
// Оптимизированное получение чатов
func (s *ChatService) GetUserChats(userID uuid.UUID) ([]Chat, error) {
    var chats []Chat
    
    err := s.db.
        Preload("User1").
        Preload("User2").
        Preload("Messages", func(db *gorm.DB) *gorm.DB {
            return db.Order("created_at DESC").Limit(1)
        }).
        Where("user1_id = ? OR user2_id = ?", userID, userID).
        Order("updated_at DESC").
        Find(&chats).Error
    
    return chats, err
}

// Оптимизированное получение сообщений
func (s *MessageService) GetChatMessages(chatID uuid.UUID, limit, offset int) ([]Message, error) {
    var messages []Message
    
    err := s.db.
        Preload("Sender").
        Where("chat_id = ?", chatID).
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&messages).Error
    
    return messages, err
}
```

## 🔒 Безопасность

### Подготовленные запросы

```go
// Использование подготовленных запросов
func (s *UserService) GetUserByPhone(phone string) (*User, error) {
    var user User
    err := s.db.Where("phone = ?", phone).First(&user).Error
    return &user, err
}
```

### Валидация данных

```go
// Валидация перед сохранением
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Проверка длины имени
    if len(u.DisplayName) < 2 || len(u.DisplayName) > 100 {
        return errors.New("display name must be between 2 and 100 characters")
    }
    
    // Проверка формата телефона
    if !isValidPhone(u.Phone) {
        return errors.New("invalid phone format")
    }
    
    return nil
}
```

### Очистка данных

```go
// Автоматическая очистка устаревших кодов
func CleanupExpiredCodes(db *gorm.DB) error {
    return db.Where("expires_at < ?", time.Now()).Delete(&VerificationCode{}).Error
}

// Периодическая очистка (запускается по cron)
func ScheduleCleanup(db *gorm.DB) {
    ticker := time.NewTicker(1 * time.Hour)
    go func() {
        for range ticker.C {
            CleanupExpiredCodes(db)
        }
    }()
}
```

## 📊 Мониторинг

### Медленные запросы

```sql
-- Включение логирования медленных запросов
ALTER SYSTEM SET log_min_duration_statement = 1000; -- 1 секунда
ALTER SYSTEM SET log_statement = 'all';
SELECT pg_reload_conf();
```

### Статистика

```sql
-- Статистика по таблицам
SELECT 
    schemaname,
    tablename,
    attname,
    n_distinct,
    correlation
FROM pg_stats
WHERE tablename IN ('users', 'chats', 'messages')
ORDER BY tablename, attname;

-- Размер таблиц
SELECT 
    table_name,
    pg_size_pretty(pg_total_relation_size(table_name)) as size
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY pg_total_relation_size(table_name) DESC;
```

## 🔄 Резервное копирование

### Автоматические бэкапы

```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
DB_NAME="tether_messenger"

# Создание бэкапа
pg_dump -h localhost -U postgres -d $DB_NAME > $BACKUP_DIR/backup_$DATE.sql

# Сжатие
gzip $BACKUP_DIR/backup_$DATE.sql

# Удаление старых бэкапов (старше 7 дней)
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
```

### Восстановление

```bash
# Восстановление из бэкапа
gunzip -c backup_20250101_120000.sql.gz | psql -h localhost -U postgres -d tether_messenger
```

## 🔮 Будущие улучшения

### Партиционирование

```sql
-- Партиционирование таблицы messages по дате
CREATE TABLE messages_2025_01 PARTITION OF messages
FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE messages_2025_02 PARTITION OF messages
FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
```

### Репликация

```sql
-- Настройка репликации
-- Primary server
ALTER SYSTEM SET wal_level = replica;
ALTER SYSTEM SET max_wal_senders = 3;
ALTER SYSTEM SET wal_keep_segments = 64;

-- Replica server
ALTER SYSTEM SET hot_standby = on;
```

### Кэширование

```go
// Интеграция с Redis для кэширования
type CacheService struct {
    redis *redis.Client
}

func (s *CacheService) GetUser(userID uuid.UUID) (*User, error) {
    // Попытка получить из кэша
    if cached, err := s.redis.Get(ctx, "user:"+userID.String()).Result(); err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    // Получение из БД
    user, err := s.userService.GetByID(userID)
    if err != nil {
        return nil, err
    }
    
    // Сохранение в кэш
    if data, err := json.Marshal(user); err == nil {
        s.redis.Set(ctx, "user:"+userID.String(), data, time.Hour)
    }
    
    return user, nil
}
```

---

Эта архитектура базы данных обеспечивает производительность, надежность и масштабируемость Tether Messenger. 