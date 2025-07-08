# API Документация

Полная документация API для Tether Messenger. Все запросы должны быть отправлены на базовый URL: `http://localhost:8081/api`

## 🔐 Аутентификация

Большинство эндпоинтов требуют аутентификации через JWT токен. Добавьте заголовок:
```
Authorization: Bearer <your_jwt_token>
```

## 📱 Аутентификация

### Регистрация нового пользователя

**POST** `/auth/register`

Создает новый аккаунт с подтверждением по SMS.

#### Запрос
```json
{
  "display_name": "Иван Петров",
  "phone": "+79991234567"
}
```

#### Ответ
```json
{
  "success": true
}
```

#### Пример с curl
```bash
curl -X POST http://localhost:8081/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "display_name": "Иван Петров",
    "phone": "+79991234567"
  }'
```

### Запрос кода для входа

**POST** `/auth/request-code`

Отправляет SMS-код для входа существующего пользователя.

#### Запрос
```json
{
  "phone": "+79991234567"
}
```

#### Ответ
```json
{
  "success": true
}
```

### Подтверждение кода

**POST** `/auth/verify-code`

Подтверждает SMS-код и возвращает JWT токен.

#### Запрос
```json
{
  "phone": "+79991234567",
  "code": "123456"
}
```

#### Ответ
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "phone": "+79991234567",
      "display_name": "Иван Петров",
      "username": "user1234"
    }
  }
}
```

## 👤 Профиль пользователя

### Получение профиля

**GET** `/profile`

Возвращает данные текущего пользователя.

#### Заголовки
```
Authorization: Bearer <jwt_token>
```

#### Ответ
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "phone": "+79991234567",
    "username": "user1234",
    "display_name": "Иван Петров",
    "bio": "Привет! Я разработчик.",
    "avatar_url": "https://example.com/avatar.jpg",
    "last_seen": "2025-07-08T03:45:00Z",
    "created_at": "2025-07-08T03:30:00Z"
  }
}
```

### Обновление профиля

**PUT** `/profile`

Обновляет данные профиля пользователя.

#### Запрос
```json
{
  "display_name": "Новое имя",
  "bio": "Новое описание"
}
```

#### Ответ
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "display_name": "Новое имя",
    "bio": "Новое описание",
    "username": "user1234"
  }
}
```

### Загрузка аватара

**POST** `/profile/avatar`

Загружает изображение профиля.

#### Заголовки
```
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
```

#### Тело запроса
```
avatar: <file>
```

#### Ответ
```json
{
  "success": true,
  "avatar_url": "https://example.com/uploads/avatar_123.jpg"
}
```

## 💬 Чаты

### Получение списка чатов

**GET** `/chats`

Возвращает все чаты текущего пользователя.

#### Ответ
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "user1": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "display_name": "Иван Петров",
        "username": "user1234",
        "avatar_url": "https://example.com/avatar1.jpg",
        "last_seen": "2025-07-08T03:45:00Z"
      },
      "user2": {
        "id": "550e8400-e29b-41d4-a716-446655440002",
        "display_name": "Мария Сидорова",
        "username": "user5678",
        "avatar_url": "https://example.com/avatar2.jpg",
        "last_seen": "2025-07-08T03:40:00Z"
      },
      "created_at": "2025-07-08T03:30:00Z"
    }
  ]
}
```

### Создание чата

**POST** `/chats`

Создает новый чат с пользователем.

#### Запрос
```json
{
  "other_user_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

#### Ответ
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "user1_id": "550e8400-e29b-41d4-a716-446655440000",
    "user2_id": "550e8400-e29b-41d4-a716-446655440002",
    "created_at": "2025-07-08T03:30:00Z"
  }
}
```

### Получение сообщений чата

**GET** `/chats/{chatId}/messages`

Возвращает все сообщения в указанном чате.

#### Параметры пути
- `chatId` - UUID чата

#### Ответ
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "chat_id": "550e8400-e29b-41d4-a716-446655440001",
      "sender_id": "550e8400-e29b-41d4-a716-446655440000",
      "content": "Привет! Как дела?",
      "is_read": true,
      "created_at": "2025-07-08T03:35:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440004",
      "chat_id": "550e8400-e29b-41d4-a716-446655440001",
      "sender_id": "550e8400-e29b-41d4-a716-446655440002",
      "content": "Привет! Все хорошо, спасибо!",
      "is_read": false,
      "created_at": "2025-07-08T03:36:00Z"
    }
  ]
}
```

## 📨 Сообщения

### Отправка сообщения

**POST** `/messages`

Отправляет новое сообщение в чат.

#### Запрос
```json
{
  "chat_id": "550e8400-e29b-41d4-a716-446655440001",
  "content": "Привет! Как дела?"
}
```

#### Ответ
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "chat_id": "550e8400-e29b-41d4-a716-446655440001",
    "sender_id": "550e8400-e29b-41d4-a716-446655440000",
    "content": "Привет! Как дела?",
    "is_read": false,
    "created_at": "2025-07-08T03:35:00Z"
  }
}
```

## 🔍 Поиск пользователей

### Поиск пользователей

**GET** `/users/search?query={searchTerm}`

Ищет пользователей по имени или username.

#### Параметры запроса
- `query` - строка поиска (минимум 2 символа)

#### Ответ
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "display_name": "Мария Сидорова",
      "username": "user5678",
      "bio": "Дизайнер",
      "avatar_url": "https://example.com/avatar2.jpg",
      "last_seen": "2025-07-08T03:40:00Z"
    }
  ]
}
```

## 🏥 Проверка здоровья

### Проверка API

**GET** `/health`

Проверяет работоспособность API.

#### Ответ
```json
{
  "success": true,
  "message": "Tether Messenger API is running"
}
```

## 📊 Коды ошибок

### HTTP статус коды

- `200` - Успешный запрос
- `201` - Ресурс создан
- `400` - Неверный запрос
- `401` - Не авторизован
- `403` - Доступ запрещен
- `404` - Ресурс не найден
- `409` - Конфликт (например, пользователь уже существует)
- `500` - Внутренняя ошибка сервера

### Формат ошибок

```json
{
  "success": false,
  "error": "Описание ошибки"
}
```

## 🔧 Примеры использования

### JavaScript (Fetch API)

```javascript
// Регистрация
const register = async (displayName, phone) => {
  const response = await fetch('/api/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ display_name: displayName, phone }),
  });
  return response.json();
};

// Вход
const login = async (phone, code) => {
  const response = await fetch('/api/auth/verify-code', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ phone, code }),
  });
  return response.json();
};

// Получение чатов
const getChats = async (token) => {
  const response = await fetch('/api/chats', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  return response.json();
};

// Отправка сообщения
const sendMessage = async (token, chatId, content) => {
  const response = await fetch('/api/messages', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ chat_id: chatId, content }),
  });
  return response.json();
};
```

### Python (requests)

```python
import requests

BASE_URL = "http://localhost:8081/api"

# Регистрация
def register(display_name, phone):
    response = requests.post(f"{BASE_URL}/auth/register", json={
        "display_name": display_name,
        "phone": phone
    })
    return response.json()

# Вход
def login(phone, code):
    response = requests.post(f"{BASE_URL}/auth/verify-code", json={
        "phone": phone,
        "code": code
    })
    return response.json()

# Получение чатов
def get_chats(token):
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(f"{BASE_URL}/chats", headers=headers)
    return response.json()

# Отправка сообщения
def send_message(token, chat_id, content):
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.post(f"{BASE_URL}/messages", json={
        "chat_id": chat_id,
        "content": content
    }, headers=headers)
    return response.json()
```

## 🔄 WebSocket (будущие версии)

В будущих версиях будет добавлена поддержка WebSocket для реального времени:

```javascript
// Пример будущего WebSocket API
const ws = new WebSocket('ws://localhost:8081/ws');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  switch (data.type) {
    case 'new_message':
      // Обработка нового сообщения
      break;
    case 'user_online':
      // Пользователь в сети
      break;
    case 'message_read':
      // Сообщение прочитано
      break;
  }
};
```

## 📝 Примечания

- Все UUID должны быть в формате RFC 4122
- Даты передаются в формате ISO 8601
- Максимальная длина сообщения: 1000 символов
- Максимальный размер аватара: 5MB
- Поддерживаемые форматы изображений: JPG, PNG, GIF

---

Для получения дополнительной помощи обратитесь к [GitHub Issues](https://github.com/M1ller13/tether-messenger/issues). 