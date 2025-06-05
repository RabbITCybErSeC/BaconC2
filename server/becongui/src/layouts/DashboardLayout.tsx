import React, { useState } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from '../partials/Sidebar';
import Header from '../partials/Header';
import Banner from '../partials/Banner';

const DashboardLayout: React.FC = () => {
  const [sidebarOpen, setSidebarOpen] = useState<boolean>(false);
  const [theme, setTheme] = useState<'light' | 'dark'>('dark');

  const toggleTheme = () => {
    setTheme((prev) => (prev === 'light' ? 'dark' : 'light'));
    document.documentElement.classList.toggle('dark');
  };

  return (
    <div
      className={`flex h-screen overflow-hidden ${theme === 'light' ? 'bg-gray-100' : 'bg-slate-800'
        }`} // Use theme here
    >
      {/* Sidebar component */}
      <Sidebar sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />

      {/* Content area */}
      <div className="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
        {/* Header component */}
        <Header
          userName="User" // Replace with actual user data
          onSidebarToggle={() => setSidebarOpen((prev) => !prev)}
          onThemeChange={toggleTheme}
        />

        <main className="grow">
          <Outlet />
        </main>

        {/* Banner component */}
        <Banner />
      </div>
    </div>
  );
};

export default DashboardLayout;
