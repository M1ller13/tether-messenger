import { useEffect, useState, useRef } from 'react';
import { Message, User } from '../types';
import { tryDecryptMessage } from '../crypto/e2ee';
import { chatAPI } from '../api/client';
import LoadingSpinner from './LoadingSpinner';

interface MessageListProps {
  chatId: string | null;
  currentUser: User;
}

const MessageList = ({ chatId, currentUser }: MessageListProps) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (chatId) {
      loadMessages();
    } else {
      setMessages([]);
    }
  }, [chatId]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const loadMessages = async () => {
    if (!chatId) return;
    
    try {
      setLoading(true);
      setError(null);
      const response = await chatAPI.getMessages(chatId);
      if (response.success) {
        const decrypted = await Promise.all((response.data as Message[]).map(async (m: any) => {
          const text = await tryDecryptMessage({
            ciphertext: (m as any).ciphertext,
            nonce: (m as any).nonce,
            alg: (m as any).alg,
            ephemeral_pub: (m as any).ephemeral_pub,
            content: m.content,
          });
          return { ...m, content: text ?? m.content } as Message;
        }));
        setMessages(decrypted);
      } else {
        setError(response.error || 'Failed to load messages');
      }
    } catch (err) {
      setError('Failed to load messages');
    } finally {
      setLoading(false);
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const formatTime = (dateString: string) => {
    return new Date(dateString).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  if (!chatId) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center text-gray-500">
          <p>Select a chat to start messaging</p>
        </div>
      </div>
    );
  }

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
            onClick={loadMessages}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-y-auto p-4 space-y-4">
      {messages.length === 0 ? (
        <div className="text-center text-gray-500 mt-8">
          No messages yet. Start the conversation!
        </div>
      ) : (
        <>
          {messages.map((message) => {
            const isOwnMessage = message.sender.id === currentUser.id;
            
            return (
              <div
                key={message.id}
                className={`flex ${isOwnMessage ? 'justify-end' : 'justify-start'}`}
              >
                <div
                  className={`max-w-xs lg:max-w-md px-4 py-2 rounded-lg ${
                    isOwnMessage
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-200 text-gray-900'
                  }`}
                >
                  <div className="flex items-end space-x-2">
                    <div className="flex-1">
                      <div className="flex items-center space-x-1">
                        <p className="text-sm">{message.content}</p>
                        {message.ciphertext && message.alg && (
                          <span 
                            className="text-xs opacity-75"
                            title="End-to-end encrypted"
                          >
                            🔒
                          </span>
                        )}
                      </div>
                    </div>
                    <div className="flex-shrink-0">
                      <span
                        className={`text-xs ${
                          isOwnMessage ? 'text-blue-100' : 'text-gray-500'
                        }`}
                      >
                        {formatTime(message.created_at)}
                      </span>
                    </div>
                  </div>
                  {!isOwnMessage && (
                    <p className="text-xs text-gray-600 mt-1">
                      {message.sender.display_name}
                    </p>
                  )}
                </div>
              </div>
            );
          })}
          <div ref={messagesEndRef} />
        </>
      )}
    </div>
  );
};

export default MessageList;