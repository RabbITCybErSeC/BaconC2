export interface Agent {
  id: string;
  name: string;
  description: string;
  status: 'Active' | 'Inactive' | 'Draft';
  lastModified: string; // ISO 8601 format
  createdBy: string;
}
