import { 
  Home, 
  Users, 
  UserCircle, 
  FolderOpen, 
  Calendar, 
  BookOpen,
  ChevronDown,
  ChevronUp
} from 'lucide-react';
import { useState } from 'react';

interface NavItem {
  icon: React.ReactNode;
  label: string;
  href: string;
  children?: NavItem[];
}

const navItems: NavItem[] = [
  {
    icon: <Home className="size-4" />,
    label: 'Dashboard',
    href: '/dashboard'
  },
  {
    icon: <Users className="size-4" />,
    label: 'Users',
    href: '/users',
    children: [
      {
        icon: <Users className="size-4" />,
        label: 'All Users',
        href: '/users/all'
      },
      {
        icon: <Users className="size-4" />,
        label: 'Active Users',
        href: '/users/active'
      }
    ]
  },
  {
    icon: <UserCircle className="size-4" />,
    label: 'Account',
    href: '/account',
    children: [
      {
        icon: <UserCircle className="size-4" />,
        label: 'Profile',
        href: '/account/profile'
      },
      {
        icon: <UserCircle className="size-4" />,
        label: 'Settings',
        href: '/account/settings'
      }
    ]
  },
  {
    icon: <FolderOpen className="size-4" />,
    label: 'Projects',
    href: '/projects'
  },
  {
    icon: <Calendar className="size-4" />,
    label: 'Calendar',
    href: '/calendar'
  },
  {
    icon: <BookOpen className="size-4" />,
    label: 'Documentation',
    href: '/docs'
  }
];

export function Sidebar() {
  const [expandedItems, setExpandedItems] = useState<string[]>([]);

  const toggleItem = (label: string) => {
    setExpandedItems(prev => 
      prev.includes(label) 
        ? prev.filter(item => item !== label)
        : [...prev, label]
    );
  };

  return (
    <div 
      id="hs-application-sidebar" 
      className="hs-overlay [--body-scroll:true] lg:[--overlay-backdrop:false] [--is-layout-affect:true] [--auto-close:lg]
        hs-overlay-open:translate-x-0
        -translate-x-full transition-all duration-300 transform
        w-65 h-full
        hidden
        fixed inset-y-0 start-0 z-60
        bg-white border-e border-gray-200
        dark:bg-neutral-800 dark:border-neutral-700"
      role="dialog"
      tabIndex={-1}
      aria-label="Sidebar"
    >
      <div className="relative flex flex-col h-full max-h-full">
        <div className="px-6 pt-4 flex items-center">
          <a className="flex-none rounded-xl text-xl inline-block font-semibold focus:outline-hidden focus:opacity-80" href="/">
            <span className="text-blue-600 dark:text-white">BaconC2</span>
          </a>
        </div>

        <div className="h-full overflow-y-auto [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-track]:bg-gray-100 [&::-webkit-scrollbar-thumb]:bg-gray-300 dark:[&::-webkit-scrollbar-track]:bg-neutral-700 dark:[&::-webkit-scrollbar-thumb]:bg-neutral-500">
          <nav className="hs-accordion-group p-3 w-full flex flex-col flex-wrap">
            <ul className="flex flex-col space-y-1">
              {navItems.map((item) => (
                <li key={item.label}>
                  {item.children ? (
                    <div className="hs-accordion">
                      <button
                        type="button"
                        className="hs-accordion-toggle w-full text-start flex items-center gap-x-3.5 py-2 px-2.5 text-sm text-gray-800 rounded-lg hover:bg-gray-100 focus:outline-hidden focus:bg-gray-100 dark:bg-neutral-800 dark:hover:bg-neutral-700 dark:focus:bg-neutral-700 dark:text-neutral-200"
                        onClick={() => toggleItem(item.label)}
                      >
                        {item.icon}
                        {item.label}
                        {expandedItems.includes(item.label) ? (
                          <ChevronDown className="size-4 ms-auto" />
                        ) : (
                          <ChevronUp className="size-4 ms-auto" />
                        )}
                      </button>
                      {expandedItems.includes(item.label) && (
                        <div className="hs-accordion-content w-full overflow-hidden transition-[height] duration-300">
                          <ul className="ps-8 pt-1 space-y-1">
                            {item.children.map((child) => (
                              <li key={child.label}>
                                <a
                                  className="flex items-center gap-x-3.5 py-2 px-2.5 text-sm text-gray-800 rounded-lg hover:bg-gray-100 focus:outline-hidden focus:bg-gray-100 dark:bg-neutral-800 dark:hover:bg-neutral-700 dark:focus:bg-neutral-700 dark:text-neutral-200"
                                  href={child.href}
                                >
                                  {child.icon}
                                  {child.label}
                                </a>
                              </li>
                            ))}
                          </ul>
                        </div>
                      )}
                    </div>
                  ) : (
                    <a
                      className="flex items-center gap-x-3.5 py-2 px-2.5 text-sm text-gray-800 rounded-lg hover:bg-gray-100 focus:outline-hidden focus:bg-gray-100 dark:bg-neutral-800 dark:hover:bg-neutral-700 dark:focus:bg-neutral-700 dark:text-neutral-200"
                      href={item.href}
                    >
                      {item.icon}
                      {item.label}
                    </a>
                  )}
                </li>
              ))}
            </ul>
          </nav>
        </div>
      </div>
    </div>
  );
} 