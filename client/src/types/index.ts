export interface User {
  id: string;
  phone: string;
  username: string;
  display_name: string;
  bio?: string;
  avatar_url?: string;
  last_seen?: string;
  created_at: string;
}

export interface Message {
  id: string;
  content: string;
  sender: User;
  chat_id: string;
  created_at: string;
}

export interface Chat {
  id: string;
  user1: User;
  user2: User;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  success: boolean;
  data: {
    token: string;
    user: User;
  };
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: string;
}