import React from "react";
import type { ReactElement } from "react";
import type { SVGProps } from "react";

import {
  Home,
  Server,
  Network,
  Settings,
  UserCircle,
  Database,
  Terminal,
  FileText,
  Puzzle
} from "lucide-react";

export interface NavItem {
  icon: ReactElement;
  label: string;
  href: string;
  children?: NavItem[];
}

const iconClasses = "shrink-0 size-5";

const cloneIcon = (icon: ReactElement<SVGProps<SVGSVGElement>>) =>
  React.cloneElement(icon, { className: iconClasses });

export const navItems: NavItem[] = [
  {
    icon: cloneIcon(<Home />),
    label: 'Dashboard',
    href: '/dashboard'
  },
  {
    icon: cloneIcon(<Server />),
    label: 'Agents',
    href: '/agents',
    children: [
      {
        icon: cloneIcon(<Server />),
        label: 'All Agents',
        href: '/agents/all'
      },
      {
        icon: cloneIcon(<Server />),
        label: 'Active Agents',
        href: '/agents/active'
      },
    ]
  },
  {
    icon: cloneIcon(<Network />),
    label: 'C2 Infrastructure',
    href: '/infrastructure'
  },
  {
    icon: cloneIcon(<Terminal />),
    label: 'Commands',
    href: '/commands'
  },
  {
    icon: cloneIcon(<Database />),
    label: 'Data Collection',
    href: '/data'
  },
  {
    icon: cloneIcon(<Puzzle />),
    label: 'Plugins',
    href: '/plugins'
  },
  {
    icon: cloneIcon(<FileText />),
    label: 'Logs',
    href: '/logs'
  },
  {
    icon: cloneIcon(<Settings />),
    label: 'Settings',
    href: '/settings',
    children: [
      {
        icon: cloneIcon(<UserCircle />),
        label: 'Profile',
        href: '/settings/profile'
      },
      {
        icon: cloneIcon(<Settings />),
        label: 'General',
        href: '/settings/general'
      },
      {
        icon: cloneIcon(<Network />),
        label: 'Network Settings',
        href: '/settings/network'
      }
    ]
  },
];
