import React from 'react';
import { useAuth } from '../hooks/useAuth';

export default function Header() {
  const { user } = useAuth();
  
  return (
    <header className="bg-white border-b border-gray-200 px-6 py-3 flex items-center justify-between">
      <div>
        <h1 className="text-xl font-semibold text-gray-800">OpAMP Control Pane</h1>
        <p className="text-sm text-gray-500">Collector Management Dashboard</p>
      </div>
      <div className="flex items-center space-x-4">
        {user && (
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
              {user.name.charAt(0).toUpperCase()}
            </div>
            <div className="text-sm">
              <p className="font-medium text-gray-700">{user.name}</p>
              <p className="text-xs text-gray-400 capitalize">{user.roles.join(', ')}</p>
            </div>
          </div>
        )}
      </div>
    </header>
  );
}