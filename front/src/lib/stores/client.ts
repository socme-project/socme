export type Client = {
  ID: string;
  Name: string;
  Logo: string;

  Host: string;

  SshPort: string;
  SshUsername: string;
  SshPassword: string;

  WazuhPort: string;
  WazuhUsername: string;
  WazuhPassword: string;

  IndexerPort: string;
  IndexerUsername: string;
  IndexerPassword: string;

  LastAlert: string;
  WazuhIsAlive: boolean;
  WazuhVersion: string;
  ConnectedAgents?: number;
  DisconnectedAgents?: number;

  Os?: string;
  MachineHost?: string;
  Kernel?: string;
  CPU?: string;
  GPU?: string;

  IP?: string;
  Uptime?: string;

  Disk?: string;
  Memory?: string;
  Swap?: string;
  CPUUsage?: string;
};
