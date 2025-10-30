import React, { useState, useMemo, useEffect } from 'react';
import type { AgentTableEntry } from '../models/Tables';
import AgentTableRow from './tables/AgentTableRow';
import InteractAgentSideBar from './modals/InteractAgentModal';
import InventoryControls from './InventoryControls';
import ErrorBanner from '../partials/ErrorBanner';
import { getToken } from '../services/authService';

// Mock data for operations not yet implemented
const MOCK_AGENTS: AgentTableEntry[] = [
  { id: 'agent-001', hostname: 'recon-host', ip: '192.168.1.1', os: 'Linux', protocol: 'HTTP', lastSeen: '2025-03-15T10:00:00Z', isActive: true },
  { id: 'agent-002', hostname: 'persist-host', ip: '192.168.1.2', os: 'Windows', protocol: 'HTTPS', lastSeen: '2025-04-01T14:30:00Z', isActive: true },
  { id: 'agent-003', hostname: 'data-exfil', ip: '192.168.1.3', os: 'Linux', protocol: 'TCP', lastSeen: '2024-11-20T09:15:00Z', isActive: false },
  { id: 'agent-004', hostname: 'command-host', ip: '192.168.1.4', os: 'Windows', protocol: 'HTTP', lastSeen: '2025-04-02T11:05:00Z', isActive: false },
  { id: 'agent-005', hostname: 'network-host', ip: '192.168.1.5', os: 'Linux', protocol: 'UDP', lastSeen: '2025-03-28T16:45:00Z', isActive: true },
];

const AgentInventory: React.FC = () => {
  const [agents, setAgents] = useState<AgentTableEntry[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedAgentIds, setSelectedAgentIds] = useState<Set<string>>(new Set());
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingAgent, setEditingAgent] = useState<AgentTableEntry | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Fetch agents from backend
  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const response = await fetch('/api/v1/frontend/agents', {
          headers: {
            'Authorization': getToken() || '',
            'Content-Type': 'application/json',
          },
        });
        if (!response.ok) {
          throw new Error(`Failed to fetch agents: ${response.statusText}`);
        }
        const apiData = await response.json();
        
        const data: AgentTableEntry[] = apiData.map((agent: any) => ({
          id: agent.id,
          hostname: agent.hostname,
          ip: agent.ip,
          os: agent.os,
          protocol: agent.protocol,
          lastSeen: agent.last_seen,
          isActive: agent.is_active,
        }));
        
        setAgents(data);
        setError(null);
      } catch (err) {
        setError('Error fetching agents. Please try again later.');
        console.error(err);
        setAgents(MOCK_AGENTS);
      }
    };
    fetchAgents();
  }, []);

  
  const filteredAgents = useMemo(() => {
    if (!searchTerm) return agents;
    return agents.filter(agent =>
      agent.hostname.toLowerCase().includes(searchTerm.toLowerCase()) ||
      agent.ip.toLowerCase().includes(searchTerm.toLowerCase()) ||
      agent.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
      agent.os.toLowerCase().includes(searchTerm.toLowerCase()) ||
      agent.protocol.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [agents, searchTerm]);

  const handleSelectRow = (id: string) => {
    setSelectedAgentIds(prev => {
      const newSelection = new Set(prev);
      if (newSelection.has(id)) {
        newSelection.delete(id);
      } else {
        newSelection.add(id);
      }
      return newSelection;
    });
  };

  const handleSelectAll = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.checked) {
      const allIds = new Set(filteredAgents.map(agent => agent.id));
      setSelectedAgentIds(allIds);
    } else {
      setSelectedAgentIds(new Set());
    }
  };

  const handleEditClick = (agent: AgentTableEntry) => {
    setEditingAgent(agent);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingAgent(null);
  };

  // const handleSaveChanges = (updatedAgent: AgentTableEntry) => {
  //   // Mock implementation
  //   console.log('Saving:', updatedAgent);
  //   setAgents(prev =>
  //     prev.map(agent => (agent.id === updatedAgent.id ? { ...agent, ...updatedAgent, lastSeen: new Date().toISOString() } : agent))
  //   );
  //   handleCloseModal();
  // };

  const handleDeleteSelected = () => {
    // Mock implementation
    console.log('Deleting:', Array.from(selectedAgentIds));
    setAgents(prev => prev.filter(agent => !selectedAgentIds.has(agent.id)));
    setSelectedAgentIds(new Set());
  };

  const updateSelectedStatus = (isActive: boolean) => {
    // Mock implementation
    console.log(`Updating status to ${isActive} for:`, Array.from(selectedAgentIds));
    setAgents(prev =>
      prev.map(agent =>
        selectedAgentIds.has(agent.id) ? { ...agent, isActive, lastSeen: new Date().toISOString() } : agent
      )
    );
    setSelectedAgentIds(new Set());
  };

  const handleActivateSelected = () => updateSelectedStatus(true);
  const handleDeactivateSelected = () => updateSelectedStatus(false);

  const isAllSelected = filteredAgents.length > 0 && selectedAgentIds.size === filteredAgents.length;

  return (
    <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
      {error && (
        <ErrorBanner
          message={error}
          onDismiss={() => setError(null)}
        />
      )}
      <InventoryControls
        searchTerm={searchTerm}
        onSearchChange={setSearchTerm}
        selectedCount={selectedAgentIds.size}
        onDeleteSelected={handleDeleteSelected}
        onActivateSelected={handleActivateSelected}
        onDeactivateSelected={handleDeactivateSelected}
      />
      <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
        <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
          <tr>
            <th scope="col" className="p-4">
              <div className="flex items-center">
                <input
                  id="checkbox-all-search"
                  type="checkbox"
                  checked={isAllSelected}
                  onChange={handleSelectAll}
                  disabled={filteredAgents.length === 0}
                  className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                />
                <label htmlFor="checkbox-all-search" className="sr-only">checkbox</label>
              </div>
            </th>
            <th scope="col" className="px-6 py-3">Hostname / Details</th>
            <th scope="col" className="px-6 py-3">IP</th>
            <th scope="col" className="px-6 py-3">OS</th>
            <th scope="col" className="px-6 py-3">Protocol</th>
            <th scope="col" className="px-6 py-3">Last Seen</th>
            <th scope="col" className="px-6 py-3">Status</th>
            <th scope="col" className="px-6 py-3">Action</th>
          </tr>
        </thead>
        <tbody>
          {filteredAgents.map(agent => (
            <AgentTableRow
              key={agent.id}
              agent={agent}
              isSelected={selectedAgentIds.has(agent.id)}
              onSelect={handleSelectRow}
              onEdit={handleEditClick}
            />
          ))}
          {filteredAgents.length === 0 && (
            <tr>
              <td colSpan={8} className="px-6 py-4 text-center text-gray-500 dark:text-gray-400">
                No agents found.
              </td>
            </tr>
          )}
        </tbody>
      </table>

      <InteractAgentSideBar
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        agent={editingAgent}
        // onSave={handleSaveChanges}
      />
    </div>
  );
};

export default AgentInventory;
