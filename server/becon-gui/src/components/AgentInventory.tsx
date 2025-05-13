import React, { useState, useMemo, useEffect } from 'react';
import { Agent } from '../models/Agent.tsx';
import AgentTableRow from './tables/AgentTableRow.tsx';
import EditAgentModal from './modals/EditAgentModal.tsx';
import InventoryControls from './InventorControls.tsx';

const MOCK_AGENTS: Agent[] = [
  { id: 'agent-001', name: 'Recon Agent', description: 'Gathers system information.', status: 'Active', lastModified: '2025-03-15T10:00:00Z', createdBy: 'admin' },
  { id: 'agent-002', name: 'Persistence Agent', description: 'Ensures continuous access.', status: 'Active', lastModified: '2025-04-01T14:30:00Z', createdBy: 'c2.ops' },
  { id: 'agent-003', name: 'Data Exfil Agent', description: 'Extracts sensitive data.', status: 'Inactive', lastModified: '2024-11-20T09:15:00Z', createdBy: 'admin' },
  { id: 'agent-004', name: 'Command Agent', description: 'Executes remote commands.', status: 'Draft', lastModified: '2025-04-02T11:05:00Z', createdBy: 'dev.team' },
  { id: 'agent-005', name: 'Network Agent', description: 'Monitors network traffic.', status: 'Active', lastModified: '2025-03-28T16:45:00Z', createdBy: 'net.sec' },
];

const AgentInventory: React.FC = () => {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedAgentIds, setSelectedAgentIds] = useState<Set<string>>(new Set());
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingAgent, setEditingAgent] = useState<Agent | null>(null);

  // Simulate fetching data
  useEffect(() => {
    // In a real app, fetch data here, e.g., using fetch or axios
    setAgents(MOCK_AGENTS);
  }, []);

  const filteredAgents = useMemo(() => {
    if (!searchTerm) return agents;
    return agents.filter(agent =>
      agent.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      agent.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
      agent.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
      agent.createdBy.toLowerCase().includes(searchTerm.toLowerCase())
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

  const handleEditClick = (agent: Agent) => {
    setEditingAgent(agent);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingAgent(null);
  };

  const handleSaveChanges = (updatedAgent: Agent) => {
    // TODO: Implement actual save logic (e.g., API call)
    console.log('Saving:', updatedAgent);
    setAgents(prev =>
      prev.map(agent => agent.id === updatedAgent.id ? { ...agent, ...updatedAgent, lastModified: new Date().toISOString() } : agent) // Update lastModified on save
    );
    handleCloseModal();
  };

  const handleDeleteSelected = () => {
    // TODO: Implement actual delete logic (e.g., API call)
    console.log('Deleting:', Array.from(selectedAgentIds));
    setAgents(prev => prev.filter(agent => !selectedAgentIds.has(agent.id)));
    setSelectedAgentIds(new Set()); // Clear selection
  };

  const updateSelectedStatus = (status: Agent['status']) => {
    // TODO: Implement actual status update logic (e.g., API call)
    console.log(`Updating status to ${status} for:`, Array.from(selectedAgentIds));
    setAgents(prev => prev.map(agent =>
      selectedAgentIds.has(agent.id) ? { ...agent, status: status, lastModified: new Date().toISOString() } : agent
    ));
    setSelectedAgentIds(new Set()); // Clear selection
  };

  const handleActivateSelected = () => updateSelectedStatus('Active');
  const handleDeactivateSelected = () => updateSelectedStatus('Inactive');

  const isAllSelected = filteredAgents.length > 0 && selectedAgentIds.size === filteredAgents.length;

  return (
    <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
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
                  className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600" />
                <label htmlFor="checkbox-all-search" className="sr-only">checkbox</label>
              </div>
            </th>
            <th scope="col" className="px-6 py-3">Name / Details</th>
            <th scope="col" className="px-6 py-3">Description</th>
            <th scope="col" className="px-6 py-3">Last Modified</th>
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
              <td colSpan={6} className="px-6 py-4 text-center text-gray-500 dark:text-gray-400">
                No agents found.
              </td>
            </tr>
          )}
        </tbody>
      </table>

      <EditAgentModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        agent={editingAgent}
        onSave={handleSaveChanges}
      />
    </div>
  );
};

export default AgentInventory;
