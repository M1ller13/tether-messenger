import { useState, KeyboardEvent } from 'react';
import { chatAPI } from '../api/client';

interface MessageInputProps {
  chatId: string | null;
  onMessageSent: () => void;
}

const MessageInput = ({ chatId, onMessageSent }: MessageInputProps) => {
  const [message, setMessage] = useState('');
  const [sending, setSending] = useState(false);

  const handleSend = async () => {
    if (!message.trim() || !chatId || sending) return;

    try {
      setSending(true);
      const response = await chatAPI.sendMessage(chatId, message.trim());
      if (response.success) {
        setMessage('');
        onMessageSent();
      } else {
        console.error('Failed to send message:', response.error);
      }
    } catch (err) {
      console.error('Error sending message:', err);
    } finally {
      setSending(false);
    }
  };

  const handleKeyPress = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  if (!chatId) {
    return (
      <div className="border-t border-gray-200 p-4 bg-white">
        <div className="text-center text-gray-500">
          Select a chat to send messages
        </div>
      </div>
    );
  }

  return (
    <div className="border-t border-gray-200 p-4 bg-white">
      <div className="flex space-x-4">
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="Type a message..."
          disabled={sending}
          className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:opacity-50"
        />
        <button
          onClick={handleSend}
          disabled={!message.trim() || sending}
          className="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-500 disabled:cursor-not-allowed"
        >
          {sending ? 'Sending...' : 'Send'}
        </button>
      </div>
    </div>
  );
};

export default MessageInput;