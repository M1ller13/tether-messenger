# Архитектура системы

Tether Messenger построен с использованием современной архитектуры, разделяющей frontend и backend для обеспечения масштабируемости и производительности.

## 🏗 Общая архитектура

```
┌─────────────────┐    HTTP/WebSocket    ┌─────────────────┐
│   Frontend      │ ◄──────────────────► │    Backend      │
│   (React)       │                      │     (Go)        │
└─────────────────┘                      └─────────────────┘
                                                │
                                                │ SQL
                                                ▼
                                        ┌─────────────────┐
                                        │   PostgreSQL    │
                                        │   Database      │
                                        └─────────────────┘
```

## 🎯 Принципы архитектуры

### 1. Разделение ответственности
- **Frontend**: Пользовательский интерфейс и клиентская логика
- **Backend**: Бизнес-логика, API и работа с данными
- **Database**: Хранение и управление данными

### 2. RESTful API
- Стандартные HTTP методы (GET, POST, PUT, DELETE)
- JSON для обмена данными
- Единообразные ответы API

### 3. Безопасность
- JWT токены для аутентификации
- Валидация всех входных данных
- Защита от SQL-инъекций

## 🖥 Frontend архитектура

### Структура проекта
```
client/
├── src/
│   ├── components/     # React компоненты
│   ├── pages/         # Страницы приложения
│   ├── hooks/         # Кастомные React хуки
│   ├── api/           # API клиент
│   ├── types/         # TypeScript типы
│   ├── utils/         # Утилиты
│   └── ws/            # WebSocket клиент
├── public/            # Статические файлы
└── package.json       # Зависимости
```

### Компонентная архитектура

#### Иерархия компонентов
```
App
├── Router
│   ├── LoginPage
│   ├── RegisterPage
│   └── ChatPage
│       ├── ChatList
│       ├── MessageList
│       └── MessageInput
└── LoadingSpinner
```

#### Основные компоненты

**App.tsx** - Корневой компонент
```tsx
const App: React.FC = () => {
  const { user, loading } = useAuth();
  
  if (loading) return <LoadingSpinner />;
  
  return (
    <Router>
      {user ? <ChatPage /> : <LoginPage />}
    </Router>
  );
};
```

**ChatPage.tsx** - Главная страница чатов
```tsx
const ChatPage: React.FC = () => {
  const [selectedChat, setSelectedChat] = useState<Chat | null>(null);
  const [chats, setChats] = useState<Chat[]>([]);
  
  return (
    <div className="flex h-screen">
      <ChatList 
        chats={chats}
        onSelectChat={setSelectedChat}
      />
      {selectedChat && (
        <div className="flex-1 flex flex-col">
          <MessageList chatId={selectedChat.id} />
          <MessageInput onSend={handleSendMessage} />
        </div>
      )}
    </div>
  );
};
```

### Управление состоянием

#### Локальное состояние
- Используется `useState` для простого состояния компонентов
- `useReducer` для сложной логики

#### Глобальное состояние
- JWT токен хранится в `localStorage`
- Данные пользователя кэшируются в памяти
- Состояние чатов синхронизируется с сервером

### API интеграция

#### API клиент
```typescript
// api/client.ts
class ApiClient {
  private baseUrl: string;
  private token: string | null;
  
  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
    this.token = localStorage.getItem('token');
  }
  
  async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options?.headers,
    };
    
    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }
    
    const response = await fetch(url, {
      ...options,
      headers,
    });
    
    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }
    
    return response.json();
  }
}
```

## 🖥 Backend архитектура

### Структура проекта
```
server/
├── main.go              # Точка входа
├── config/
│   └── config.go        # Конфигурация
├── database/
│   └── database.go      # Подключение к БД
├── models/
│   ├── user.go          # Модель пользователя
│   ├── chat.go          # Модель чата
│   ├── message.go       # Модель сообщения
│   └── verification_code.go
├── handlers/
│   ├── auth.go          # Обработчики аутентификации
│   └── chat.go          # Обработчики чатов
├── middleware/
│   └── auth.go          # Middleware аутентификации
├── routes/
│   └── routes.go        # Маршруты API
├── utils/
│   ├── jwt.go           # JWT утилиты
│   └── password.go      # Утилиты паролей
└── ws/
    └── websocket.go     # WebSocket обработчик
```

### Слои архитектуры

#### 1. HTTP слой (Fiber)
```go
func main() {
    app := fiber.New(fiber.Config{
        ErrorHandler: customErrorHandler,
    })
    
    // Middleware
    app.Use(cors.New())
    app.Use(logger.New())
    
    // Routes
    routes.SetupRoutes(app)
    
    app.Listen(":8081")
}
```

#### 2. Обработчики (Handlers)
```go
type AuthHandler struct {
    userService *services.UserService
    jwtService  *utils.JWTService
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "error":   "Invalid request body",
        })
    }
    
    user, err := h.userService.CreateUser(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "error":   err.Error(),
        })
    }
    
    return c.Status(201).JSON(fiber.Map{
        "success": true,
        "data":    user,
    })
}
```

#### 3. Сервисы (Services)
```go
type UserService struct {
    db *gorm.DB
}

func (s *UserService) CreateUser(req RegisterRequest) (*User, error) {
    // Валидация
    if err := s.validateUser(req); err != nil {
        return nil, err
    }
    
    // Создание пользователя
    user := &User{
        ID:          uuid.New(),
        Phone:       req.Phone,
        DisplayName: req.DisplayName,
        Username:    generateUsername(),
        CreatedAt:   time.Now(),
    }
    
    if err := s.db.Create(user).Error; err != nil {
        return nil, err
    }
    
    return user, nil
}
```

#### 4. Модели (Models)
```go
type User struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    Phone       string    `json:"phone" gorm:"uniqueIndex;not null"`
    Username    string    `json:"username" gorm:"uniqueIndex;not null"`
    DisplayName string    `json:"display_name" gorm:"not null"`
    Bio         string    `json:"bio"`
    AvatarURL   string    `json:"avatar_url"`
    LastSeen    time.Time `json:"last_seen"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Middleware

#### Аутентификация
```go
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        if token == "" {
            return c.Status(401).JSON(fiber.Map{
                "success": false,
                "error":   "Authorization header required",
            })
        }
        
        // Убираем "Bearer " префикс
        token = strings.TrimPrefix(token, "Bearer ")
        
        claims, err := utils.ValidateJWT(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "success": false,
                "error":   "Invalid token",
            })
        }
        
        // Добавляем user_id в контекст
        c.Locals("user_id", claims.UserID)
        return c.Next()
    }
}
```

## 🗄 База данных

### Схема базы данных

```sql
-- Пользователи
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

-- Чаты
CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user1_id, user2_id)
);

-- Сообщения
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Коды верификации
CREATE TABLE verification_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) NOT NULL,
    code VARCHAR(6) NOT NULL,
    display_name VARCHAR(100),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Индексы
```sql
-- Индексы для производительности
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_chats_user1 ON chats(user1_id);
CREATE INDEX idx_chats_user2 ON chats(user2_id);
CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_verification_codes_phone ON verification_codes(phone);
CREATE INDEX idx_verification_codes_expires_at ON verification_codes(expires_at);
```

## 🔄 Потоки данных

### Регистрация пользователя
```
1. Frontend → POST /api/auth/register
2. Backend → Валидация данных
3. Backend → Создание кода верификации
4. Backend → Сохранение в БД
5. Backend → Возврат успеха
6. Frontend → Запрос кода
7. Frontend → POST /api/auth/verify-code
8. Backend → Проверка кода
9. Backend → Создание пользователя
10. Backend → Генерация JWT
11. Backend → Возврат токена
12. Frontend → Сохранение токена
```

### Отправка сообщения
```
1. Frontend → POST /api/messages
2. Backend → Проверка JWT
3. Backend → Валидация данных
4. Backend → Сохранение в БД
5. Backend → Возврат сообщения
6. Frontend → Обновление UI
```

### Получение чатов
```
1. Frontend → GET /api/chats
2. Backend → Проверка JWT
3. Backend → Запрос к БД
4. Backend → Форматирование данных
5. Backend → Возврат чатов
6. Frontend → Отображение списка
```

## 🔒 Безопасность

### Аутентификация
- JWT токены с временем жизни
- Автоматическая очистка невалидных токенов
- Безопасное хранение в localStorage

### Валидация данных
- Проверка всех входных параметров
- Санитизация данных
- Защита от SQL-инъекций через GORM

### CORS
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:3000",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowMethods: "GET, POST, PUT, DELETE",
}))
```

## 📈 Масштабируемость

### Горизонтальное масштабирование
- Stateless backend (без сессий)
- База данных может быть реплицирована
- Load balancer для распределения нагрузки

### Кэширование
- Redis для кэширования (планируется)
- Кэширование профилей пользователей
- Кэширование списков чатов

### Мониторинг
- Логирование всех запросов
- Метрики производительности
- Health checks

## 🔮 Будущие улучшения

### WebSocket интеграция
```go
// ws/websocket.go
func HandleWebSocket(c *fiber.Ctx) error {
    return websocket.New(func(conn *websocket.Conn) {
        // Обработка real-time сообщений
        for {
            messageType, message, err := conn.ReadMessage()
            if err != nil {
                break
            }
            
            // Обработка сообщения
            handleMessage(conn, messageType, message)
        }
    })(c)
}
```

### Микросервисная архитектура
- Разделение на отдельные сервисы
- API Gateway
- Service Discovery
- Message Queue

---

Эта архитектура обеспечивает надежность, масштабируемость и простоту разработки Tether Messenger. 