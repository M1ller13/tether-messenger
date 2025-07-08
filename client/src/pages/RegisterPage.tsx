import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

const RegisterPage = () => {
  const [step, setStep] = useState<'register' | 'code'>('register');
  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [code, setCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ display_name: name, phone }),
      });
      const data = await res.json();
      if (data.success) {
        setStep('code');
      } else {
        setError(data.error || 'Ошибка регистрации');
      }
    } catch {
      setError('Ошибка сети');
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyCode = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await fetch('/api/auth/verify-code', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone, code }),
      });
      const data = await res.json();
      if (data.success && data.data && data.data.token) {
        localStorage.setItem('token', data.data.token);
        navigate('/chats');
      } else {
        setError(data.error || 'Ошибка подтверждения кода');
      }
    } catch {
      setError('Ошибка сети');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <form onSubmit={step === 'register' ? handleRegister : handleVerifyCode} className="bg-white p-8 rounded shadow-md w-full max-w-md space-y-4">
        <h2 className="text-2xl font-bold mb-4">Регистрация</h2>
        {error && <div className="text-red-600 text-sm">{error}</div>}
        {step === 'register' ? (
          <>
            <input
              type="text"
              className="w-full px-3 py-2 border rounded"
              placeholder="Ваше имя"
              value={name}
              onChange={e => setName(e.target.value)}
              required
            />
            <input
              type="tel"
              className="w-full px-3 py-2 border rounded"
              placeholder="Телефон (+7...)"
              value={phone}
              onChange={e => setPhone(e.target.value)}
              required
            />
            <button
              type="submit"
              className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
              disabled={loading}
            >
              {loading ? 'Отправка...' : 'Зарегистрироваться'}
            </button>
            <div className="text-center mt-4">
              <Link to="/login" className="text-blue-600 hover:underline">
                Уже есть аккаунт? Войти
              </Link>
            </div>
          </>
        ) : (
          <>
            <input
              type="text"
              className="w-full px-3 py-2 border rounded"
              placeholder="Код из SMS"
              value={code}
              onChange={e => setCode(e.target.value)}
              required
              maxLength={6}
            />
            <button
              type="submit"
              className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
              disabled={loading}
            >
              {loading ? 'Проверка...' : 'Подтвердить'}
            </button>
            <button
              type="button"
              className="w-full text-blue-600 mt-2 hover:underline"
              onClick={() => setStep('register')}
            >
              Изменить данные
            </button>
          </>
        )}
      </form>
    </div>
  );
};

export default RegisterPage;