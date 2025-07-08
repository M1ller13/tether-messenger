import { useState, useEffect, useRef } from 'react';
import { useAuth } from '../hooks/useAuth';
import ChatList from '../components/ChatList';
import MessageList from '../components/MessageList';
import MessageInput from '../components/MessageInput';
import { Chat, User } from '../types';
import { chatAPI } from '../api/client';

const ChatPage = () => {
  const { logout } = useAuth();
  const [user, setUser] = useState<User | null>(null);
  const [selectedChat, setSelectedChat] = useState<Chat | null>(null);
  const [search, setSearch] = useState('');
  const [searchResults, setSearchResults] = useState<User[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState('');
  const searchBoxRef = useRef<HTMLDivElement>(null);
  const [showResults, setShowResults] = useState(false);
  const [profileOpen, setProfileOpen] = useState(false);
  const [profile, setProfile] = useState<User | null>(null);
  const [profileEdit, setProfileEdit] = useState<Partial<User>>({});
  const [profileError, setProfileError] = useState('');
  const [profileLoading, setProfileLoading] = useState(false);
  const [avatarUploading, setAvatarUploading] = useState(false);

  // Загрузка профиля пользователя при монтировании
  useEffect(() => {
    loadUserProfile();
  }, []);

  const loadUserProfile = async () => {
    try {
      const res = await chatAPI.getProfile();
      if (res.success) {
        setUser(res.data);
      } else {
        console.error('Failed to load user profile:', res.error);
      }
    } catch (err) {
      console.error('Error loading user profile:', err);
    }
  };

  // Поиск пользователей
  const handleSearch = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearch(value);
    setSearchError('');
    if (value.length < 2) {
      setSearchResults([]);
      setShowResults(false);
      return;
    }
    setSearchLoading(true);
    setShowResults(true);
    try {
      const res = await chatAPI.searchUsers(value);
      if (res.success) {
        setSearchResults(res.data);
        if (res.data.length === 0) setSearchError('Пользователь не найден');
      } else {
        setSearchResults([]);
        setSearchError('Ошибка поиска');
      }
    } catch {
      setSearchResults([]);
      setSearchError('Ошибка поиска');
    } finally {
      setSearchLoading(false);
    }
  };

  // Создание чата с выбранным пользователем
  const handleUserSelect = async (userToChat: User) => {
    try {
      const res = await chatAPI.createChat(userToChat.id);
      if (res.success && res.data) {
        setSelectedChat(res.data);
        setSearch('');
        setSearchResults([]);
        setShowResults(false);
      }
    } catch {
      setSearchError('Ошибка создания чата');
    }
  };

  // Автоматическое закрытие выпадающего списка при клике вне
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (searchBoxRef.current && !searchBoxRef.current.contains(event.target as Node)) {
        setShowResults(false);
      }
    };
    if (showResults) {
      document.addEventListener('mousedown', handleClickOutside);
    } else {
      document.removeEventListener('mousedown', handleClickOutside);
    }
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [showResults]);

  const openProfile = async () => {
    setProfileError('');
    setProfileLoading(true);
    setProfileOpen(true);
    try {
      const res = await chatAPI.getProfile();
      if (res.success) {
        setProfile(res.data);
        setProfileEdit(res.data);
      } else {
        setProfileError(res.error || 'Ошибка загрузки профиля');
      }
    } catch {
      setProfileError('Ошибка сети');
    } finally {
      setProfileLoading(false);
    }
  };

  const saveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setProfileError('');
    setProfileLoading(true);
    try {
      const res = await chatAPI.updateProfile(profileEdit);
      if (res.success) {
        setProfile(res.data);
        setProfileEdit(res.data);
        setProfileOpen(false);
      } else {
        setProfileError(res.error || 'Ошибка сохранения');
      }
    } catch {
      setProfileError('Ошибка сети');
    } finally {
      setProfileLoading(false);
    }
  };

  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files || !e.target.files[0]) return;
    setAvatarUploading(true);
    setProfileError('');
    try {
      const res = await chatAPI.uploadAvatar(e.target.files[0]);
      if (res.success) {
        setProfileEdit(p => ({ ...p, avatar_url: res.avatar_url }));
        setProfile(p => p ? { ...p, avatar_url: res.avatar_url } : p);
      } else {
        setProfileError(res.error || 'Ошибка загрузки аватара');
      }
    } catch {
      setProfileError('Ошибка сети');
    } finally {
      setAvatarUploading(false);
    }
  };

  const handleMessageSent = () => {
    // Сообщение отправлено, можно обновить UI если нужно
  };

  return (
    <div className="h-screen flex bg-gray-100">
      {!user ? (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Загрузка...</p>
          </div>
        </div>
      ) : (
        <>
          {/* Sidebar */}
          <div className="w-80 bg-white border-r border-gray-200 flex flex-col">
            {/* Header */}
            <div className="p-4 border-b border-gray-200">
              <div className="flex items-center justify-between">
                <h1 className="text-xl font-semibold text-gray-900">Tether</h1>
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-gray-600 cursor-pointer" onClick={openProfile} title="Профиль">
                    {user?.display_name || user?.username}
                  </span>
                  <button
                    onClick={logout}
                    className="text-sm text-red-600 hover:text-red-800"
                  >
                    Logout
                  </button>
                </div>
              </div>
            </div>

            {/* Search bar */}
            <div className="p-2 border-b border-gray-200" ref={searchBoxRef}>
              <input
                type="text"
                className="w-full px-3 py-2 border rounded focus:outline-none"
                placeholder="Поиск пользователей..."
                value={search}
                onChange={handleSearch}
                onFocus={() => search.length >= 2 && setShowResults(true)}
              />
              {searchLoading && <div className="text-xs text-gray-400 mt-1">Поиск...</div>}
              {showResults && (
                <div className="bg-white border rounded shadow mt-1 max-h-40 overflow-y-auto z-10">
                  {searchResults.length > 0 ? (
                    searchResults.map((u) => (
                      <div
                        key={u.id}
                        className="px-3 py-2 hover:bg-blue-50 cursor-pointer transition-colors"
                        onClick={() => handleUserSelect(u)}
                      >
                        <div className="font-medium text-gray-900">{u.display_name}</div>
                        <div className="text-xs text-blue-600">@{u.username}</div>
                        <div className="text-xs text-gray-500">{u.bio}</div>
                      </div>
                    ))
                  ) : (
                    searchError && <div className="px-3 py-2 text-gray-400">{searchError}</div>
                  )}
                </div>
              )}
            </div>

            {/* Chat List */}
            <ChatList
              currentUser={user}
              selectedChat={selectedChat}
              onChatSelect={setSelectedChat}
            />
          </div>

          {/* Main Chat Area */}
          <div className="flex-1 flex flex-col">
            {selectedChat ? (
              <>
                <MessageList chatId={selectedChat.id} currentUser={user} />
                <MessageInput chatId={selectedChat.id} onMessageSent={handleMessageSent} />
              </>
            ) : (
              <div className="flex-1 flex items-center justify-center">
                <div className="text-center">
                  <h2 className="text-xl font-semibold text-gray-900 mb-2">
                    Welcome to Tether
                  </h2>
                  <p className="text-gray-600">Select a chat to start messaging</p>
                </div>
              </div>
            )}
          </div>
        </>
      )}

      {/* Модалка профиля */}
      {profileOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
          <form onSubmit={saveProfile} className="bg-white p-6 rounded shadow-md w-full max-w-md space-y-4 relative">
            <button type="button" onClick={() => setProfileOpen(false)} className="absolute top-2 right-2 text-gray-400 hover:text-gray-700">✕</button>
            <h2 className="text-xl font-bold mb-2">Профиль</h2>
            {profileError && <div className="text-red-600 text-sm">{profileError}</div>}
            
            {/* Аватар */}
            <div className="flex flex-col items-center mb-2">
              {profileEdit.avatar_url ? (
                <img src={profileEdit.avatar_url} alt="avatar" className="w-20 h-20 rounded-full object-cover border mb-2" />
              ) : (
                <div className="w-20 h-20 rounded-full bg-blue-200 flex items-center justify-center text-3xl text-blue-700 font-bold mb-2">
                  {profileEdit.display_name?.charAt(0).toUpperCase() || profileEdit.username?.charAt(0).toUpperCase() || '?'}
                </div>
              )}
              <label className="cursor-pointer text-blue-600 hover:underline text-sm">
                {avatarUploading ? 'Загрузка...' : 'Сменить аватар'}
                <input type="file" accept="image/*" className="hidden" onChange={handleAvatarChange} disabled={avatarUploading} />
              </label>
            </div>
            
            <input
              type="text"
              className="w-full px-3 py-2 border rounded"
              placeholder="Отображаемое имя"
              value={profileEdit.display_name || ''}
              onChange={e => setProfileEdit(p => ({ ...p, display_name: e.target.value }))}
            />
            <input
              type="text"
              className="w-full px-3 py-2 border rounded bg-gray-100"
              placeholder="@username"
              value={profileEdit.username || ''}
              disabled
            />
            <textarea
              className="w-full px-3 py-2 border rounded"
              placeholder="О себе"
              value={profileEdit.bio || ''}
              onChange={e => setProfileEdit(p => ({ ...p, bio: e.target.value }))}
              rows={3}
            />
            
            <button
              type="submit"
              className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
              disabled={profileLoading}
            >
              {profileLoading ? 'Сохранение...' : 'Сохранить'}
            </button>
          </form>
        </div>
      )}
    </div>
  );
};

export default ChatPage;