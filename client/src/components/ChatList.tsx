import { useEffect, useState } from 'react';
import { Chat, User } from '../types';
import { chatAPI } from '../api/client';
import LoadingSpinner from './LoadingSpinner';

interface ChatListProps {
  currentUser: User;
  selectedChat: Chat | null;
  onChatSelect: (chat: Chat) => void;
}

const ChatList = ({ currentUser, selectedChat, onChatSelect }: ChatListProps) => {
  const [chats, setChats] = useState<Chat[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadChats();
  }, []);

  const loadChats = async () => {
    try {
      setLoading(true);
      const response = await chatAPI.getChats();
      if (response.success) {
        setChats(response.data);
      } else {
        setError(response.error || 'Failed to load chats');
      }
    } catch (err) {
      setError('Failed to load chats');
    } finally {
      setLoading(false);
    }
  };

  const getOtherUser = (chat: Chat) => {
    return chat.user1.id === currentUser.id ? chat.user2 : chat.user1;
  };

  const formatLastSeen = (lastSeen: string) => {
    const date = new Date(lastSeen);
    const now = new Date();
    const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60));
    
    if (diffInMinutes < 1) return 'now';
    if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
    if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}h ago`;
    return date.toLocaleDateString();
  };

  if (loading) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 mb-2">{error}</p>
          <button 
            onClick={loadChats}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-y-auto">
      {chats.length === 0 ? (
        <div className="p-4 text-center text-gray-500">
          No chats yet. Start a conversation!
        </div>
      ) : (
        <div className="divide-y divide-gray-200">
          {chats.map((chat) => {
            const otherUser = getOtherUser(chat);
            return (
              <div
                key={chat.id}
                onClick={() => onChatSelect(chat)}
                className={`p-4 cursor-pointer hover:bg-gray-50 transition-colors ${
                  selectedChat?.id === chat.id ? 'bg-blue-50 border-r-2 border-blue-500' : ''
                }`}
              >
                <div className="flex items-center space-x-3">
                  <div className="flex-shrink-0">
                    {otherUser.avatar_url ? (
                      <img 
                        src={otherUser.avatar_url} 
                        alt={otherUser.display_name}
                        className="w-10 h-10 rounded-full object-cover"
                      />
                    ) : (
                      <div className="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center">
                        <span className="text-white font-medium">
                          {otherUser.display_name.charAt(0).toUpperCase()}
                        </span>
                      </div>
                    )}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex justify-between items-start">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {otherUser.display_name}
                      </p>
                      {otherUser.last_seen && (
                        <span className="text-xs text-gray-500">
                          {formatLastSeen(otherUser.last_seen)}
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-gray-500 truncate">
                      @{otherUser.username}
                    </p>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default ChatList;