import { API_URL } from './config';

export { API_URL };

export interface LoginInput {
  email: string;
  password: string;
}

export interface LoginResult {
  token: string;
  user: User;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'agent' | 'end_user';
  tenant_id: string;
  active: boolean;
  created_at: string;
}

export interface Summary {
  total_tickets: number;
  open_tickets: number;
  resolved_tickets: number;
  sla_breached: number;
}

export interface Ticket {
  id: string;
  tenant_id: string;
  title: string;
  description: string;
  status: string;
  priority: string;
  category: string;
  created_by: string;
  creator_name?: string;
  assigned_to?: string;
  assignee_name?: string;
  sla_breached: boolean;
  created_at: string;
  updated_at: string;
}

export interface Article {
  id: string;
  tenant_id: string;
  title: string;
  content: string;
  category: string;
  tags: string[];
  published: boolean;
  created_by: string;
  author_name?: string;
  created_at: string;
  updated_at: string;
}

export interface SLAPolicy {
  id?: string;
  tenant_id?: string;
  priority: string;
  response_time_minutes: number;
  resolution_time_minutes: number;
  escalation_minutes: number;
  created_at?: string;
  updated_at?: string;
}

export interface AutomationRule {
  id?: string;
  tenant_id?: string;
  name: string;
  description: string;
  trigger_type: string;
  conditions: Record<string, any>;
  actions: any[];
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface VolumePoint {
  date: string;
  count: number;
}

export interface StatusDistribution {
  status: string;
  count: number;
}

export interface ResolutionMetrics {
  priority: string;
  avg_hours: number;
}

export interface ApiError {
  error: string;
}

function getHeaders(token?: string): HeadersInit {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    'X-Tenant-ID': 'default',
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  return headers;
}

export async function api<T>(path: string, init: RequestInit & { token?: string } = {}): Promise<T> {
  const { token, ...rest } = init;
  const res = await fetch(`${API_URL}/api/v1${path}`, {
    ...rest,
    headers: getHeaders(token),
    credentials: 'include',
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({})) as ApiError;
    throw new Error(body.error || `Request failed with status ${res.status}`);
  }
  if (res.status === 204) {
    return undefined as T;
  }
  return res.json() as Promise<T>;
}

export async function login(input: LoginInput): Promise<LoginResult> {
  return api<LoginResult>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export async function getMe(token: string): Promise<User> {
  return api<User>('/auth/me', { token });
}

export async function getSummary(token: string): Promise<Summary> {
  return api<Summary>('/dashboard/summary', { token });
}

export async function getTickets(token: string): Promise<Ticket[]> {
  return api<Ticket[]>('/tickets', { token });
}

export async function getTicket(token: string, id: string): Promise<Ticket> {
  return api<Ticket>(`/tickets/${id}`, { token });
}

export async function createTicket(token: string, data: { title: string; description: string; priority: string; category: string }): Promise<Ticket> {
  return api<Ticket>('/tickets', { token, method: 'POST', body: JSON.stringify(data) });
}

export async function getArticles(token: string, query = ''): Promise<Article[]> {
  return api<Article[]>(`/articles?q=${encodeURIComponent(query)}`, { token });
}

export async function getArticle(token: string, id: string): Promise<Article> {
  return api<Article>(`/articles/${id}`, { token });
}

export async function createArticle(token: string, data: Partial<Article>): Promise<Article> {
  return api<Article>('/articles', { token, method: 'POST', body: JSON.stringify(data) });
}

export async function getUsers(token: string): Promise<User[]> {
  return api<User[]>('/admin/users', { token });
}

export async function createUser(token: string, data: { name: string; email: string; role: string; password: string }): Promise<User> {
  return api<User>('/admin/users', { token, method: 'POST', body: JSON.stringify(data) });
}

export async function getSLAPolicies(token: string): Promise<SLAPolicy[]> {
  return api<SLAPolicy[]>('/admin/sla', { token });
}

export async function upsertSLAPolicy(token: string, priority: string, data: Omit<SLAPolicy, 'priority'>): Promise<SLAPolicy> {
  return api<SLAPolicy>(`/admin/sla/${priority}`, { token, method: 'PUT', body: JSON.stringify(data) });
}

export async function getAutomations(token: string): Promise<AutomationRule[]> {
  return api<AutomationRule[]>('/admin/automations', { token });
}

export async function createAutomation(token: string, data: AutomationRule): Promise<AutomationRule> {
  return api<AutomationRule>('/admin/automations', { token, method: 'POST', body: JSON.stringify(data) });
}

export async function getAnalyticsVolume(token: string): Promise<VolumePoint[]> {
  return api<VolumePoint[]>('/admin/analytics/volume', { token });
}

export async function getAnalyticsStatus(token: string): Promise<StatusDistribution[]> {
  return api<StatusDistribution[]>('/admin/analytics/status-distribution', { token });
}

export async function getAnalyticsResolution(token: string): Promise<ResolutionMetrics[]> {
  return api<ResolutionMetrics[]>('/admin/analytics/resolution-time', { token });
}
