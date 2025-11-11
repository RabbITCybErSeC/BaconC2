import React, { useState, useEffect } from 'react';
import ErrorBanner from '../partials/ErrorBanner';
import { getToken } from '../services/authService';

interface PluginMetadata {
  name: string;
  version: string;
  author: string;
  description: string;
  capabilities: string[];
  dependencies: string[];
  loaded_at: string;
  status: string;
  file_path?: string;
  hash?: string;
}

interface PluginInfo {
  metadata: PluginMetadata;
  file_path: string;
  file_size: number;
  hash: string;
}

const PluginsPage: React.FC = () => {
  const [plugins, setPlugins] = useState<PluginInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPlugins = async () => {
      try {
        const response = await fetch('/api/v1/frontend/plugins', {
          headers: {
            'Authorization': getToken() || '',
            'Content-Type': 'application/json',
          },
        });
        if (!response.ok) {
          throw new Error(`Failed to fetch plugins: ${response.statusText}`);
        }
        const data = await response.json();
        setPlugins(data || []);
        setError(null);
      } catch (err) {
        setError('Error fetching plugins. Please try again later.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    fetchPlugins();
  }, []);

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const getStatusColor = (status: string): string => {
    switch (status.toLowerCase()) {
      case 'loaded':
      case 'active':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300';
      case 'loading':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300';
      case 'error':
        return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300';
    }
  };

  return (
    <div className="px-4 sm:px-6 lg:px-8 py-8 w-full max-w-9xl mx-auto">
      {/* Page header */}
      <div className="sm:flex sm:justify-between sm:items-center mb-8">
        <div>
          <h1 className="text-2xl md:text-3xl text-gray-800 dark:text-gray-100 font-bold">Plugins</h1>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            Manage and monitor loaded plugins
          </p>
        </div>
        <div className="grid grid-flow-col sm:auto-cols-max justify-start sm:justify-end gap-2">
          {/* Future: Add plugin upload button */}
        </div>
      </div>

      {/* Error banner */}
      {error && <ErrorBanner message={error} onDismiss={() => setError(null)} />}

      {/* Loading state */}
      {loading && (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 dark:border-gray-100"></div>
          <p className="mt-4 text-gray-600 dark:text-gray-400">Loading plugins...</p>
        </div>
      )}

      {/* Empty state */}
      {!loading && plugins.length === 0 && (
        <div className="bg-white dark:bg-slate-800 shadow-md rounded-sm border border-gray-200 dark:border-slate-700 p-8 text-center">
          <svg
            className="w-16 h-16 mx-auto text-gray-400 dark:text-gray-600 mb-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
            />
          </svg>
          <h3 className="text-lg font-semibold text-gray-800 dark:text-gray-100 mb-2">
            No plugins loaded
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            No plugins are currently loaded in the system.
          </p>
        </div>
      )}

      {/* Plugin grid */}
      {!loading && plugins.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {plugins.map((plugin) => (
            <div
              key={plugin.hash}
              className="bg-white dark:bg-slate-800 shadow-md rounded-lg border border-gray-200 dark:border-slate-700 overflow-hidden hover:shadow-lg transition-shadow duration-200"
            >
              {/* Plugin header */}
              <div className="p-6 border-b border-gray-200 dark:border-slate-700">
                <div className="flex items-start justify-between mb-2">
                  <h3 className="text-lg font-semibold text-gray-800 dark:text-gray-100 truncate">
                    {plugin.metadata.name}
                  </h3>
                  <span
                    className={`px-2 py-1 text-xs font-medium rounded-full ${getStatusColor(
                      plugin.metadata.status
                    )}`}
                  >
                    {plugin.metadata.status}
                  </span>
                </div>
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
                  {plugin.metadata.description}
                </p>
                <div className="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-500">
                  <span>v{plugin.metadata.version}</span>
                  <span>•</span>
                  <span>{plugin.metadata.author}</span>
                </div>
              </div>

              {/* Plugin details */}
              <div className="p-6 space-y-4">
                {/* Capabilities */}
                {plugin.metadata.capabilities && plugin.metadata.capabilities.length > 0 && (
                  <div>
                    <h4 className="text-xs font-semibold text-gray-700 dark:text-gray-300 uppercase mb-2">
                      Capabilities
                    </h4>
                    <div className="flex flex-wrap gap-2">
                      {plugin.metadata.capabilities.map((cap, idx) => (
                        <span
                          key={idx}
                          className="px-2 py-1 text-xs bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300 rounded"
                        >
                          {cap}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Dependencies */}
                {plugin.metadata.dependencies && plugin.metadata.dependencies.length > 0 && (
                  <div>
                    <h4 className="text-xs font-semibold text-gray-700 dark:text-gray-300 uppercase mb-2">
                      Dependencies
                    </h4>
                    <div className="flex flex-wrap gap-2">
                      {plugin.metadata.dependencies.map((dep, idx) => (
                        <span
                          key={idx}
                          className="px-2 py-1 text-xs bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300 rounded"
                        >
                          {dep}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* File info */}
                <div className="pt-4 border-t border-gray-200 dark:border-slate-700">
                  <div className="space-y-2 text-xs">
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">File Size:</span>
                      <span className="text-gray-800 dark:text-gray-200 font-medium">
                        {formatFileSize(plugin.file_size)}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">Loaded:</span>
                      <span className="text-gray-800 dark:text-gray-200 font-medium">
                        {formatDate(plugin.metadata.loaded_at)}
                      </span>
                    </div>
                    {plugin.hash && (
                      <div className="flex flex-col gap-1">
                        <span className="text-gray-600 dark:text-gray-400">Hash:</span>
                        <code className="text-xs text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-900 p-1 rounded break-all">
                          {plugin.hash.substring(0, 16)}...
                        </code>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default PluginsPage;
