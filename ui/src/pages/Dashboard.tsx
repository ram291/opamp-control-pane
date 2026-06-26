import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getAgents, getVersions, AgentDTO, AvailableVersion } from '../api/client';
import { useAuth } from '../hooks/useAuth';

export default function Dashboard() {
  const [agents, setAgents] = useState<AgentDTO[]>([]);
  const [versions, setVersions] = useState<AvailableVersion[]>([]);
  const [loading, setLoading] = useState(true);
  const { hasPermission } = useAuth();

  useEffect(() => {
    Promise.all([
      getAgents(),
      hasPermission('upgrade:view') ? getVersions() : Promise.resolve([]),
    ])
      .then(([agentsData, versionsData]) => {
        setAgents(agentsData);
        setVersions(versionsData);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [hasPermission]);

  if (loading) {
    return <div className="text-center py-12 text-gray-500">Loading dashboard...</div>;
  }

  const healthyCount = agents.filter(a => a.status === 'running').length;
  const totalCount = agents.length;

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-gray-800">Dashboard</h2>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="text-sm text-gray-500 mb-1">Total Agents</div>
          <div className="text-3xl font-bold text-gray-800">{totalCount}</div>
        </div>
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="text-sm text-gray-500 mb-1">Healthy</div>
          <div className="text-3xl font-bold text-green-600">{healthyCount}</div>
        </div>
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="text-sm text-gray-500 mb-1">Available Versions</div>
          <div className="text-3xl font-bold text-blue-600">{versions.length}</div>
        </div>
      </div>

      {/* Agents Table */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-semibold text-gray-800">Managed Collectors</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Version</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Hostname</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {agents.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-8 text-center text-gray-500">
                    No agents connected yet
                  </td>
                </tr>
              ) : (
                agents.map(agent => (
                  <tr key={agent.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm font-mono text-gray-900">{agent.id.slice(0, 8)}...</td>
                    <td className="px-6 py-4 text-sm text-gray-600">{agent.version}</td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        agent.status === 'running' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                      }`}>
                        {agent.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">{agent.hostname}</td>
                    <td className="px-6 py-4 text-sm">
                      <Link to={`/agents/${agent.id}`} className="text-blue-600 hover:text-blue-800">
                        View Details →
                      </Link>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}