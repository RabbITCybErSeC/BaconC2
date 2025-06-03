import React, { useState, useEffect } from 'react';
import { type Agent } from '../../models/Agent';

interface EditAgentModalProps {
  isOpen: boolean;
  onClose: () => void;
  agent: Agent | null;
  onSave: (agent: Agent) => void;
}

const EditAgentModal: React.FC<EditAgentModalProps> = ({ isOpen, onClose, agent, onSave }) => {
  const [formData, setFormData] = useState<Agent>({
    id: '',
    hostname: '',
    ip: '',
    lastSeen: '',
    os: '',
    isActive: false,
    protocol: '',
  });

  useEffect(() => {
    if (agent) {
      setFormData(agent);
    } else {
      setFormData({
        id: '',
        hostname: '',
        ip: '',
        lastSeen: '',
        os: '',
        isActive: false,
        protocol: '',
      });
    }
  }, [agent]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value, type } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? (e.target as HTMLInputElement).checked : value,
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (agent) {
      // Update existing agent
      onSave({ ...formData, id: agent.id, lastSeen: new Date().toISOString() });
    }
  };

  if (!isOpen) return null;

  return (
    <>
      <div className="fixed inset-0 bg-gray-600 bg-opacity-50 z-50 flex items-center justify-center">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-md p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
              {agent ? 'Edit Agent' : 'Create Agent'}
            </h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
            >
              <svg
                className="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
          <form onSubmit={handleSubmit}>
            <div className="space-y-4">
              <div>
                <label
                  htmlFor="hostname"
                  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Hostname
                </label>
                <input
                  type="text"
                  name="hostname"
                  id="hostname"
                  value={formData.hostname}
                  onChange={handleChange}
                  className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  placeholder="Enter hostname..."
                  required
                />
              </div>
              <div>
                <label
                  htmlFor="ip"
                  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  IP Address
                </label>
                <input
                  type="text"
                  name="ip"
                  id="ip"
                  value={formData.ip}
                  onChange={handleChange}
                  className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  placeholder="Enter IP address..."
                  required
                />
              </div>
              <div>
                <label
                  htmlFor="os"
                  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Operating System
                </label>
                <input
                  type="text"
                  name="os"
                  id="os"
                  value={formData.os}
                  onChange={handleChange}
                  className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  placeholder="Enter OS..."
                />
              </div>
              <div>
                <label
                  htmlFor="protocol"
                  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Protocol
                </label>
                <input
                  type="text"
                  name="protocol"
                  id="protocol"
                  value={formData.protocol}
                  onChange={handleChange}
                  className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  placeholder="Enter protocol..."
                />
              </div>
              <div>
                <label
                  htmlFor="isActive"
                  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Status
                </label>
                <div className="flex items-center">
                  <input
                    type="checkbox"
                    name="isActive"
                    id="isActive"
                    checked={formData.isActive}
                    onChange={handleChange}
                    className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                  />
                  <label
                    htmlFor="isActive"
                    className="ml-2 text-sm text-gray-900 dark:text-white"
                  >
                    {formData.isActive ? 'Active' : 'Inactive'}
                  </label>
                </div>
              </div>
            </div>
            <div className="flex items-center space-x-3 pt-4">
              <button
                type="submit"
                className="w-full bg-blue-600 text-white py-2.5 px-5 rounded-lg hover:bg-blue-700 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium text-sm text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
              >
                Save Changes
              </button>
              <button
                type="button"
                onClick={onClose}
                className="w-full bg-gray-200 text-gray-700 py-2.5 px-5 rounded-lg hover:bg-gray-300 focus:ring-4 focus:outline-none focus:ring-gray-300 font-medium text-sm text-center dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 dark:focus:ring-gray-600"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </>
  );
};

export default EditAgentModal;
