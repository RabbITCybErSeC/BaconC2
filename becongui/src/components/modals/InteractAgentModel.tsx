import React, { useState, useEffect } from 'react';
import { type Agent } from '../../models/Agent.tsx';
import SidebarModal from './SideBarModal.tsx';

interface InteractAgentSideBarProps {
  isOpen: boolean;
  onClose: () => void;
  agent: Agent | null;
}

type TabOption = 'terminal' | 'extended' | 'actions';

const InteractionAgentSideBar: React.FC<InteractAgentSideBarProps> = ({ isOpen, onClose, agent }) => {
  const [terminalOutput, setTerminalOutput] = useState<string[]>([]);
  const [activeTab, setActiveTab] = useState<TabOption>('terminal');

  useEffect(() => {
    if (agent) {
      // Reset terminal when new agent is selected
      setTerminalOutput(["[+] Connected to agent...", "[+] Awaiting instructions..."]); 
      setActiveTab('terminal');
    } else {
      setTerminalOutput([]);
    }
  }, [agent]);

  const handleDummyInstruction = () => {
    setTerminalOutput(prev => [...prev, "[>] Sending instruction: Dump Credentials", "[+] Response: credentials dumped (dummy)"]);
  };

  return (
    <SidebarModal
      isOpen={isOpen}
      onClose={onClose}
      title={agent ? `Agent: ${agent.hostname}` : 'No Agent Selected'}
    >
      {agent ? (
        <div className="space-y-6">
          {/* General Info */}
          <div className="space-y-2">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">General Info</h3>
            <p><strong>Hostname:</strong> {agent.hostname}</p>
            <p><strong>IP Address:</strong> {agent.ip}</p>
            <p><strong>OS:</strong> {agent.os}</p>
            <p><strong>Protocol:</strong> {agent.protocol}</p>
            <p><strong>Status:</strong> {agent.isActive ? 'Active' : 'Inactive'}</p>
            <p><strong>Last Seen:</strong> {new Date(agent.lastSeen).toLocaleString()}</p>
          </div>

          {/* Tabs */}
          <div>
            <div className="flex border-b border-gray-300 dark:border-gray-700">
              <button
                className={`px-4 py-2 text-sm font-medium ${activeTab === 'terminal' ? 'border-b-2 border-blue-600 text-blue-600 dark:text-blue-400' : 'text-gray-600 dark:text-gray-300'}`}
                onClick={() => setActiveTab('terminal')}
              >
                Terminal
              </button>
              <button
                className={`px-4 py-2 text-sm font-medium ${activeTab === 'extended' ? 'border-b-2 border-blue-600 text-blue-600 dark:text-blue-400' : 'text-gray-600 dark:text-gray-300'}`}
                onClick={() => setActiveTab('extended')}
              >
                Extended Info
              </button>
              <button
                className={`px-4 py-2 text-sm font-medium ${activeTab === 'actions' ? 'border-b-2 border-blue-600 text-blue-600 dark:text-blue-400' : 'text-gray-600 dark:text-gray-300'}`}
                onClick={() => setActiveTab('actions')}
              >
                Actions
              </button>
            </div>

            <div className="mt-4">
              {activeTab === 'terminal' && (
                <>
                  {/* Terminal Output */}
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">Terminal Output</h3>
                    <div className="bg-black text-green-400 font-mono text-sm p-3 rounded-lg h-64 overflow-y-auto">
                      {terminalOutput.map((line, idx) => (
                        <div key={idx}>{line}</div>
                      ))}
                    </div>
                  </div>
                </>
              )}

              {activeTab === 'extended' && (
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">Extended Info</h3>
                  <p className="text-gray-600 dark:text-gray-300 text-sm">Extended agent details will appear here (to be implemented).</p>
                </div>
              )}

              {activeTab === 'actions' && (
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">Actions</h3>
                  <button
                    type="button"
                    onClick={handleDummyInstruction}
                    className="w-full bg-blue-600 text-white py-2.5 px-5 rounded-lg hover:bg-blue-700 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium text-sm text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
                  >
                    Dump Credentials (Dummy)
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      ) : (
        <p className="text-gray-500 dark:text-gray-400">No agent selected</p>
      )}
    </SidebarModal>
  );
};

export default InteractionAgentSideBar;
