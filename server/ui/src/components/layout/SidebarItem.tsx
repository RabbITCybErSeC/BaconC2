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
        ? 'bg-stone-100 text-stone-800'
        : 'text-stone-700 hover:bg-stone-100 focus:bg-stone-100'
      }`}
      href={href}
    >
      <span className="shrink-0 size-4">{icon}</span>
      {label}
    </a>
  );
}; 