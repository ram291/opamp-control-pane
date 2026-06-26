const API_BASE = '/api';

export interface AgentDTO {
  id: string;
  version: string;
  status: string;
  uptime: string;
  hostname: string;
}

export interface AvailableVersion {
  version: string;
  date: string;
  sha256: string;
}

export interface UserDTO {
  id: string;
  name: string;
  email: string;
  roles: string[];
  permissions: string[];
}

export interface AgentsResponse {
  agents: AgentDTO[];
}

export interface VersionsResponse {
  versions: AvailableVersion[];
}

export interface ConfigResponse {
  config: {
    content: string;
    source: string;
  };
}

export interface UpgradeRequest {
  version: string;
}

export async function getCurrentUser(): Promise<UserDTO> {
  const resp = await fetch(`${API_BASE}/me`);
  if (!resp.ok) throw new Error(`API error: ${resp.status}`);
  return resp.json();
}

export async function getAgents(): Promise<AgentDTO[]> {
  const resp = await fetch(`${API_BASE}/agents`);
  if (!resp.ok) throw new Error(`API error: ${resp.status}`);
  const data: AgentsResponse = await resp.json();
  return data.agents;
}

export async function getAgent(id: string): Promise<AgentDTO> {
  const resp = await fetch(`${API_BASE}/agents/${id}`);
  if (!resp.ok) throw new Error(`API error: ${resp.status}`);
  return resp.json();
}

export async function getVersions(): Promise<AvailableVersion[]> {
  const resp = await fetch(`${API_BASE}/versions`);
  if (!resp.ok) throw new Error(`API error: ${resp.status}`);
  const data: VersionsResponse = await resp.json();
  return data.versions;
}

export async function getAgentConfig(id: string): Promise<ConfigResponse> {
  const resp = await fetch(`${API_BASE}/agents/${id}/config`);
  if (!resp.ok) throw new Error(`API error: ${resp.status}`);
  return resp.json();
}

export async function upgradeAgent(id: string, version: string): Promise<void> {
  const resp = await fetch(`${API_BASE}/agents/${id}/upgrade`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ version }),
  });
  if (!resp.ok) throw new Error(`Upgrade failed: ${resp.status}`);
}