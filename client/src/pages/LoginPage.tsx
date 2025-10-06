import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const LoginPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      const data = await res.json();
      if (data.success && data.data) {
        login(data.data.access_token, data.data.refresh_token);
        navigate('/chats');
      } else {
        setError(data.error || 'Ошибка входа');
      }
    } catch {
      setError('Ошибка сети');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <form onSubmit={handleLogin} className="bg-white p-8 rounded shadow-md w-full max-w-md space-y-4">
        <h2 className="text-2xl font-bold mb-4">Вход в Tether Messenger</h2>
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
          type="password"
          className="w-full px-3 py-2 border rounded"
          placeholder="Пароль"
          value={password}
          onChange={e => setPassword(e.target.value)}
          required
        />
        
        <button
          type="submit"
          className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
          disabled={loading}
        >
          {loading ? 'Вход...' : 'Войти'}
        </button>
        
        <div className="text-center mt-4 space-y-2">
          <Link to="/register" className="text-blue-600 hover:underline block">
            Нет аккаунта? Зарегистрироваться
          </Link>
          <Link to="/forgot-password" className="text-blue-600 hover:underline block">
            Забыли пароль?
          </Link>
        </div>
      </form>
    </div>
  );
};

export default LoginPage;