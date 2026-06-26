import React from 'react';
import { NavLink } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const navItems = [
  { path: '/', label: 'Dashboard', icon: '📊', permission: 'agent:list' },
  { path: '/upgrades', label: 'Upgrades', icon: '⬆️', permission: 'upgrade:view' },
  { path: '/config', label: 'Configuration', icon: '⚙️', permission: 'config:view' },
  { path: '/settings', label: 'Settings', icon: '🔧', permission: 'admin:settings' },
];

export default function Sidebar() {
  const { hasPermission } = useAuth();

  const visibleItems = navItems.filter(item => hasPermission(item.permission));

  return (
    <aside className="w-64 bg-white border-r border-gray-200 flex flex-col">
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center space-x-2">
          <span className="text-2xl">🛡️</span>
          <span className="font-bold text-lg text-gray-800">OpAMP Control</span>
        </div>
      </div>
      <nav className="flex-1 p-4 space-y-1">
        {visibleItems.map(item => (
          <NavLink
            key={item.path}
            to={item.path}
            end={item.path === '/'}
            className={({ isActive }) =>
              `flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                isActive
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
              }`
            }
          >
            <span>{item.icon}</span>
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>
      <div className="p-4 border-t border-gray-200">
        <p className="text-xs text-gray-400">OpAMP Control Pane v0.1.0</p>
      </div>
    </aside>
  );
}