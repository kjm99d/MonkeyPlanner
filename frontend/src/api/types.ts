export type IssueStatus = 'Pending' | 'Approved' | 'InProgress' | 'Done' | 'Rejected';

export type Criterion = { text: string; done: boolean };

export interface Issue {
  id: string;
  boardId: string;
  parentId?: string | null;
  title: string;
  body: string;
  instructions: string;
  status: IssueStatus;
  properties: Record<string, unknown>;
  criteria: Criterion[];
  position: number;
  createdAt: string;
  updatedAt: string;
  approvedAt?: string | null;
  completedAt?: string | null;
  blockedBy?: string[];
}

export interface Comment {
  id: string;
  issueId: string;
  body: string;
  createdAt: string;
}

export type PropertyType = 'text' | 'number' | 'select' | 'multi_select' | 'date' | 'checkbox';

export interface BoardProperty {
  id: string;
  boardId: string;
  name: string;
  type: PropertyType;
  options: string[];
  position: number;
  createdAt: string;
}

export interface Board {
  id: string;
  name: string;
  viewType: 'kanban' | 'list';
  createdAt: string;
}

export type WebhookEvent = 'issue.created' | 'issue.approved' | 'issue.status_changed' | 'issue.deleted';

export interface Webhook {
  id: string;
  boardId: string;
  name: string;
  url: string;
  events: WebhookEvent[];
  enabled: boolean;
  createdAt: string;
}

export interface DayCount {
  date: string;
  created: number;
  approved: number;
  completed: number;
}

export interface DayStats {
  created: Issue[];
  approved: Issue[];
  completed: Issue[];
}
