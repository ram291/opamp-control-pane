import React from 'react';
import { useParams } from 'react-router-dom';

export default function AgentDetail() {
  const { id } = useParams();

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-gray-800">Agent Detail</h2>
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <p className="text-gray-500">Agent ID: <span className="font-mono text-gray-800">{id}</span></p>
        <p className="mt-4 text-gray-500">Detailed agent information and management controls will be displayed here.</p>
      </div>
    </div>
  );
}