import { ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  return (
    <div className="bg-gray-50 transition-all duration-300 lg:hs-overlay-layout-open:ps-65 dark:bg-neutral-900">
      <main id="content" className="w-full lg:ps-64">
        <div className="p-4 sm:p-6 space-y-4 sm:space-y-6">
          {children}
        </div>
      </main>
    </div>
  );
} 