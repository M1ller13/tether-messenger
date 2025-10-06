# 🎬 Демонстрация Tether Messenger

## 📸 Скриншоты интерфейса

### 1. Страница входа
```
┌─────────────────────────────────────┐
│        Вход в Tether Messenger      │
├─────────────────────────────────────┤
│  Email: [user@example.com        ] │
│  Пароль: [••••••••••••••••••••••] │
│                                     │
│  [           Войти           ]      │
│                                     │
│  Нет аккаунта? Зарегистрироваться   │
│  Забыли пароль?                     │
└─────────────────────────────────────┘
```

### 2. Страница регистрации
```
┌─────────────────────────────────────┐
│              Регистрация            │
├─────────────────────────────────────┤
│  Email: [user@example.com        ] │
│  Отображаемое имя: [Иван Иванов  ] │
│  Имя пользователя: [ivan123      ] │
│  Пароль: [••••••••••••••••••••••] │
│  Подтвердите пароль: [••••••••••] │
│                                     │
│  [        Зарегистрироваться   ]    │
│                                     │
│  Уже есть аккаунт? Войти            │
└─────────────────────────────────────┘
```

### 3. Успешная регистрация
```
┌─────────────────────────────────────┐
│        ✅ Регистрация успешна!      │
├─────────────────────────────────────┤
│                                     │
│  Мы отправили письмо с подтверждением│
│  на ваш email. Пожалуйста, проверьте│
│  почту и перейдите по ссылке для    │
│  активации аккаунта.                │
│                                     │
│  [      Перейти к входу      ]      │
└─────────────────────────────────────┘
```

### 4. Подтверждение email
```
┌─────────────────────────────────────┐
│        ✅ Email подтвержден!        │
├─────────────────────────────────────┤
│                                     │
│  Ваш email успешно подтвержден.     │
│  Вы будете перенаправлены на        │
│  страницу входа.                    │
│                                     │
└─────────────────────────────────────┘
```

### 5. Сброс пароля
```
┌─────────────────────────────────────┐
│           Сброс пароля              │
├─────────────────────────────────────┤
│                                     │
│  Введите ваш email, и мы отправим   │
│  ссылку для сброса пароля.          │
│                                     │
│  Email: [user@example.com        ] │
│                                     │
│  [        Отправить ссылку     ]    │
│                                     │
│  [      Вернуться к входу      ]    │
└─────────────────────────────────────┘
```

---

## 🎯 Текущий функционал

### ✅ Аутентификация
- **Email регистрация** с подтверждением
- **Вход по email/паролю** с JWT токенами
- **Refresh tokens** для безопасности
- **Сброс пароля** через email
- **Валидация форм** на frontend

### ✅ API для досок задач
- **Создание досок** (личные, командные, CRM)
- **Управление колонками** (To Do, In Progress, Done)
- **Карточки с CRM полями**:
  - Lead name, contact email, company
  - Value, priority, status
  - Назначение ответственных
- **Права доступа** и безопасность

### ✅ База данных
- **PostgreSQL** с UUID первичными ключами
- **Модели**: User, Board, Column, Card, Workspace
- **Связи** между сущностями
- **Миграции** через GORM

---

## 🚀 Как запустить

### Вариант 1: Локальный запуск
```bash
# 1. Запустите PostgreSQL
# 2. Создайте базу данных tether_messenger

# 3. Запустите backend
cd server
go run main.go

# 4. Запустите frontend
cd client
npm run dev
```

### Вариант 2: Docker (когда Docker Desktop запустится)
```bash
# Development окружение
docker-compose -f docker-compose.dev.yml up -d

# Полный запуск
docker-compose up -d --build
```

---

## 🧪 Тестирование API

### 1. Регистрация пользователя
```bash
curl -X POST http://localhost:8081/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "display_name": "Test User"
  }'
```

### 2. Вход в систему
```bash
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Создание доски
```bash
curl -X POST http://localhost:8081/api/boards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "My Project Board",
    "description": "Board for project management",
    "type": "personal",
    "color": "#3B82F6"
  }'
```

### 4. Создание карточки
```bash
curl -X POST http://localhost:8081/api/cards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "title": "New Lead",
    "description": "Potential customer inquiry",
    "column_id": "COLUMN_UUID",
    "lead_name": "John Doe",
    "contact_email": "john@company.com",
    "company": "Acme Corp",
    "value": 50000,
    "priority": "high",
    "status": "new"
  }'
```

---

## 📊 Структура данных

### User (Пользователь)
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "username": "user123",
  "display_name": "John Doe",
  "bio": "Software Developer",
  "avatar_url": "/uploads/avatar.jpg",
  "email_verified": true,
  "created_at": "2025-10-06T06:00:00Z"
}
```

### Board (Доска)
```json
{
  "id": "uuid",
  "name": "Project Board",
  "description": "Main project board",
  "type": "personal",
  "owner_id": "uuid",
  "workspace_id": null,
  "is_public": false,
  "color": "#3B82F6",
  "columns": [...],
  "created_at": "2025-10-06T06:00:00Z"
}
```

### Card (Карточка)
```json
{
  "id": "uuid",
  "title": "New Feature",
  "description": "Implement user authentication",
  "position": 0,
  "column_id": "uuid",
  "assignee_id": "uuid",
  "created_by_id": "uuid",
  "due_date": "2025-10-15T00:00:00Z",
  "lead_name": "Jane Smith",
  "contact_email": "jane@company.com",
  "company": "Tech Corp",
  "value": 25000,
  "priority": "medium",
  "status": "in_progress"
}
```

---

## 🎨 UI/UX Особенности

### Дизайн
- **Современный интерфейс** с Tailwind CSS
- **Адаптивный дизайн** для всех устройств
- **Темная/светлая тема** (планируется)
- **Анимации** и переходы

### Пользовательский опыт
- **Интуитивная навигация**
- **Валидация форм** в реальном времени
- **Обратная связь** для всех действий
- **Обработка ошибок** с понятными сообщениями

---

## 🔮 Следующие шаги

### 1. UI для досок задач
- Kanban доска с drag & drop
- Модальные окна для редактирования
- Фильтры и поиск

### 2. Рабочие пространства
- Создание команд
- Управление ролями
- Приглашения пользователей

### 3. Уведомления
- Email уведомления
- Push уведомления
- In-app уведомления

### 4. Интеграции
- API для внешних систем
- Webhooks
- Импорт/экспорт данных

---

## 🎉 Заключение

Tether Messenger успешно развивается! Основная функциональность аутентификации и API для досок задач реализована. Проект готов к дальнейшей разработке UI компонентов и расширению функциональности.

**Текущий статус**: MVP готов, можно тестировать основную функциональность через API и базовый UI для аутентификации.
