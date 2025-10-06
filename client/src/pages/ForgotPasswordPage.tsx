import { useState } from 'react';
import { Link } from 'react-router-dom';

const ForgotPasswordPage = () => {
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleRequestReset = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const res = await fetch('/api/auth/request-password-reset', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      });
      const data = await res.json();
      if (data.success) {
        setSuccess(true);
      } else {
        setError(data.error || 'Ошибка отправки запроса');
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
          <h2 className="text-2xl font-bold mb-4 text-green-600">Письмо отправлено!</h2>
          <p className="text-gray-600 mb-4">
            Если аккаунт с таким email существует, мы отправили ссылку для сброса пароля.
            Проверьте почту и перейдите по ссылке.
          </p>
          <Link to="/login" className="text-blue-600 hover:underline">
            Вернуться к входу
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <form onSubmit={handleRequestReset} className="bg-white p-8 rounded shadow-md w-full max-w-md space-y-4">
        <h2 className="text-2xl font-bold mb-4">Сброс пароля</h2>
        {error && <div className="text-red-600 text-sm">{error}</div>}
        
        <p className="text-gray-600 text-sm">
          Введите ваш email, и мы отправим ссылку для сброса пароля.
        </p>
        
        <input
          type="email"
          className="w-full px-3 py-2 border rounded"
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          required
        />
        
        <button
          type="submit"
          className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
          disabled={loading}
        >
          {loading ? 'Отправка...' : 'Отправить ссылку'}
        </button>
        
        <div className="text-center mt-4">
          <Link to="/login" className="text-blue-600 hover:underline">
            Вернуться к входу
          </Link>
        </div>
      </form>
    </div>
  );
};

export default ForgotPasswordPage;
