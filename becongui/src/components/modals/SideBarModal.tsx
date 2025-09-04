import React from 'react';
import { X } from 'lucide-react';

interface SidebarModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
  footer?: React.ReactNode;
  className?: string;
}

const SidebarModal: React.FC<SidebarModalProps> = ({
  isOpen,
  onClose,
  title,
  children,
  footer,
  className = '',
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      {/* Sidebar panel */}
      <div
        className={`relative w-full md:w-2/3 max-w-4xl border-l border-gray-200 bg-white dark:bg-gray-800 shadow-xl h-full animate-slide-in-right ${className}`}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-100">{title}</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
            aria-label="Close"
          >
            <X size={20} className="text-gray-600 dark:text-gray-300" />
          </button>
        </div>

        {/* Scrollable Content */}
        <div className="h-[calc(100%-8rem)] overflow-y-auto p-4">{children}</div>

        {/* Footer */}
        {footer && (
          <div className="p-4 border-t border-gray-200 dark:border-gray-700">
            {footer}
          </div>
        )}
      </div>
    </div>
  );
};

export default SidebarModal;
