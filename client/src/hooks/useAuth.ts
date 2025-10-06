import { useState, useEffect } from 'react';

export function useAuth() {
  const [accessToken, setAccessToken] = useState(() => localStorage.getItem('access_token'));
  const [refreshToken, setRefreshToken] = useState(() => localStorage.getItem('refresh_token'));
  const [isValidating, setIsValidating] = useState(true);

  // Проверка валидности токена при загрузке
  useEffect(() => {
    if (accessToken) {
      validateToken();
    } else {
      setIsValidating(false);
    }
  }, [accessToken]);

  const validateToken = async () => {
    try {
      const response = await fetch('/api/profile', {
        headers: {
          'Authorization': `Bearer ${accessToken}`,
        },
      });
      
      if (!response.ok) {
        // Access token невалиден, пробуем обновить
        if (response.status === 401 && refreshToken) {
          await refreshAccessToken();
        } else {
          logout();
        }
      }
    } catch (error) {
      // Ошибка сети, очищаем токены
      logout();
    } finally {
      setIsValidating(false);
    }
  };

  const refreshAccessToken = async () => {
    try {
      const response = await fetch('/api/auth/refresh-token', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (response.ok) {
        const data = await response.json();
        if (data.success) {
          setAccessToken(data.data.access_token);
          setRefreshToken(data.data.refresh_token);
          localStorage.setItem('access_token', data.data.access_token);
          localStorage.setItem('refresh_token', data.data.refresh_token);
          return;
        }
      }
      
      // Если refresh не удался, выходим
      logout();
    } catch (error) {
      logout();
    }
  };

  const login = (accessToken: string, refreshToken: string) => {
    localStorage.setItem('access_token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
    setAccessToken(accessToken);
    setRefreshToken(refreshToken);
  };

  const logout = async () => {
    // Отзываем refresh token на сервере
    if (refreshToken) {
      try {
        await fetch('/api/auth/logout', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ refresh_token: refreshToken }),
        });
      } catch (error) {
        // Игнорируем ошибки при logout
      }
    }

    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    setAccessToken(null);
    setRefreshToken(null);
  };

  const isAuthenticated = !!accessToken && !isValidating;

  return { 
    accessToken, 
    refreshToken, 
    login, 
    logout, 
    isAuthenticated, 
    isValidating,
    refreshAccessToken 
  };
}