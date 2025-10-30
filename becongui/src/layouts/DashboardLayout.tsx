import React, { useState, useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from '../partials/Sidebar';
import Header from '../partials/Header';
import Banner from '../partials/Banner';

const DashboardLayout: React.FC = () => {
  const [sidebarOpen, setSidebarOpen] = useState<boolean>(false);
  
  // Always use dark theme
  useEffect(() => {
    // Add dark class to html element
    const html = document.documentElement;
    html.classList.add('dark');
    // Set data-theme attribute for any components that might use it
    html.setAttribute('data-theme', 'dark');
    // Ensure background color is set
    document.body.classList.add('bg-slate-900');
  }, []);

  return (
    <div className="flex h-screen overflow-hidden bg-slate-800">
      {/* Sidebar component */}
      <Sidebar sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />

      {/* Content area */}
      <div className="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
        {/* Header component */}
        <Header
          userName="User" // Replace with actual user data
          onSidebarToggle={() => setSidebarOpen((prev) => !prev)}
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
