import React, { useState } from 'react';
import { Send } from 'lucide-react';

interface CommandInputProps {
  onSend: (command: string, type?: string) => void;
}

const CommandInput: React.FC<CommandInputProps> = ({ onSend }) => {
  const [input, setInput] = useState('');
  const [type, setType] = useState('shell');

  const handleSend = () => {
    if (!input.trim()) return;
    onSend(input.trim(), type);
    setInput('');
  };

  return (
    <div className="flex items-center space-x-2 mt-3">
      <select
        value={type}
        onChange={e => setType(e.target.value)}
        className="px-2 py-2 text-sm rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200 focus:ring-2 focus:ring-violet-500 focus:outline-none"
      >
        <option value="shell">Shell</option>
        <option value="intern">Internal</option>
      </select>
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