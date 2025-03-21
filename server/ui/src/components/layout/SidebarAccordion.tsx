import React, { ReactNode, useState } from 'react';
import { ChevronUp, ChevronDown } from 'lucide-react';

interface AccordionItem {
  id?: string;
  label: string;
  href?: string;
  items?: AccordionItem[];
}

interface SidebarAccordionProps {
  id: string;
  icon: ReactNode;
  label: string;
  items: AccordionItem[];
}

export const SidebarAccordion: React.FC<SidebarAccordionProps> = ({ 
  id,
  icon,
  label,
  items
}) => {
  const [isOpen, setIsOpen] = useState(false);

  const toggleAccordion = () => {
    setIsOpen(!isOpen);
  };

  const renderItems = (accordionItems: AccordionItem[], level: number = 0) => {
    return (
      <ul className={`${level > 0 ? 'ps-8' : ''} pt-1 space-y-1`}>
        {accordionItems.map((item, index) => {
          if (item.items) {
            // This is a nested accordion
            return (
              <li key={item.id || `${id}-sub-${index}`} className="hs-accordion">
                <button 
                  type="button" 
                  className="hs-accordion-toggle w-full text-start flex items-center gap-x-3.5 py-2 px-2.5 text-sm text-stone-700 rounded-lg hover:bg-stone-100 focus:outline-hidden focus:bg-stone-100"
                  onClick={() => {
                    // Handle nested accordion
                  }}
                >
                  {item.label}
                  <ChevronUp className="hs-accordion-active:block ms-auto hidden size-4" />
                  <ChevronDown className="hs-accordion-active:hidden ms-auto block size-4" />
                </button>
                
                <div className="hs-accordion-content w-full overflow-hidden transition-[height] duration-300 hidden">
                  {renderItems(item.items, level + 1)}
                </div>
              </li>
            );
          } else {
            // This is a regular item
            return (
              <li key={`${id}-item-${index}`}>
                <a 
                  className="flex items-center gap-x-3.5 py-2 px-2.5 text-sm text-stone-700 rounded-lg hover:bg-stone-100 focus:outline-hidden focus:bg-stone-100" 
                  href={item.href || '#'}
                >
                  {item.label}
                </a>
              </li>
            );
          }
        })}
      </ul>
    );
  };

  return (
    <li className="hs-accordion">
      <button 
        type="button" 
        className="hs-accordion-toggle w-full text-start flex items-center gap-x-3.5 py-2 px-2.5 text-sm text-stone-700 rounded-lg hover:bg-stone-100 focus:outline-hidden focus:bg-stone-100"
        onClick={toggleAccordion}
      >
        <span className="shrink-0 size-4">{icon}</span>
        {label}
        
        {isOpen ? (
          <ChevronUp className="ms-auto size-4" />
        ) : (
          <ChevronDown className="ms-auto size-4" />
        )}
      </button>

      <div 
        className={`w-full overflow-hidden transition-all duration-300 ${isOpen ? 'block' : 'hidden'}`}
      >
        {renderItems(items)}
      </div>
    </li>
  );
}; 