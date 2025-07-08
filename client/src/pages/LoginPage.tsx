import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

const LoginPage = () => {
  const [step, setStep] = useState<'phone' | 'code'>('phone');
  const [phone, setPhone] = useState('');
  const [code, setCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleRequestCode = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await fetch('/api/auth/request-code', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone }),
      });
      const data = await res.json();
      if (data.success) {
        setStep('code');
      } else {
        setError(data.error || 'Ошибка отправки кода');
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
      <form onSubmit={step === 'phone' ? handleRequestCode : handleVerifyCode} className="bg-white p-8 rounded shadow-md w-full max-w-md space-y-4">
        <h2 className="text-2xl font-bold mb-4">Вход по номеру телефона</h2>
        {error && <div className="text-red-600 text-sm">{error}</div>}
        {step === 'phone' ? (
          <>
            <input
              type="tel"
              className="w-full px-3 py-2 border rounded"
              placeholder="+7XXXXXXXXXX"
              value={phone}
              onChange={e => setPhone(e.target.value)}
              required
            />
            <button
              type="submit"
              className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
              disabled={loading}
            >
              {loading ? 'Отправка...' : 'Получить код'}
            </button>
            <div className="text-center mt-4">
              <Link to="/register" className="text-blue-600 hover:underline">
                Нет аккаунта? Зарегистрироваться
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
              {loading ? 'Проверка...' : 'Войти'}
            </button>
            <button
              type="button"
              className="w-full text-blue-600 mt-2 hover:underline"
              onClick={() => setStep('phone')}
            >
              Изменить номер
            </button>
          </>
        )}
      </form>
    </div>
  );
};

export default LoginPage;