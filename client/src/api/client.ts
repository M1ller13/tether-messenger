import axios from 'axios';
import { AuthResponse, ApiResponse, Chat, Message, User } from '../types';

const API_BASE_URL = '/api';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  register: async (username: string, password: string): Promise<AuthResponse> => {
    const response = await apiClient.post('/auth/register', { username, password });
    return response.data;
  },

  login: async (username: string, password: string): Promise<AuthResponse> => {
    const response = await apiClient.post('/auth/login', { username, password });
    return response.data;
  },
};

export const chatAPI = {
  async getChats() {
    const token = localStorage.getItem('access_token');
    const res = await fetch('/api/chats', {
      headers: { 'Authorization': `Bearer ${token}` },
    });
    return res.json();
  },
  async createChat(otherUserId: string) {
    const token = localStorage.getItem('access_token');
    const res = await fetch('/api/chats', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ other_user_id: otherUserId }),
    });
    return res.json();
  },
  async getMessages(chatId: string) {
    const token = localStorage.getItem('access_token');
    const res = await fetch(`/api/chats/${chatId}/messages`, {
      headers: { 'Authorization': `Bearer ${token}` },
    });
    return res.json();
  },

  async getChat(chatId: string) {
    const token = localStorage.getItem('access_token');
    const res = await fetch(`/api/chats/${chatId}`, {
      headers: { 'Authorization': `Bearer ${token}` },
    });
    return res.json();
  },
  async sendMessage(chatId: string, content: string) {
    const token = localStorage.getItem('access_token');
    const res = await fetch('/api/messages', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ chat_id: chatId, content }),
    });
    return res.json();
  },

  async sendEncryptedMessage(payload: {
    chat_id: string;
    ciphertext: string;
    nonce: string;
    alg: string;
    ephemeral_pub: string;
  }) {
    const token = localStorage.getItem('access_token');
    const res = await fetch('/api/messages', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });
    return res.json();
  },

  async searchUsers(query: string) {
    const token = localStorage.getItem('access_token');
    const res = await fetch(`/api/users/search?query=${encodeURIComponent(query)}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return res.json();
  },

  async getProfile() {
    const token = localStorage.getItem('access_token');
    const res = await fetch('/api/profile', {
      headers: { 'Authorization': `Bearer ${token}` },
    });
    return res.json();
  },

  async updateProfile(profile: Partial<User>) {
    const token = localStorage.getItem('access_token');
    const res = await fetch('/api/profile', {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(profile),
    });
    return res.json();
  },

  async uploadAvatar(file: File) {
    const token = localStorage.getItem('access_token');
    const formData = new FormData();
    formData.append('avatar', file);
    const res = await fetch('/api/profile/avatar', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` },
      body: formData,
    });
    return res.json();
  },
};

export const e2eeAPI = {
  async publishDeviceKeys(payload: {
    device_id: string;
    identity_key_public: string;
    signed_prekey_public: string;
    signed_prekey_signature: string;
    one_time_prekeys?: { key_id: number; public_key: string }[];
  }) {
    const token = localStorage.getItem('access_token');
    const res = await fetch('/api/e2ee/device-keys', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });
    return res.json();
  },

  async fetchPreKeyBundle(userId: string, deviceId?: string) {
    const token = localStorage.getItem('access_token');
    const url = deviceId
      ? `/api/e2ee/prekey-bundle/${userId}?device_id=${encodeURIComponent(deviceId)}`
      : `/api/e2ee/prekey-bundle/${userId}`;
    const res = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return res.json();
  },
};

export default apiClient;