import React, { useState, useEffect } from 'react';
import { type Agent, type AgentSession } from '../../models/Agent.tsx';
import { type AgentTableEntry } from '../../models/Tables.tsx';
import SidebarModal from './SideBarModal.tsx'

interface InteractAgentSideBarProps {
  isOpen: boolean;
  onClose: () => void;
  agent: Agent | null;
  onSave: (agent: Agent) => void;
}

const InteractionAgentSideBar: React.FC<InteractAgentSideBarProps> = ({ isOpen, onClose, agent, onSave }) => {
  const [formData, setFormData] = useState<AgentTableEntry>({
    id: '',
    hostname: '',
    ip: '',
    lastSeen: '',
    os: '',
    isActive: false,
    protocol: '',
  });
  const [showExtendedInfo, setShowExtendedInfo] = useState(false);
  const [sessions, setSessions] = useState<AgentSession[]>([]);

  useEffect(() => {
    if (agent) {
      setFormData(agent);
      // setSessions(agent.sessions || []);
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
      setSessions([]);
    }
  }, [agent]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? (e.target as HTMLInputElement).checked : value,
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // if (agent) {
    //   onSave({ ...formData, id: agent.id, lastSeen: new Date().toISOString() });
    // }
  };

  return (
    <SidebarModal
      isOpen={isOpen}
      onClose={onClose}
      title={agent ? 'Edit Agent' : 'Create Agent'}
      footer={
        <div className="flex items-center space-x-3">
          <button
            type="submit"
            form="edit-agent-form"
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
      }
    >
      <form id="edit-agent-form" onSubmit={handleSubmit} className="space-y-6">
        <div className="space-y-4">
          <div>
            <label htmlFor="hostname" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
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
            <label htmlFor="ip" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
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
            <label htmlFor="os" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
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
            <label htmlFor="protocol" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
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
            <label htmlFor="isActive" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
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
              <label htmlFor="isActive" className="ml-2 text-sm text-gray-900 dark:text-white">
                {formData.isActive ? 'Active' : 'Inactive'}
              </label>
            </div>
          </div>
          <div>
            <button
              type="button"
              onClick={() => setShowExtendedInfo(!showExtendedInfo)}
              className="w-full bg-gray-200 text-gray-700 py-2.5 px-5 rounded-lg hover:bg-gray-300 focus:ring-4 focus:outline-none focus:ring-gray-300 font-medium text-sm text-center dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 dark:focus:ring-gray-600"
            >
              {showExtendedInfo ? 'Hide Extended Info' : 'Show Extended Info'}
            </button>
            {showExtendedInfo && (
              <div className="mt-4 space-y-4">
                <div>
                  <h3 className="text-lg font-medium text-gray-900 dark:text-white">Sessions</h3>
                  {sessions.length > 0 ? (
                    <ul className="mt-2 space-y-2">
                      {sessions.map((session) => (
                        <li key={session.id} className="p-2 bg-gray-50 dark:bg-gray-700 rounded-lg">
                          <p><strong>Session ID:</strong> {session.session_id}</p>
                          <p><strong>Start Time:</strong> {new Date(session.start_time).toLocaleString()}</p>
                          <p><strong>End Time:</strong> {session.end_time ? new Date(session.end_time).toLocaleString() : 'Active'}</p>
                          <p><strong>IP Address:</strong> {session.ip_address}</p>
                          <p><strong>User Agent:</strong> {session.user_agent}</p>
                          <p><strong>Status:</strong> {session.is_active ? 'Active' : 'Inactive'}</p>
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <p className="text-sm text-gray-500 dark:text-gray-400">No sessions available</p>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </form>
    </SidebarModal>
  );
};

export default InteractionAgentSideBar;