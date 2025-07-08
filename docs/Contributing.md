# Вклад в проект

Спасибо за интерес к Tether Messenger! Мы приветствуем вклад от сообщества.

## 🤝 Как внести вклад

### Типы вкладов

Мы приветствуем различные типы вкладов:

- 🐛 **Исправление багов**
- ✨ **Новые функции**
- 📚 **Улучшение документации**
- 🎨 **Улучшение UI/UX**
- ⚡ **Оптимизация производительности**
- 🧪 **Тесты**

### Процесс вклада

#### 1. Подготовка

1. **Форкните репозиторий**
   ```bash
   git clone https://github.com/YOUR_USERNAME/tether-messenger.git
   cd tether-messenger
   ```

2. **Создайте ветку**
   ```bash
   git checkout -b feature/your-feature-name
   # или
   git checkout -b fix/your-bug-fix
   ```

3. **Настройте окружение**
   ```bash
   # Backend
   cd server
   go mod download
   
   # Frontend
   cd ../client
   npm install
   ```

#### 2. Разработка

1. **Следуйте стандартам кода**
   - Используйте ESLint и Prettier для frontend
   - Следуйте Go conventions для backend
   - Пишите понятные комментарии

2. **Тестируйте изменения**
   - Убедитесь, что все тесты проходят
   - Протестируйте функциональность вручную
   - Проверьте на разных браузерах

3. **Обновляйте документацию**
   - Обновите README если нужно
   - Добавьте комментарии к новым функциям
   - Обновите API документацию

#### 3. Коммиты

Используйте понятные сообщения коммитов:

```bash
# Хорошо
git commit -m "feat: add user search functionality"
git commit -m "fix: resolve authentication token issue"
git commit -m "docs: update API documentation"

# Плохо
git commit -m "fix stuff"
git commit -m "update"
```

**Префиксы коммитов:**
- `feat:` - новая функция
- `fix:` - исправление бага
- `docs:` - изменения в документации
- `style:` - форматирование кода
- `refactor:` - рефакторинг
- `test:` - добавление тестов
- `chore:` - обновление зависимостей

#### 4. Push и Pull Request

```bash
git push origin feature/your-feature-name
```

Создайте Pull Request на GitHub с описанием изменений.

## 📋 Стандарты кода

### Backend (Go)

#### Структура кода
```go
// Хорошо
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "error":   "Invalid request body",
        })
    }
    
    // Валидация
    if req.DisplayName == "" {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "error":   "Display name is required",
        })
    }
    
    // Логика
    user, err := h.userService.CreateUser(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "error":   "Failed to create user",
        })
    }
    
    return c.Status(201).JSON(fiber.Map{
        "success": true,
        "data":    user,
    })
}
```

#### Обработка ошибок
```go
// Всегда возвращайте структурированные ошибки
return c.Status(400).JSON(fiber.Map{
    "success": false,
    "error":   "Описание ошибки",
})
```

### Frontend (React/TypeScript)

#### Компоненты
```tsx
// Хорошо
interface MessageInputProps {
  onSend: (message: string) => void;
  disabled?: boolean;
}

export const MessageInput: React.FC<MessageInputProps> = ({ 
  onSend, 
  disabled = false 
}) => {
  const [message, setMessage] = useState('');
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (message.trim()) {
      onSend(message.trim());
      setMessage('');
    }
  };
  
  return (
    <form onSubmit={handleSubmit} className="flex gap-2">
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        disabled={disabled}
        className="flex-1 px-3 py-2 border rounded-lg"
        placeholder="Введите сообщение..."
      />
      <button
        type="submit"
        disabled={disabled || !message.trim()}
        className="px-4 py-2 bg-blue-500 text-white rounded-lg disabled:opacity-50"
      >
        Отправить
      </button>
    </form>
  );
};
```

#### Хуки
```tsx
// Хорошо
export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      validateToken(token);
    } else {
      setLoading(false);
    }
  }, []);
  
  const validateToken = async (token: string) => {
    try {
      const response = await api.get('/profile', { token });
      setUser(response.data);
    } catch (error) {
      localStorage.removeItem('token');
    } finally {
      setLoading(false);
    }
  };
  
  return { user, loading };
};
```

## 🧪 Тестирование

### Backend тесты
```go
func TestCreateUser(t *testing.T) {
    app := fiber.New()
    handler := &Handler{}
    
    req := httptest.NewRequest("POST", "/api/auth/register", 
        strings.NewReader(`{"display_name":"Test","phone":"+79991234567"}`))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

### Frontend тесты
```tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { MessageInput } from './MessageInput';

test('sends message when form is submitted', () => {
  const mockOnSend = jest.fn();
  render(<MessageInput onSend={mockOnSend} />);
  
  const input = screen.getByPlaceholderText('Введите сообщение...');
  const button = screen.getByText('Отправить');
  
  fireEvent.change(input, { target: { value: 'Hello' } });
  fireEvent.click(button);
  
  expect(mockOnSend).toHaveBeenCalledWith('Hello');
});
```

## 📝 Документация

### Комментарии к коду
```go
// CreateUser создает нового пользователя в системе
// Возвращает созданного пользователя или ошибку
func (s *UserService) CreateUser(req CreateUserRequest) (*User, error) {
    // Проверяем уникальность телефона
    if s.userExists(req.Phone) {
        return nil, errors.New("user already exists")
    }
    
    // Создаем пользователя
    user := &User{
        ID:          uuid.New(),
        Phone:       req.Phone,
        DisplayName: req.DisplayName,
        Username:    generateUsername(),
        CreatedAt:   time.Now(),
    }
    
    return user, s.db.Create(user).Error
}
```

### README обновления
При добавлении новых функций обновляйте:
- Основной README.md
- Соответствующие страницы Wiki
- API документацию

## 🔍 Code Review

### Что проверяется

1. **Функциональность**
   - Код работает как ожидается
   - Обрабатываются edge cases
   - Нет регрессий

2. **Качество кода**
   - Следует стандартам проекта
   - Понятный и читаемый код
   - Правильная обработка ошибок

3. **Безопасность**
   - Нет уязвимостей
   - Правильная валидация данных
   - Безопасная работа с БД

4. **Производительность**
   - Нет утечек памяти
   - Эффективные алгоритмы
   - Оптимальные запросы к БД

### Процесс Review

1. **Создание PR**
   - Четкое описание изменений
   - Ссылки на связанные issues
   - Скриншоты для UI изменений

2. **Review комментарии**
   - Конструктивная обратная связь
   - Предложения по улучшению
   - Объяснение причин изменений

3. **Исправления**
   - Внесите исправления в ту же ветку
   - Ответьте на комментарии
   - Обновите PR

## 🚀 Запуск проекта

### Требования
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### Локальная разработка
```bash
# Backend
cd server
go run main.go

# Frontend (в новом терминале)
cd client
npm run dev
```

## 📞 Получение помощи

### Вопросы по разработке
- Создайте Issue с тегом "question"
- Используйте Discussions
- Обратитесь к документации

### Проблемы с окружением
- Проверьте [Installation.md](Installation.md)
- Создайте Issue с подробным описанием
- Укажите версии ПО и ОС

## 🎯 Приоритетные задачи

### Высокий приоритет
- [ ] WebSocket для real-time сообщений
- [ ] Групповые чаты
- [ ] Поиск по сообщениям
- [ ] Push уведомления

### Средний приоритет
- [ ] Удаление сообщений
- [ ] Редактирование сообщений
- [ ] Статус "печатает..."
- [ ] Экспорт чатов

### Низкий приоритет
- [ ] Темная тема
- [ ] Кастомные эмодзи
- [ ] Голосовые сообщения
- [ ] Видеозвонки

## 📄 Лицензия

Проект распространяется под лицензией MIT. Внося вклад, вы соглашаетесь с условиями лицензии.

---

**Спасибо за ваш вклад в Tether Messenger!** 🎉

Вместе мы создаем лучший мессенджер для всех. 