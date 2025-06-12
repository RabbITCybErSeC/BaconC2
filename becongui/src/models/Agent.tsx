export interface Agent {
  id: string;
  hostname: string;
  ip: string;
  lastSeen: string;
  os: string;
  isActive: boolean;
  protocol: string;
}
