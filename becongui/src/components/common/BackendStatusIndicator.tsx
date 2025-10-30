import React, { useState, useEffect } from 'react';
import { Check, X, Loader2 } from 'lucide-react';

type Status = 'checking' | 'ok' | 'error';

const BackendStatusIndicator: React.FC = () => {
  const [status, setStatus] = useState<Status>('checking');

  useEffect(() => {
    const checkStatus = async () => {
      try {
        const response = await fetch('/api/v1/general/health');
        if (response.ok) {
          const data = await response.json();
          setStatus(data.status === 'ok' ? 'ok' : 'error');
        } else {
          setStatus('error');
        }
      } catch (error) {
        console.error('Error checking backend health:', error);
        setStatus('error');
      }
    };

    checkStatus();
    const intervalId = setInterval(checkStatus, 30000);

    return () => clearInterval(intervalId);
  }, []);

  const renderIndicator = () => {
    switch (status) {
      case 'ok':
        return (
          <div
            className="flex items-center justify-center w-6 h-6 rounded-full bg-green-500/10 text-green-500"
            title="Backend Status: Connected"
          >
            <Check size={14} strokeWidth={2.5} />
          </div>
        );
      case 'error':
        return (
          <div
            className="flex items-center justify-center w-6 h-6 rounded-full bg-red-500/10 text-red-500"
            title="Backend Status: Error"
          >
            <X size={14} strokeWidth={2.5} />
          </div>
        );
      case 'checking':
      default:
        return (
          <div
            className="flex items-center justify-center w-6 h-6 rounded-full bg-yellow-500/10 text-yellow-500"
            title="Backend Status: Checking..."
          >
            <Loader2 size={14} className="animate-spin" strokeWidth={2.5} />
          </div>
        );
    }
  };

  return (
    <div className="flex items-center" aria-live="polite" aria-label={`Backend status is ${status}`}>
      {renderIndicator()}
    </div>
  );
};

export default BackendStatusIndicator;
