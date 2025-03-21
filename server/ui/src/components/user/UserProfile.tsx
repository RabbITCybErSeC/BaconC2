import { useState, useRef, useEffect } from 'react';
import { User, Settings, LogOut } from 'lucide-react';

interface UserProfileProps {
  username: string;
  organization: string;
  onLogout: () => void;
}

export const UserProfile: React.FC<UserProfileProps> = ({ 
  username, 
  organization,
  onLogout
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  
  const toggleDropdown = () => {
    setIsOpen(!isOpen);
  };
  
  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };
    
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <div className="relative" ref={dropdownRef}>
      <button 
        type="button" 
        className="flex items-center gap-x-2 text-sm text-stone-700 hover:text-stone-600 focus:outline-none"
        onClick={toggleDropdown}
      >
        <span className="inline-flex items-center justify-center size-8 rounded-full bg-stone-200">
          <User className="size-4 text-stone-600" />
        </span>
        <div className="grow ms-2">
          <span className="block text-sm font-medium text-stone-800">{username}</span>
          <span className="block text-xs text-stone-500">{organization}</span>
        </div>
        <svg className="shrink-0 size-4 text-stone-600" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="m6 9 6 6 6-6"/>
        </svg>
      </button>
      
      {isOpen && (
        <div className="absolute z-10 mt-1 end-0 rounded-lg shadow-lg w-44 py-1.5 bg-white border border-stone-200">
          <a className="flex items-center gap-x-3.5 py-2 px-3 text-sm text-stone-700 hover:bg-stone-100" href="#">
            <User className="size-4" />
            Profile
          </a>
          <a className="flex items-center gap-x-3.5 py-2 px-3 text-sm text-stone-700 hover:bg-stone-100" href="#">
            <Settings className="size-4" />
            Settings
          </a>
          <div className="my-1 border-t border-stone-200"></div>
          <button 
            className="w-full flex items-center gap-x-3.5 py-2 px-3 text-sm text-stone-700 hover:bg-stone-100"
            onClick={onLogout}
          >
            <LogOut className="size-4" />
            Logout
          </button>
        </div>
      )}
    </div>
  );
}; 