import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

const RegisterPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [username, setUsername] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const navigate = useNavigate();

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    // Validation
    if (password !== confirmPassword) {
      setError('Пароли не совпадают');
      setLoading(false);
      return;
    }

    if (password.length < 6) {
      setError('Пароль должен содержать минимум 6 символов');
      setLoading(false);
      return;
    }

    try {
      const res = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          email, 
          password, 
          display_name: displayName,
          username: username || undefined
        }),
      });
      const data = await res.json();
      if (data.success) {
        setSuccess(true);
      } else {
        setError(data.error || 'Ошибка регистрации');
      }
    } catch {
      setError('Ошибка сети');
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-100">
        <div className="bg-white p-8 rounded shadow-md w-full max-w-md text-center">
          <h2 className="text-2xl font-bold mb-4 text-green-600">Регистрация успешна!</h2>
          <p className="text-gray-600 mb-4">
            Мы отправили письмо с подтверждением на ваш email. 
            Пожалуйста, проверьте почту и перейдите по ссылке для активации аккаунта.
          </p>
          <Link to="/login" className="text-blue-600 hover:underline">
            Перейти к входу
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <form onSubmit={handleRegister} className="bg-white p-8 rounded shadow-md w-full max-w-md space-y-4">
        <h2 className="text-2xl font-bold mb-4">Регистрация</h2>
        {error && <div className="text-red-600 text-sm">{error}</div>}
        
        <input
          type="email"
          className="w-full px-3 py-2 border rounded"
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          required
        />
        
        <input
          type="text"
          className="w-full px-3 py-2 border rounded"
          placeholder="Отображаемое имя"
          value={displayName}
          onChange={e => setDisplayName(e.target.value)}
          required
        />
        
        <input
          type="text"
          className="w-full px-3 py-2 border rounded"
          placeholder="Имя пользователя (необязательно)"
          value={username}
          onChange={e => setUsername(e.target.value)}
        />
        
        <input
          type="password"
          className="w-full px-3 py-2 border rounded"
          placeholder="Пароль"
          value={password}
          onChange={e => setPassword(e.target.value)}
          required
        />
        
        <input
          type="password"
          className="w-full px-3 py-2 border rounded"
          placeholder="Подтвердите пароль"
          value={confirmPassword}
          onChange={e => setConfirmPassword(e.target.value)}
          required
        />
        
        <button
          type="submit"
          className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
          disabled={loading}
        >
          {loading ? 'Регистрация...' : 'Зарегистрироваться'}
        </button>
        
        <div className="text-center mt-4">
          <Link to="/login" className="text-blue-600 hover:underline">
            Уже есть аккаунт? Войти
          </Link>
        </div>
      </form>
    </div>
  );
};

export default RegisterPage;