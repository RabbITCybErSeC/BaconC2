import React, { ReactNode } from 'react';

interface SidebarItemProps {
  icon: ReactNode;
  label: string;
  href: string;
  active?: boolean;
}

export const SidebarItem: React.FC<SidebarItemProps> = ({ 
  icon, 
  label, 
  href, 
  active = false 
}) => {
  return (
    <a 
      className={`flex items-center gap-x-3.5 py-2 px-2.5 text-sm rounded-lg focus:outline-hidden 
      ${active 
        ? 'bg-gray-100 text-gray-800'
        : 'text-gray-700 hover:bg-gray-100 focus:bg-gray-100'
      }`}
      href={href}
    >
      <span className="shrink-0 size-4">{icon}</span>
      {label}
    </a>
  );
}; 