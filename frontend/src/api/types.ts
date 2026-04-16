export type IssueStatus = 'Pending' | 'Approved' | 'InProgress' | 'Done';

export interface Issue {
  id: string;
  boardId: string;
  parentId?: string | null;
  title: string;
  body: string;
  status: IssueStatus;
  createdAt: string;
  updatedAt: string;
  approvedAt?: string | null;
  completedAt?: string | null;
}

export interface Board {
  id: string;
  name: string;
  viewType: 'kanban' | 'list';
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
