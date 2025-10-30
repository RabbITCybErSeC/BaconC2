export interface Agent {
  id: string;
  hostname: string;
  ip: string;
  os: string;
  protocol: string;
  isActive: boolean;
  lastSeen: string;
  // extended_info: string;
}

export interface AgentSession {
  id: number;
  agent_id: string;
  session_id: string;
  start_time: string;
  end_time?: string | null;
  ip_address: string;
  user_agent: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface NetworkInterface {
  name: string;
  mac_address?: string;
  ip_addresses?: string[];
  netmask?: string;
  gateway?: string;
}

export interface ExtendedAgentInfo {
  agent_id: string;
  network_interfaces: NetworkInterface[];
  architecture: string;
  cpu_info: string;
  memory_total: number;
  memory_free: number;
  disk_total: number;
  disk_free: number;
  uptime: number;
  process_count: number;
  username: string;
  domain: string;
  last_boot_time: string;
}

