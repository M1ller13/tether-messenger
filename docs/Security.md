# Безопасность

Tether Messenger разработан с учетом современных стандартов безопасности. Этот документ описывает меры безопасности, реализованные в системе.

## 🔒 Обзор безопасности

### Принципы безопасности

1. **Defense in Depth** - Многоуровневая защита
2. **Principle of Least Privilege** - Минимальные права доступа
3. **Secure by Default** - Безопасность по умолчанию
4. **Fail Securely** - Безопасный отказ

### Угрозы и контрмеры

| Угроза | Контрмера |
|--------|-----------|
| SQL Injection | Подготовленные запросы, ORM |
| XSS | Валидация и санитизация |
| CSRF | JWT токены, CORS |
| Brute Force | Ограничение попыток |
| Session Hijacking | Secure cookies, HTTPS |
| Data Breach | Шифрование, RBAC |

## 🔐 Аутентификация

### JWT токены

```go
// Генерация JWT токена
func GenerateJWT(userID uuid.UUID) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID.String(),
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
        "iat":     time.Now().Unix(),
        "iss":     "tether-messenger",
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.JWTSecret))
}

// Валидация JWT токена
func ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.JWTSecret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return &claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

### Безопасность токенов

- **Время жизни**: 7 дней
- **Алгоритм**: HMAC-SHA256
- **Хранение**: localStorage (клиент)
- **Передача**: Authorization header

### Очистка токенов

```go
// Middleware для очистки невалидных токенов
func CleanupInvalidTokens() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        if token != "" {
            token = strings.TrimPrefix(token, "Bearer ")
            if _, err := ValidateJWT(token); err != nil {
                // Токен невалиден, но продолжаем обработку
                // Клиент должен очистить токен при получении 401
            }
        }
        return c.Next()
    }
}
```

## 🛡️ Валидация данных

### Входная валидация

```go
// Структура для валидации запросов
type RegisterRequest struct {
    DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
    Phone       string `json:"phone" validate:"required,e164"`
}

type LoginRequest struct {
    Phone string `json:"phone" validate:"required,e164"`
    Code  string `json:"code" validate:"required,len=6,numeric"`
}

// Валидация в обработчиках
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "error":   "Invalid request body",
        })
    }
    
    // Валидация с помощью validator
    if err := h.validator.Struct(req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "error":   "Validation failed",
            "details": err.Error(),
        })
    }
    
    // Дополнительная бизнес-логика валидации
    if !isValidPhone(req.Phone) {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "error":   "Invalid phone number format",
        })
    }
    
    // Создание пользователя...
}
```

### Санитизация данных

```go
// Очистка HTML тегов
func sanitizeHTML(input string) string {
    // Удаление всех HTML тегов
    re := regexp.MustCompile(`<[^>]*>`)
    return re.ReplaceAllString(input, "")
}

// Очистка специальных символов
func sanitizeInput(input string) string {
    // Экранирование специальных символов
    return html.EscapeString(input)
}

// Валидация телефона
func isValidPhone(phone string) bool {
    // Проверка формата +7XXXXXXXXXX
    re := regexp.MustCompile(`^\+7\d{10}$`)
    return re.MatchString(phone)
}
```

## 🗄️ Безопасность базы данных

### Подготовленные запросы

```go
// Использование GORM для предотвращения SQL injection
func (s *UserService) GetUserByPhone(phone string) (*User, error) {
    var user User
    err := s.db.Where("phone = ?", phone).First(&user).Error
    return &user, err
}

// Прямые SQL запросы с параметрами
func (s *UserService) GetUsersByIDs(ids []uuid.UUID) ([]User, error) {
    var users []User
    query := "SELECT * FROM users WHERE id = ANY($1)"
    err := s.db.Raw(query, pq.Array(ids)).Scan(&users).Error
    return users, err
}
```

### Права доступа

```sql
-- Создание пользователя с ограниченными правами
CREATE USER tether_app WITH PASSWORD 'secure_password';

-- Предоставление только необходимых прав
GRANT CONNECT ON DATABASE tether_messenger TO tether_app;
GRANT USAGE ON SCHEMA public TO tether_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO tether_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO tether_app;

-- Ограничение прав на системные таблицы
REVOKE ALL ON ALL TABLES IN SCHEMA information_schema FROM tether_app;
REVOKE ALL ON ALL TABLES IN SCHEMA pg_catalog FROM tether_app;
```

### Шифрование данных

```go
// Шифрование чувствительных данных
func encryptSensitiveData(data string) (string, error) {
    key := []byte(config.EncryptionKey)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    ciphertext := make([]byte, aes.BlockSize+len(data))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }
    
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))
    
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Расшифровка данных
func decryptSensitiveData(encryptedData string) (string, error) {
    key := []byte(config.EncryptionKey)
    ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return "", err
    }
    
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    if len(ciphertext) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }
    
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]
    
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)
    
    return string(ciphertext), nil
}
```

## 🌐 Web Security

### CORS настройки

```go
// Настройка CORS для безопасности
app.Use(cors.New(cors.Config{
    AllowOrigins:     "https://your-domain.com,https://www.your-domain.com",
    AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
    AllowCredentials: true,
    MaxAge:           300, // 5 минут
}))
```

### Security Headers

```go
// Middleware для security headers
func SecurityHeaders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        c.Set("X-Content-Type-Options", "nosniff")
        c.Set("X-Frame-Options", "DENY")
        c.Set("X-XSS-Protection", "1; mode=block")
        c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        return c.Next()
    }
}
```

### Rate Limiting

```go
// Ограничение количества запросов
func RateLimit() fiber.Handler {
    limiter := rate.NewLimiter(rate.Every(time.Second), 10) // 10 запросов в секунду
    
    return func(c *fiber.Ctx) error {
        if !limiter.Allow() {
            return c.Status(429).JSON(fiber.Map{
                "success": false,
                "error":   "Too many requests",
            })
        }
        return c.Next()
    }
}

// Специальное ограничение для аутентификации
func AuthRateLimit() fiber.Handler {
    limiter := rate.NewLimiter(rate.Every(time.Minute), 5) // 5 попыток в минуту
    
    return func(c *fiber.Ctx) error {
        clientIP := c.IP()
        key := fmt.Sprintf("auth:%s", clientIP)
        
        if !limiter.Allow() {
            return c.Status(429).JSON(fiber.Map{
                "success": false,
                "error":   "Too many authentication attempts",
            })
        }
        return c.Next()
    }
}
```

## 🔍 Аудит и логирование

### Логирование безопасности

```go
// Структура для логов безопасности
type SecurityLog struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"user_id"`
    Action    string    `json:"action"`
    IP        string    `json:"ip"`
    UserAgent string    `json:"user_agent"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details"`
}

// Логирование событий безопасности
func LogSecurityEvent(userID uuid.UUID, action, details string, c *fiber.Ctx) {
    log := SecurityLog{
        ID:        uuid.New(),
        UserID:    userID,
        Action:    action,
        IP:        c.IP(),
        UserAgent: c.Get("User-Agent"),
        Timestamp: time.Now(),
        Details:   details,
    }
    
    // Сохранение в базу данных
    db.Create(&log)
    
    // Логирование в файл
    logEntry := fmt.Sprintf("[SECURITY] %s - User: %s, Action: %s, IP: %s, Details: %s\n",
        log.Timestamp.Format("2006-01-02 15:04:05"),
        log.UserID,
        log.Action,
        log.IP,
        log.Details)
    
    securityLogger.Info(logEntry)
}
```

### Мониторинг подозрительной активности

```go
// Обнаружение подозрительной активности
func DetectSuspiciousActivity(userID uuid.UUID, action string, c *fiber.Ctx) bool {
    // Проверка частоты действий
    recentActions := getRecentActions(userID, time.Hour)
    if len(recentActions) > 100 {
        LogSecurityEvent(userID, "SUSPICIOUS_ACTIVITY", "Too many actions", c)
        return true
    }
    
    // Проверка необычных IP адресов
    if isUnusualIP(c.IP(), userID) {
        LogSecurityEvent(userID, "SUSPICIOUS_IP", "Unusual IP address", c)
        return true
    }
    
    // Проверка подозрительных паттернов
    if isSuspiciousPattern(action, userID) {
        LogSecurityEvent(userID, "SUSPICIOUS_PATTERN", "Suspicious action pattern", c)
        return true
    }
    
    return false
}
```

## 🔐 Шифрование

### Хеширование паролей (для будущих версий)

```go
// Хеширование паролей с bcrypt
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// Проверка пароля
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### Шифрование конфиденциальных данных

```go
// Шифрование сообщений (для будущих версий)
func EncryptMessage(message string, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    ciphertext := make([]byte, aes.BlockSize+len(message))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }
    
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(message))
    
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}
```

## 🚨 Инцидент-менеджмент

### Процедура реагирования на инциденты

1. **Обнаружение**
   - Автоматическое обнаружение через мониторинг
   - Ручное сообщение через систему тикетов

2. **Оценка**
   - Определение типа и серьезности инцидента
   - Оценка потенциального ущерба

3. **Реагирование**
   - Блокировка подозрительных аккаунтов
   - Временное отключение уязвимых функций
   - Уведомление пользователей при необходимости

4. **Восстановление**
   - Исправление уязвимостей
   - Восстановление нормальной работы
   - Обновление систем безопасности

### Контакты для инцидентов

```go
// Структура для уведомлений об инцидентах
type SecurityIncident struct {
    ID          uuid.UUID `json:"id"`
    Type        string    `json:"type"`
    Severity    string    `json:"severity"`
    Description string    `json:"description"`
    Timestamp   time.Time `json:"timestamp"`
    Status      string    `json:"status"`
}

// Уведомление о критических инцидентах
func NotifySecurityIncident(incident SecurityIncident) {
    // Отправка email администраторам
    sendSecurityAlert(incident)
    
    // Логирование в систему мониторинга
    logCriticalIncident(incident)
    
    // Уведомление через Slack/Discord (если настроено)
    sendSlackNotification(incident)
}
```

## 📋 Чек-лист безопасности

### Развертывание

- [ ] HTTPS настроен и принудительно используется
- [ ] Security headers настроены
- [ ] CORS настроен правильно
- [ ] Rate limiting включен
- [ ] Логирование безопасности настроено

### Мониторинг

- [ ] Система мониторинга настроена
- [ ] Алерты на подозрительную активность
- [ ] Регулярные проверки логов
- [ ] Мониторинг производительности

### Обновления

- [ ] Регулярные обновления зависимостей
- [ ] Сканирование уязвимостей
- [ ] Тестирование безопасности
- [ ] План обновлений

### Резервное копирование

- [ ] Автоматические бэкапы настроены
- [ ] Шифрование бэкапов
- [ ] Тестирование восстановления
- [ ] Хранение бэкапов в безопасном месте

## 🔮 Будущие улучшения

### Двухфакторная аутентификация
- TOTP (Time-based One-Time Password)
- SMS коды (уже частично реализовано)
- Backup коды

### End-to-End шифрование
- Шифрование сообщений на клиенте
- Perfect Forward Secrecy
- Проверка ключей

### Аудит безопасности
- Детальное логирование всех действий
- Анализ поведения пользователей
- Машинное обучение для обнаружения аномалий

### Соответствие стандартам
- GDPR compliance
- SOC 2 Type II
- ISO 27001

---

Безопасность Tether Messenger является приоритетом. Мы постоянно работаем над улучшением мер защиты и следуем лучшим практикам индустрии. 