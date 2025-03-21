import { useState, useEffect } from 'react';
import { Home, Users, Shield, Briefcase, Calendar, BookOpen, Menu, Sun, Moon } from 'lucide-react';
import { Breadcrumb } from '../components/layout/Breadcrumb';
import { SidebarItem } from '../components/layout/SidebarItem';
import { SidebarAccordion } from '../components/layout/SidebarAccordion';
import { UserProfile } from '../components/user/UserProfile';

export default function Dashboard() {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [darkMode, setDarkMode] = useState(() => {
    // Check if user has a preference stored in localStorage
    const savedTheme = localStorage.getItem('theme');
    // Check if user prefers dark mode via system preference
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    return savedTheme === 'dark' || (!savedTheme && prefersDark);
  });
  
  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };
  
  const toggleDarkMode = () => {
    setDarkMode(!darkMode);
  };
  
  // Update body classes when sidebar state or dark mode changes
  useEffect(() => {
    if (sidebarOpen) {
      document.body.classList.add('sidebar-open');
    } else {
      document.body.classList.remove('sidebar-open');
    }
    
    if (darkMode) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  }, [sidebarOpen, darkMode]);

  return (
    <div className={`bg-stone-50 dark:bg-stone-900 transition-all duration-300 ${sidebarOpen ? 'lg:pl-64' : ''}`}>
      {/* Main Content */}
      <main id="content">
        {/* Breadcrumb */}
        <div className="sticky top-0 inset-x-0 z-20 bg-white dark:bg-stone-800 border-y border-stone-200 dark:border-stone-700 px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between py-2">
            <div className="flex items-center">
              {/* Navigation Toggle */}
              <button 
                type="button" 
                className="size-8 flex justify-center items-center gap-x-2 border border-stone-200 dark:border-stone-700 text-stone-700 dark:text-stone-300 hover:text-stone-500 dark:hover:text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-stone-500 disabled:opacity-50 disabled:pointer-events-none" 
                aria-expanded={sidebarOpen}
                aria-controls="hs-application-sidebar" 
                aria-label="Toggle navigation" 
                onClick={toggleSidebar}
              >
                <span className="sr-only">Toggle Navigation</span>
                <Menu className="shrink-0 size-4" />
              </button>
              {/* End Navigation Toggle */}
              
              {/* Breadcrumb */}
              <Breadcrumb 
                items={[
                  { label: 'Application Layout', href: '/' },
                  { label: 'Dashboard', href: '#', current: true }
                ]} 
              />
              {/* End Breadcrumb */}
            </div>
            
            {/* Dark Mode Toggle */}
            <button
              type="button"
              className="size-8 flex justify-center items-center gap-x-2 border border-stone-200 dark:border-stone-700 text-stone-700 dark:text-stone-300 hover:text-stone-500 dark:hover:text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-stone-500 disabled:opacity-50 disabled:pointer-events-none"
              aria-label="Toggle dark mode"
              onClick={toggleDarkMode}
            >
              <span className="sr-only">Toggle Dark Mode</span>
              {darkMode ? (
                <Sun className="shrink-0 size-4" />
              ) : (
                <Moon className="shrink-0 size-4" />
              )}
            </button>
            {/* End Dark Mode Toggle */}
          </div>
        </div>
        {/* End Breadcrumb */}

        {/* Sidebar */}
        <div 
          id="hs-application-sidebar" 
          className={`transition-all duration-300 transform
            w-64 h-full
            fixed inset-y-0 start-0 z-50
            bg-white dark:bg-stone-800 border-e border-stone-200 dark:border-stone-700
            ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}
          role="navigation" 
          aria-label="Sidebar"
        >
          <div className="relative flex flex-col h-full max-h-full">
            <div className="px-6 pt-4 flex items-center">
              {/* Logo */}
              <a className="flex-none rounded-xl text-xl inline-block font-semibold focus:outline-none focus:opacity-80" href="#" aria-label="BaconC2">
                <span className="text-stone-600 dark:text-stone-400">BaconC2</span>           </a>
              {/* End Logo */}
            </div>

            {/* Content */}
            <div className="h-full overflow-y-auto [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-track]:bg-stone-100 dark:bg-stone-800 [&::-webkit-scrollbar-thumb]:bg-stone-300 dark:[&::-webkit-scrollbar-thumb]:bg-stone-600">
              <nav className="hs-accordion-group p-3 w-full flex flex-col flex-wrap" data-hs-accordion-always-open>
                <ul className="flex flex-col space-y-1">
                  <li>
                    <SidebarItem 
                      icon={<Home size={16} className="text-stone-600 dark:text-stone-400" />}
                      label="Dashboard"
                      href="#"
                      active={true}
                    />
                  </li>

                  <SidebarAccordion 
                    id="users-accordion"
                    icon={<Users size={16} className="text-stone-600 dark:text-stone-400" />}
                    label="Users"
                    items={[
                      { 
                        id: "users-accordion-sub-1",
                        label: "Sub Menu 1",
                        items: [
                          { label: "Link 1", href: "#" },
                          { label: "Link 2", href: "#" },
                          { label: "Link 3", href: "#" }
                        ]
                      },
                      {
                        id: "users-accordion-sub-2",
                        label: "Sub Menu 2",
                        items: [
                          { label: "Link 1", href: "#" },
                          { label: "Link 2", href: "#" },
                          { label: "Link 3", href: "#" }
                        ]
                      }
                    ]}
                  />

                  <SidebarAccordion 
                    id="account-accordion"
                    icon={<Shield size={16} className="text-stone-600 dark:text-stone-400" />}
                    label="Account"
                    items={[
                      { label: "Link 1", href: "#" },
                      { label: "Link 2", href: "#" },
                      { label: "Link 3", href: "#" }
                    ]}
                  />

                  <SidebarAccordion 
                    id="projects-accordion"
                    icon={<Briefcase size={16} className="text-stone-600 dark:text-stone-400" />}
                    label="Projects"
                    items={[
                      { label: "Link 1", href: "#" },
                      { label: "Link 2", href: "#" },
                      { label: "Link 3", href: "#" }
                    ]}
                  />

                  <li>
                    <SidebarItem 
                      icon={<Calendar size={16} className="text-stone-600 dark:text-stone-400" />}
                      label="Calendar"
                      href="#"
                    />
                  </li>
                  
                  <li>
                    <SidebarItem 
                      icon={<BookOpen size={16} className="text-stone-600 dark:text-stone-400" />}
                      label="Documentation"
                      href="#"
                    />
                  </li>
                </ul>
              </nav>
            </div>
            {/* End Content */}
            
            {/* User Profile */}
            <div className="p-4 border-t border-stone-200 dark:border-stone-700">
              <UserProfile 
                username="admin" 
                organization="Administrator"
                onLogout={() => window.location.href = '/'}
              />
            </div>
            {/* End User Profile */}
          </div>
        </div>
        {/* End Sidebar */}

        {/* Content */}
        <div className="w-full transition-all duration-300">
          <div className="p-4 sm:p-6 space-y-4 sm:space-y-6">
            {/* Dashboard Content Here */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div className="bg-white dark:bg-stone-800 p-6 rounded-lg shadow-sm border border-stone-200 dark:border-stone-700">
                <h2 className="text-lg font-semibold mb-4 text-stone-800 dark:text-white">Welcome to BaconC2</h2>
                <p className="text-stone-600 dark:text-stone-300">
                  This is your new dashboard. Start adding your content here.
                </p>
              </div>
              
              <div className="bg-white dark:bg-stone-800 p-6 rounded-lg shadow-sm border border-stone-200 dark:border-stone-700">
                <h2 className="text-lg font-semibold mb-4 text-stone-800 dark:text-white">Agent Statistics</h2>
                <p className="text-stone-600 dark:text-stone-300">
                  View your agent statistics and metrics here.
                </p>
              </div>
              
              <div className="bg-white dark:bg-stone-800 p-6 rounded-lg shadow-sm border border-stone-200 dark:border-stone-700">
                <h2 className="text-lg font-semibold mb-4 text-stone-800 dark:text-white">Command History</h2>
                <p className="text-stone-600 dark:text-stone-300">
                  Your recent command history will appear here.
                </p>
              </div>
            </div>
          </div>
        </div>
        {/* End Content */}
      </main>
      {/* End Main Content */}
    </div>
  );
}
