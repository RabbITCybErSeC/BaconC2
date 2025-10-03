import React from 'react';
import { CheckCircle2, Loader2, XCircle } from 'lucide-react';

export interface CommandEntry {
  id: string;
  command: string;
  status: string;
  createdAt: string;
  result?: string;
}

interface CommandTimelineProps {
  commands: CommandEntry[];
}

const CommandTimeline: React.FC<CommandTimelineProps> = ({ commands }) => {
  if (commands.length === 0) {
    return (
      <p className="text-sm text-gray-500 dark:text-gray-400">
        No commands executed yet.
      </p>
    );
  }

  return (
    <div className="space-y-4 max-h-72 overflow-y-auto pr-2">
      {commands.map(cmd => (
        <div
          key={cmd.id}
          className="border-l-2 border-violet-500 pl-4 relative"
        >
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium text-gray-800 dark:text-gray-200">
              {cmd.command}
            </p>
            <span className="text-xs text-gray-500 dark:text-gray-400">
              {new Date(cmd.createdAt).toLocaleTimeString()}
            </span>
          </div>
          <div className="flex items-center space-x-2 mt-1">
            {cmd.status === 'completed' && (
              <CheckCircle2 className="w-4 h-4 text-green-500" />
            )}
            {cmd.status === 'failed' && (
              <XCircle className="w-4 h-4 text-red-500" />
            )}
            {cmd.status !== 'completed' &&
              cmd.status !== 'failed' && (
                <Loader2 className="w-4 h-4 text-blue-500 animate-spin" />
              )}
            <span className="text-xs text-gray-600 dark:text-gray-400">
              {cmd.status}
            </span>
          </div>
          {cmd.result && (
            <pre className="bg-gray-900 text-green-400 font-mono text-xs mt-2 p-2 rounded-lg whitespace-pre-wrap border border-gray-700 dark:border-gray-600">
              {cmd.result}
            </pre>
          )}
        </div>
      ))}
    </div>
  );
};

export default CommandTimeline;