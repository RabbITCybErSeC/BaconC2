import React, { useState, useEffect, useRef } from 'react';
import {
  Menu,
  Search,
  User,
  ChevronDown,
  LogOut,
} from 'lucide-react';
import BackendStatusIndicator from './../components/common/BackendStatusIndicator';
import SearchModal from './../components/search/SearchModal';
import { logout } from '../services/authService';

interface HeaderProps {
  userName: string;
  onSidebarToggle: () => void;
  onThemeChange?: () => void;
}

const Header: React.FC<HeaderProps> = ({
  userName,
  onSidebarToggle,
  // onThemeChange,
}) => {
  const [userMenuOpen, setUserMenuOpen] = useState<boolean>(false);
  const [modalOpen, setModalOpen] = useState<boolean>(false);
  const userMenuRef = useRef<HTMLDivElement>(null);

  // const handleThemeToggle = onThemeChange ?? (() => console.warn('Theme toggle handler (onThemeChange) not provided to Header component'));

  const handleUserMenuToggle = () => {
    setUserMenuOpen(prev => !prev);
  };

  const handleSearchClick = () => {
    setModalOpen(true);
  };

  const handleCloseModal = () => {
    setModalOpen(false);
  };

  const handleLogout = () => {
    logout();
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
        setUserMenuOpen(false);
      }
    };

    if (userMenuOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [userMenuOpen]);

  return (
    <>
      <header className="flex items-center justify-between px-6 py-3 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 relative z-10">
        <div className="flex items-center">
          <button
            onClick={onSidebarToggle}
            className="p-2 rounded-md text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-violet-500 mr-4 lg:hidden"
            aria-label="Toggle sidebar"
          >
            <Menu size={24} />
          </button>
        </div>

        <div className="flex items-center space-x-4">
          <BackendStatusIndicator />

          <button
            onClick={handleSearchClick}
            className="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 focus:outline-none focus:text-violet-600 dark:focus:text-violet-400"
            aria-label="Open search"
          >
            <Search size={20} />
          </button>
          {/**/}
          {/* <button */}
          {/*   onClick={handleThemeToggle} */}
          {/*   className="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 focus:outline-none focus:text-violet-600 dark:focus:text-violet-400" */}
          {/*   aria-label="Toggle theme" */}
          {/* > */}
          {/*   <Sun className="hidden dark:block" size={20} /> */}
          {/*   <Moon className="dark:hidden" size={20} /> */}
          {/* </button> */}

          <div className="relative" ref={userMenuRef}>
            <button
              onClick={handleUserMenuToggle}
              className="flex items-center focus:outline-none rounded-md p-1 focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-violet-500 dark:focus-visible:ring-offset-gray-800"
              aria-haspopup="true"
              aria-expanded={userMenuOpen}
              aria-label={`User menu for ${userName}`}
            >
              <span className={`mr-2 flex h-9 w-9 items-center justify-center overflow-hidden rounded-full transition-colors duration-150 ${userMenuOpen ? 'bg-violet-100 dark:bg-violet-900/30' : 'bg-gray-200 dark:bg-gray-700'} text-gray-600 dark:text-gray-400`}>
                <User size={20} className={`transition-colors duration-150 ${userMenuOpen ? 'text-violet-600 dark:text-violet-400' : 'text-gray-600 dark:text-gray-400'}`} />
              </span>
              <span className={`hidden md:block mr-1 font-medium text-sm transition-colors duration-150 ${userMenuOpen ? 'text-violet-600 dark:text-violet-400' : 'text-gray-800 dark:text-gray-100'}`}>
                {userName}
              </span>
              <ChevronDown
                className={`transition-transform duration-200 ease-in-out ${userMenuOpen ? 'rotate-180 text-violet-500 dark:text-violet-300' : 'text-gray-500 dark:text-gray-400'}`}
                size={18}
                strokeWidth={1.5}
                aria-hidden="true"
              />
            </button>

            {/* Dropdown Menu */}
            {userMenuOpen && (
              <div className="absolute right-0 mt-2 w-56 origin-top-right z-50">
                <div className="rounded-xl border border-gray-200/60 dark:border-gray-700/60 bg-white/90 dark:bg-gray-900/70 backdrop-blur supports-[backdrop-filter]:bg-white/60 dark:supports-[backdrop-filter]:bg-gray-900/50 shadow-xl shadow-black/20 ring-1 ring-black/5">
                  <div className="py-2">
                    <button
                      onClick={handleLogout}
                      className="flex items-center w-full px-4 py-2 text-sm rounded-md text-gray-700 dark:text-gray-200 hover:text-violet-600 dark:hover:text-violet-400 hover:bg-gray-100/70 dark:hover:bg-gray-700/50 transition-colors duration-150"
                    >
                      <LogOut size={16} className="mr-2" />
                      Logout
                    </button>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </header>

      <SearchModal isOpen={modalOpen} onClose={handleCloseModal} />
    </>
  );
};

export default Header;
