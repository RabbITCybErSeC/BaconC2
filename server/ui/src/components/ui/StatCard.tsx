import React, { ReactNode } from 'react';

interface StatCardProps {
  title: string;
  value: number | string;
  subtitle?: string;
  icon?: ReactNode;
}

export const StatCard: React.FC<StatCardProps> = ({ 
  title, 
  value, 
  subtitle,
  icon
}) => {
  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-500">{title}</p>
          <h3 className="text-2xl font-bold text-gray-800 mt-1">{value}</h3>
          {subtitle && (
            <p className="text-xs text-gray-500 mt-1">{subtitle}</p>
          )}
        </div>
        {icon && (
          <div className="size-10 flex items-center justify-center bg-gray-100 rounded-full">
            {icon}
          </div>
        )}
      </div>
    </div>
  );
}; 