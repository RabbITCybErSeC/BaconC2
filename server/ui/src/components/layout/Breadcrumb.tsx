import React from 'react';

interface BreadcrumbItem {
  label: string;
  href?: string;
  current?: boolean;
}

interface BreadcrumbProps {
  items: BreadcrumbItem[];
}

export const Breadcrumb: React.FC<BreadcrumbProps> = ({ items }) => {
  return (
    <ol className="ms-3 flex items-center whitespace-nowrap">
      {items.map((item, index) => (
        <li 
          key={index} 
          className={`flex items-center text-sm ${
            item.current 
              ? 'font-semibold text-stone-800' 
              : 'text-stone-700'
          }`}
          {...(item.current ? { 'aria-current': 'page' } : {})}
        >
          {item.href && !item.current ? (
            <a href={item.href} className="hover:text-stone-600">
              {item.label}
            </a>
          ) : (
            item.label
          )}
          
          {index < items.length - 1 && (
            <svg className="shrink-0 mx-3 overflow-visible size-2.5 text-stone-400" width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M5 1L10.6869 7.16086C10.8637 7.35239 10.8637 7.64761 10.6869 7.83914L5 14" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
            </svg>
          )}
        </li>
      ))}
    </ol>
  );
}; 