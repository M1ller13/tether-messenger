import { useState, KeyboardEvent } from 'react';
import { chatAPI, e2eeAPI } from '../api/client';
import { encryptForBundle } from '../crypto/e2ee';

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
      // Try to encrypt using recipient prekey bundle
      const currentUserId = localStorage.getItem('user_id');
      let response: any;
      if (currentUserId) {
        try {
          // Get chat info to determine peer user
          const chatRes = await chatAPI.getChat(chatId);
          if (chatRes?.success && chatRes?.data) {
            const chat = chatRes.data;
            const peerUserId = chat.user1.id === currentUserId ? chat.user2.id : chat.user1.id;
            
            const prekeyRes = await e2eeAPI.fetchPreKeyBundle(peerUserId);
            if (prekeyRes?.success && prekeyRes?.data) {
              const enc = await encryptForBundle(prekeyRes.data, message.trim());
              response = await chatAPI.sendEncryptedMessage({
                chat_id: chatId,
                ciphertext: enc.ciphertext,
                nonce: enc.nonce,
                alg: enc.alg,
                ephemeral_pub: enc.ephemeral_pub,
              });
            } else {
              response = await chatAPI.sendMessage(chatId, message.trim());
            }
          } else {
            response = await chatAPI.sendMessage(chatId, message.trim());
          }
        } catch {
          response = await chatAPI.sendMessage(chatId, message.trim());
        }
      } else {
        response = await chatAPI.sendMessage(chatId, message.trim());
      }
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