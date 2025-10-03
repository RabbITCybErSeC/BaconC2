import React, { useState } from 'react';
import { Send } from 'lucide-react';

interface CommandInputProps {
  onSend: (command: string) => void;
}

const CommandInput: React.FC<CommandInputProps> = ({ onSend }) => {
  const [input, setInput] = useState('');

  const handleSend = () => {
    if (!input.trim()) return;
    onSend(input.trim());
    setInput('');
  };

  return (
    <div className="flex items-center space-x-2 mt-3">
      <input
        type="text"
        value={input}
        onChange={e => setInput(e.target.value)}
        placeholder="Enter command..."
        className="flex-1 px-3 py-2 text-sm rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200 focus:ring-2 focus:ring-violet-500 focus:outline-none"
      />
      <button
        onClick={handleSend}
        className="flex items-center bg-violet-600 text-white px-3 py-2 rounded-lg hover:bg-violet-700 focus:ring-2 focus:ring-violet-400"
      >
        <Send className="w-4 h-4 mr-1" />
        Send
      </button>
    </div>
  );
};

export default CommandInput;