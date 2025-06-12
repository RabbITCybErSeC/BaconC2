import React, { useState } from 'react';
import { X } from 'lucide-react';

interface ErrorBannerProps {
  message: string;
  onDismiss: () => void;
}

const ErrorBanner: React.FC<ErrorBannerProps> = ({ message, onDismiss }) => {
  const [isOpen, setIsOpen] = useState<boolean>(true);

  const handleDismiss = () => {
    setIsOpen(false);
    onDismiss();
  };

  return (
    <>
      {isOpen && (
        <div className="fixed bottom-0 right-0 w-full md:bottom-8 md:right-12 md:w-auto z-50">
          <div className="bg-red-800 border border-transparent dark:border-red-700/60 text-red-50 text-sm p-3 md:rounded-sm shadow-lg flex justify-between items-center">
            <span>{message}</span>
            <button
              className="text-red-400 hover:text-red-300 pl-2 ml-3 border-l border-red-700/60"
              onClick={handleDismiss}
              aria-label="Dismiss error"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}
    </>
  );
};

export default ErrorBanner;