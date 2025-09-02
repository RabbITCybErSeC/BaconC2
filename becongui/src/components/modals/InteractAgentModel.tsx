import React, { useState, useEffect } from 'react';
import { type Agent } from '../../models/Agent.tsx';
import SidebarModal from './SideBarModal.tsx';
import { Terminal, Info, Zap } from 'lucide-react';

interface InteractAgentSideBarProps {
  isOpen: boolean;
  onClose: () => void;
  agent: Agent | null;
}

type TabOption = 'terminal' | 'extended' | 'actions';

const InteractionAgentSideBar: React.FC<InteractAgentSideBarProps> = ({
  isOpen,
  onClose,
  agent,
}) => {
  const [terminalOutput, setTerminalOutput] = useState<string[]>([]);
  const [activeTab, setActiveTab] = useState<TabOption>('terminal');

  useEffect(() => {
    if (agent) {
      setTerminalOutput(['[+] Connected to agent...', '[+] Awaiting instructions...']);
      setActiveTab('terminal');
    } else {
      setTerminalOutput([]);
    }
  }, [agent]);

  const handleDummyInstruction = () => {
    setTerminalOutput((prev) => [
      ...prev,
      '[>] Sending instruction: Dump Credentials',
      '[+] Response: credentials dumped (dummy)',
    ]);
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
            <h3 className="text-sm uppercase text-gray-400 dark:text-gray-500 font-semibold">
              General Info
            </h3>
            <div className="text-gray-800 dark:text-gray-200 text-sm space-y-1">
              <p>
                <strong>Hostname:</strong> {agent.hostname}
              </p>
              <p>
                <strong>IP Address:</strong> {agent.ip}
              </p>
              <p>
                <strong>OS:</strong> {agent.os}
              </p>
              <p>
                <strong>Protocol:</strong> {agent.protocol}
              </p>
              <p>
                <strong>Status:</strong> {agent.isActive ? 'Active' : 'Inactive'}
              </p>
              <p>
                <strong>Last Seen:</strong> {new Date(agent.lastSeen).toLocaleString()}
              </p>
            </div>
          </div>

          {/* Tabs */}
          <div>
            <div className="flex border-b border-gray-200 dark:border-gray-700">
              {[
                { key: 'terminal', label: 'Terminal', icon: <Terminal className="w-4 h-4 mr-2" /> },
                { key: 'extended', label: 'Extended Info', icon: <Info className="w-4 h-4 mr-2" /> },
                { key: 'actions', label: 'Actions', icon: <Zap className="w-4 h-4 mr-2" /> },
              ].map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => setActiveTab(tab.key as TabOption)}
                  className={`flex items-center px-3 py-2 text-sm rounded-t-lg transition-colors
                    ${
                      activeTab === tab.key
                        ? 'text-violet-600 dark:text-violet-400 bg-violet-100 dark:bg-violet-900/40 font-medium'
                        : 'text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700/60'
                    }`}
                >
                  {tab.icon}
                  {tab.label}
                </button>
              ))}
            </div>

            <div className="mt-4">
              {activeTab === 'terminal' && (
                <div>
                  <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2">
                    Terminal Output
                  </h3>
                  <div className="bg-black dark:bg-gray-900 text-green-400 font-mono text-sm p-3 rounded-lg h-64 overflow-y-auto border border-gray-700 dark:border-gray-600">
                    {terminalOutput.map((line, idx) => (
                      <div key={idx}>{line}</div>
                    ))}
                  </div>
                </div>
              )}

              {activeTab === 'extended' && (
                <div>
                  <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2">
                    Extended Info
                  </h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Extended agent details will appear here (to be implemented).
                  </p>
                </div>
              )}

              {activeTab === 'actions' && (
                <div>
                  <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2">
                    Actions
                  </h3>
                  <button
                    type="button"
                    onClick={handleDummyInstruction}
                    className="w-full bg-violet-600 text-white py-2.5 px-5 rounded-lg hover:bg-violet-700 focus:ring-4 focus:outline-none focus:ring-violet-300 font-medium text-sm text-center dark:bg-violet-500 dark:hover:bg-violet-600 dark:focus:ring-violet-800"
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
