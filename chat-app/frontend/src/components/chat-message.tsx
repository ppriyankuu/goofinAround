import React from 'react';

interface ChatMessageProps {
  user: string;
  content: string;
  isUser?: boolean; 
}

const ChatMessage: React.FC<ChatMessageProps> = ({ user, content, isUser = false }) => {
  return (
    <div className={`flex ${isUser ? 'justify-end' : 'justify-start'}`}>
      <div
        className={`max-w-md p-2 rounded-lg mb-2 ${
          isUser ? 'bg-blue-500 text-white' : 'bg-green-500 text-white'
        }`}
      >
        <span className="font-bold">{user}:</span> {content}
      </div>
    </div>
  );
};

export default ChatMessage;
