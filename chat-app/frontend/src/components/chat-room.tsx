import React, { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/router';
import { io, Socket } from 'socket.io-client';
import ChatMessage from './chat-message';
import ChatInput from './chat-input';
import UsernamePrompt from './username-prompt';

interface ChatRoomProps {
  room: string;
}

const ChatRoom: React.FC<ChatRoomProps> = ({ room }) => {
  const router = useRouter();
  const [username, setUsername] = useState<string | null>(null);
  const [messages, setMessages] = useState<{ id: string; user: string; content: string }[]>([]);
  const [socket, setSocket] = useState<Socket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const storedUsername = localStorage.getItem('username');
    if (storedUsername) {
      setUsername(storedUsername);
    }
  }, []);

  useEffect(() => {
    if (!username) return;

    const newSocket = io('http://localhost:8080', {
      query: { room, username },
    });

    setSocket(newSocket);

    newSocket.on('connect', () => {
      console.log('Connected to server');
    });

    newSocket.on('message', (data) => {
      setMessages((prevMessages) => [...prevMessages, data]);
    });

    // Fetch existing messages
    fetch(`/api/messages?room=${room}`)
      .then((res) => res.json())
      .then((data) => {
        setMessages(data);
      })
      .catch((err) => {
        console.error('Error fetching messages:', err);
      });

    return () => {
      newSocket.disconnect();
    };
  }, [username, room]);

  const sendMessage = (content: string) => {
    if (socket && username) {
      socket.emit('sendMessage', { user: username, content });
    }
  };

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  if (!username) {
    return <UsernamePrompt onEnter={(name) => setUsername(name)} />;
  }

  return (
    <div className="flex flex-col h-screen">
      <header className="bg-blue-500 text-white p-4 text-center">
        <h1 className="text-2xl font-bold">{room.charAt(0).toUpperCase() + room.slice(1)} Room</h1>
      </header>
      <main className="flex-1 overflow-y-auto p-4 space-y-2">
        {messages && messages.map((msg, index) => (
          <ChatMessage key={index} user={msg.user} content={msg.content} isUser={msg.user === username} />
        ))}
        <div ref={messagesEndRef} />
      </main>
      <footer className="p-4">
        <ChatInput onSendMessage={sendMessage} />
      </footer>
    </div>
  );
};

export default ChatRoom;