import { useState, useEffect } from 'react';

export function useAuth() {
  const [token, setToken] = useState(() => localStorage.getItem('token'));
  const [isValidating, setIsValidating] = useState(true);

  // Проверка валидности токена при загрузке
  useEffect(() => {
    if (token) {
      validateToken();
    } else {
      setIsValidating(false);
    }
  }, [token]);

  const validateToken = async () => {
    try {
      const response = await fetch('/api/profile', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      
      if (!response.ok) {
        // Токен невалиден, очищаем его
        logout();
      }
    } catch (error) {
      // Ошибка сети, очищаем токен
      logout();
    } finally {
      setIsValidating(false);
    }
  };

  const login = (jwt: string) => {
    localStorage.setItem('token', jwt);
    setToken(jwt);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
  };

  const isAuthenticated = !!token && !isValidating;

  // Можно добавить декодирование JWT для получения user info
  const getUser = () => {
    // Пример: return jwt_decode(token) или хранить user отдельно
    return null;
  };

  return { token, login, logout, isAuthenticated, getUser, isValidating };
}