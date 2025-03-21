import { Menu } from 'lucide-react';

interface HeaderProps {
  title: string;
  breadcrumbs: {
    label: string;
    href?: string;
  }[];
}

export function Header({ title, breadcrumbs }: HeaderProps) {
  return (
    <div className="sticky top-0 inset-x-0 z-20 bg-white border-y border-gray-200 px-4 sm:px-6 lg:px-8 dark:bg-neutral-800 dark:border-neutral-700">
      <div className="flex items-center py-2">
        {/* Navigation Toggle */}
        <button
          type="button"
          className="size-8 flex justify-center items-center gap-x-2 border border-gray-200 text-gray-800 hover:text-gray-500 rounded-lg focus:outline-hidden focus:text-gray-500 disabled:opacity-50 disabled:pointer-events-none dark:border-neutral-700 dark:text-neutral-200 dark:hover:text-neutral-500 dark:focus:text-neutral-500"
          aria-haspopup="dialog"
          aria-expanded="false"
          aria-controls="hs-application-sidebar"
          aria-label="Toggle navigation"
          data-hs-overlay="#hs-application-sidebar"
        >
          <span className="sr-only">Toggle Navigation</span>
          <Menu className="size-4" />
        </button>

        {/* Breadcrumb */}
        <ol className="ms-3 flex items-center whitespace-nowrap">
          {breadcrumbs.map((crumb, index) => (
            <li key={crumb.label} className="flex items-center text-sm text-gray-800 dark:text-neutral-400">
              {index > 0 && (
                <svg
                  className="shrink-0 mx-3 overflow-visible size-2.5 text-gray-400 dark:text-neutral-500"
                  width="16"
                  height="16"
                  viewBox="0 0 16 16"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M5 1L10.6869 7.16086C10.8637 7.35239 10.8637 7.64761 10.6869 7.83914L5 14"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                  />
                </svg>
              )}
              {crumb.href ? (
                <a
                  href={crumb.href}
                  className="hover:text-gray-500 dark:hover:text-neutral-300"
                >
                  {crumb.label}
                </a>
              ) : (
                <span className="font-semibold text-gray-800 dark:text-neutral-400">
                  {crumb.label}
                </span>
              )}
            </li>
          ))}
        </ol>
      </div>
    </div>
  );
} 