#!/usr/bin/env python3
"""
Тестовый скрипт для проверки API endpoints Go сервера.
"""

import requests
import json
import os
import random

# Отключаем прокси для локальных запросов
os.environ['NO_PROXY'] = 'localhost,127.0.0.1'

BASE_URL = "http://localhost:8081"

# Генерируем уникальное имя пользователя для каждого теста
username = f"testuser_{random.randint(1000, 9999)}"
password = "testpassword"

def test_health():
    """Тест health check endpoint"""
    print("🔍 Тестирую health check...")
    try:
        response = requests.get(f"{BASE_URL}/health", proxies={'http': None, 'https': None})
        print(f"✅ Health check: {response.status_code} - {response.json()}")
        return True
    except Exception as e:
        print(f"❌ Health check failed: {e}")
        return False

def test_register():
    """Тест регистрации"""
    print("\n🔍 Тестирую регистрацию...")
    try:
        data = {
            "username": username,
            "password": password
        }
        response = requests.post(f"{BASE_URL}/api/auth/register", json=data, proxies={'http': None, 'https': None})
        print(f"✅ Register: {response.status_code}")
        if response.status_code == 200:
            result = response.json()
            print(f"   Token: {result['data']['token'][:20]}...")
            print(f"   User: {result['data']['user']['username']}")
            return result['data']['token']
        else:
            print(f"   Error: {response.json()}")
            return None
    except Exception as e:
        print(f"❌ Register failed: {e}")
        return None

def test_login():
    """Тест входа"""
    print("\n🔍 Тестирую вход...")
    try:
        data = {
            "username": username,
            "password": password
        }
        response = requests.post(f"{BASE_URL}/api/auth/login", json=data, proxies={'http': None, 'https': None})
        print(f"✅ Login: {response.status_code}")
        if response.status_code == 200:
            result = response.json()
            print(f"   Token: {result['data']['token'][:20]}...")
            print(f"   User: {result['data']['user']['username']}")
            return result['data']['token']
        else:
            print(f"   Error: {response.json()}")
            return None
    except Exception as e:
        print(f"❌ Login failed: {e}")
        return None

def test_chats(token):
    """Тест получения чатов"""
    print("\n🔍 Тестирую получение чатов...")
    try:
        headers = {"Authorization": f"Bearer {token}"}
        response = requests.get(f"{BASE_URL}/api/chats", headers=headers, proxies={'http': None, 'https': None})
        print(f"✅ Get chats: {response.status_code}")
        if response.status_code == 200:
            result = response.json()
            print(f"   Chats count: {len(result['data'])}")
            for chat in result['data']:
                print(f"   Chat ID: {chat['id']}, Users: {chat['user1']['username']} <-> {chat['user2']['username']}")
        else:
            print(f"   Error: {response.json()}")
    except Exception as e:
        print(f"❌ Get chats failed: {e}")

def test_messages(token):
    """Тест получения сообщений"""
    print("\n🔍 Тестирую получение сообщений...")
    try:
        headers = {"Authorization": f"Bearer {token}"}
        response = requests.get(f"{BASE_URL}/api/chats/1/messages", headers=headers, proxies={'http': None, 'https': None})
        print(f"✅ Get messages: {response.status_code}")
        if response.status_code == 200:
            result = response.json()
            print(f"   Messages count: {len(result['data'])}")
            for msg in result['data']:
                print(f"   Message: {msg['sender']['username']}: {msg['content']}")
        else:
            print(f"   Error: {response.json()}")
    except Exception as e:
        print(f"❌ Get messages failed: {e}")

def test_create_chat(token):
    """Тест создания чата"""
    print("\n🔍 Тестирую создание чата...")
    try:
        headers = {"Authorization": f"Bearer {token}"}
        data = {"other_user_id": 2}
        response = requests.post(f"{BASE_URL}/api/chats", json=data, headers=headers, proxies={'http': None, 'https': None})
        print(f"✅ Create chat: {response.status_code}")
        if response.status_code == 200:
            result = response.json()
            print(f"   Chat created: ID {result['data']['id']}")
        else:
            print(f"   Error: {response.json()}")
    except Exception as e:
        print(f"❌ Create chat failed: {e}")

def main():
    print("🚀 Тестирование API Tether Messenger")
    print("=" * 50)
    
    # Тест health check
    if not test_health():
        print("❌ Сервер недоступен!")
        return
    
    # Тест регистрации
    token = test_register()
    if not token:
        print("❌ Регистрация не удалась!")
        return
    
    # Тест входа
    login_token = test_login()
    if not login_token:
        print("❌ Вход не удался!")
        return
    
    # Тест защищенных endpoints
    test_chats(token)
    test_messages(token)
    test_create_chat(token)
    
    print("\n✅ Все тесты завершены!")

if __name__ == "__main__":
    main()