import { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

const VerifyEmailPage = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(true);
  const [token, setToken] = useState('');

  useEffect(() => {
    const tokenParam = searchParams.get('token');
    if (!tokenParam) {
      setError('Недействительная ссылка для верификации');
      setLoading(false);
      return;
    }
    setToken(tokenParam);
    verifyEmail(tokenParam);
  }, [searchParams]);

  const verifyEmail = async (token: string) => {
    try {
      const res = await fetch(`/api/auth/verify-email?token=${token}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });
      const data = await res.json();
      if (data.success) {
        setSuccess(true);
        setTimeout(() => {
          navigate('/login');
        }, 3000);
      } else {
        setError(data.error || 'Ошибка верификации email');
      }
    } catch {
      setError('Ошибка сети');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-100">
        <div className="bg-white p-8 rounded shadow-md w-full max-w-md text-center">
          <h2 className="text-2xl font-bold mb-4">Верификация email...</h2>
          <p className="text-gray-600">Пожалуйста, подождите</p>
        </div>
      </div>
    );
  }

  if (success) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-100">
        <div className="bg-white p-8 rounded shadow-md w-full max-w-md text-center">
          <h2 className="text-2xl font-bold mb-4 text-green-600">Email подтвержден!</h2>
          <p className="text-gray-600 mb-4">
            Ваш email успешно подтвержден. Вы будете перенаправлены на страницу входа.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <div className="bg-white p-8 rounded shadow-md w-full max-w-md text-center">
        <h2 className="text-2xl font-bold mb-4 text-red-600">Ошибка верификации</h2>
        <p className="text-gray-600 mb-4">
          {error}
        </p>
        <button
          onClick={() => navigate('/login')}
          className="text-blue-600 hover:underline"
        >
          Перейти к входу
        </button>
      </div>
    </div>
  );
};

export default VerifyEmailPage;
