import { useState } from 'react';
import { Home, Users, Shield, Briefcase, Calendar, BookOpen } from 'lucide-react';
import { Breadcrumb } from '../components/layout/Breadcrumb';
import { SidebarItem } from '../components/layout/SidebarItem';
import { SidebarAccordion } from '../components/layout/SidebarAccordion';
import { UserProfile } from '../components/user/UserProfile';

export default function Dashboard() {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  
  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  return (
    <div className="bg-gray-50 transition-all duration-300 lg:hs-overlay-layout-open:ps-65">
      {/* Main Content */}
      <main id="content">
        {/* Breadcrumb */}
        <div className="sticky top-0 inset-x-0 z-20 bg-white border-y border-gray-200 px-4 sm:px-6 lg:px-8">
          <div className="flex items-center py-2">
            {/* Navigation Toggle */}
            <button 
              type="button" 
              className="size-8 flex justify-center items-center gap-x-2 border border-gray-200 text-gray-700 hover:text-gray-500 rounded-lg focus:outline-hidden focus:text-gray-500 disabled:opacity-50 disabled:pointer-events-none" 
              aria-haspopup="dialog"
              aria-expanded={sidebarOpen}
              aria-controls="hs-application-sidebar" 
              aria-label="Toggle navigation" 
              onClick={toggleSidebar}
              data-hs-overlay="#hs-application-sidebar"
            >
              <span className="sr-only">Toggle Navigation</span>
              <svg className="shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M15 3v18"/><path d="m8 9 3 3-3 3"/></svg>
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
        </div>
        {/* End Breadcrumb */}

        {/* Sidebar */}
        <div 
          id="hs-application-sidebar" 
          className={`hs-overlay [--body-scroll:true] lg:[--overlay-backdrop:false] [--is-layout-affect:true] [--auto-close:lg]
            hs-overlay-open:translate-x-0
            -translate-x-full transition-all duration-300 transform
            w-65 h-full
            fixed inset-y-0 start-0 z-60
            bg-white border-e border-gray-200
            ${sidebarOpen ? 'translate-x-0' : ''}`}
          role="dialog" 
          tabIndex={-1} 
          aria-label="Sidebar"
        >
          <div className="relative flex flex-col h-full max-h-full">
            <div className="px-6 pt-4 flex items-center">
              {/* Logo */}
              <a className="flex-none rounded-xl text-xl inline-block font-semibold focus:outline-hidden focus:opacity-80" href="#" aria-label="BaconC2">
                <span className="text-blue-600">BaconC2</span>
              </a>
              {/* End Logo */}
            </div>

            {/* Content */}
            <div className="h-full overflow-y-auto [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-track]:bg-gray-100 [&::-webkit-scrollbar-thumb]:bg-gray-300">
              <nav className="hs-accordion-group p-3 w-full flex flex-col flex-wrap" data-hs-accordion-always-open>
                <ul className="flex flex-col space-y-1">
                  <li>
                    <SidebarItem 
                      icon={<Home size={16} />}
                      label="Dashboard"
                      href="#"
                      active={true}
                    />
                  </li>

                  <SidebarAccordion 
                    id="users-accordion"
                    icon={<Users size={16} />}
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
                    icon={<Shield size={16} />}
                    label="Account"
                    items={[
                      { label: "Link 1", href: "#" },
                      { label: "Link 2", href: "#" },
                      { label: "Link 3", href: "#" }
                    ]}
                  />

                  <SidebarAccordion 
                    id="projects-accordion"
                    icon={<Briefcase size={16} />}
                    label="Projects"
                    items={[
                      { label: "Link 1", href: "#" },
                      { label: "Link 2", href: "#" },
                      { label: "Link 3", href: "#" }
                    ]}
                  />

                  <li>
                    <SidebarItem 
                      icon={<Calendar size={16} />}
                      label="Calendar"
                      href="#"
                    />
                  </li>
                  
                  <li>
                    <SidebarItem 
                      icon={<BookOpen size={16} />}
                      label="Documentation"
                      href="#"
                    />
                  </li>
                </ul>
              </nav>
            </div>
            {/* End Content */}
            
            {/* User Profile */}
            <div className="p-4 border-t border-gray-200">
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
        <div className="w-full lg:ps-64">
          <div className="p-4 sm:p-6 space-y-4 sm:space-y-6">
            {/* Dashboard Content Here */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
                <h2 className="text-lg font-semibold mb-4 text-gray-800">Welcome to BaconC2</h2>
                <p className="text-gray-600">
                  This is your new dashboard. Start adding your content here.
                </p>
              </div>
              
              <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
                <h2 className="text-lg font-semibold mb-4 text-gray-800">Agent Statistics</h2>
                <p className="text-gray-600">
                  View your agent statistics and metrics here.
                </p>
              </div>
              
              <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
                <h2 className="text-lg font-semibold mb-4 text-gray-800">Command History</h2>
                <p className="text-gray-600">
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
