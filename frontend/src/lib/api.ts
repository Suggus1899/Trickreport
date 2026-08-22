import { API_URL } from './config';

export { API_URL };

export interface LoginInput {
  email: string;
  password: string;
}

export interface LoginResult {
  token?: string;
  refresh_token?: string;
  user?: User;
  requires_mfa?: boolean;
  mfa_token?: string;
}

export interface MFASetupResult {
  secret: string;
  qr_url: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'agent' | 'end_user';
  tenant_id: string;
  active: boolean;
  created_at: string;
  mfa_enabled?: boolean;
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

export interface Comment {
  id: string;
  ticket_id: string;
  author_id: string;
  author_name?: string;
  content: string;
  is_internal: boolean;
  created_at: string;
}

export interface HistoryEntry {
  id: string;
  ticket_id: string;
  field: string;
  old_value: string;
  new_value: string;
  changed_by: string;
  changed_by_name?: string;
  created_at: string;
}

export interface Notification {
  id: string;
  tenant_id: string;
  user_id: string;
  title: string;
  body: string;
  type: string;
  ref_id?: string;
  ref_type?: string;
  read: boolean;
  created_at: string;
}

export interface Attachment {
  id: string;
  ticket_id: string;
  filename: string;
  size: number;
  content_type: string;
  uploaded_by: string;
  uploaded_by_name?: string;
  created_at: string;
}

/* ─── Typed API error classes ───────────────────────────────────────── */

export class ApiAuthError extends Error {
  status = 401;
  constructor(message: string) {
    super(message);
    this.name = 'ApiAuthError';
  }
}

export class ApiValidationError extends Error {
  status: number;
  details?: unknown;
  constructor(message: string, status = 422, details?: unknown) {
    super(message);
    this.name = 'ApiValidationError';
    this.status = status;
    this.details = details;
  }
}

export class ApiNetworkError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ApiNetworkError';
  }
}

/* ─── Token refresh support ─────────────────────────────────────────── */

let refreshTokenFn: (() => Promise<string | null>) | null = null;

export function setRefreshTokenStrategy(fn: (() => Promise<string | null>) | null) {
  refreshTokenFn = fn;
}

/* ─── Core request with timeout, retry, and token refresh ───────────── */

const DEFAULT_TIMEOUT = 30000;
const MAX_RETRIES = 3;

const CSRF_COOKIE_NAME = 'trickreport_csrf';
const CSRF_HEADER_NAME = 'X-CSRF-Token';
const SAFE_METHODS = new Set(['GET', 'HEAD', 'OPTIONS']);

/** Reads the double-submit CSRF cookie. Only present in a browser context. */
function getCsrfToken(): string | undefined {
  if (typeof document === 'undefined') return undefined;
  const match = document.cookie.match(new RegExp(`(?:^|; )${CSRF_COOKIE_NAME}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : undefined;
}

function getHeaders(token: string | undefined, skipJson: boolean, method: string): HeadersInit {
  const headers: HeadersInit = {
    'X-Tenant-ID': 'default',
  };
  if (!skipJson) {
    headers['Content-Type'] = 'application/json';
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  if (!SAFE_METHODS.has(method)) {
    const csrfToken = getCsrfToken();
    if (csrfToken) {
      headers[CSRF_HEADER_NAME] = csrfToken;
    }
  }
  return headers;
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function doFetch<T>(
  path: string,
  init: RequestInit & { token?: string; signal?: AbortSignal },
  attempt = 0,
): Promise<T> {
  const { token, signal, ...rest } = init;
  const method = (rest.method || 'GET').toUpperCase();
  const skipJson = rest.body instanceof FormData;
  const headers = getHeaders(token, skipJson, method);

  // Combine caller signal with a timeout signal.
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), DEFAULT_TIMEOUT);
  if (signal) {
    if (signal.aborted) controller.abort();
    else signal.addEventListener('abort', () => controller.abort(), { once: true });
  }

  let res: Response;
  try {
    res = await fetch(`${API_URL}/api/v1${path}`, {
      ...rest,
      headers,
      credentials: 'include',
      signal: controller.signal,
    });
  } catch (err) {
    clearTimeout(timeoutId);
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw new ApiNetworkError('Request timed out');
    }
    throw new ApiNetworkError(err instanceof Error ? err.message : 'Network request failed');
  }
  clearTimeout(timeoutId);

  // 401: attempt token refresh once, then retry.
  if (res.status === 401 && token && refreshTokenFn && attempt === 0) {
    const newToken = await refreshTokenFn();
    if (newToken) {
      return doFetch<T>(path, { ...init, token: newToken }, 1);
    }
    const body = await res.json().catch(() => ({})) as ApiError;
    throw new ApiAuthError(body.error || 'Authentication required');
  }

  if (res.status === 401) {
    const body = await res.json().catch(() => ({})) as ApiError;
    throw new ApiAuthError(body.error || 'Authentication required');
  }

  if (res.status >= 400 && res.status < 500 && res.status !== 429) {
    const body = await res.json().catch(() => ({})) as ApiError;
    throw new ApiValidationError(body.error || `Request failed with status ${res.status}`, res.status, body);
  }

  if (res.status >= 500 && attempt < MAX_RETRIES) {
    const backoff = Math.pow(2, attempt) * 500;
    await sleep(backoff);
    return doFetch<T>(path, init, attempt + 1);
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({})) as ApiError;
    throw new Error(body.error || `Request failed with status ${res.status}`);
  }

  if (res.status === 204) {
    return undefined as T;
  }
  return res.json() as Promise<T>;
}

export async function api<T>(
  path: string,
  init: RequestInit & { token?: string; signal?: AbortSignal } = {},
): Promise<T> {
  return doFetch<T>(path, init);
}

/* ─── Auth ──────────────────────────────────────────────────────────── */

export async function login(input: LoginInput): Promise<LoginResult> {
  return api<LoginResult>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export async function register(input: {
  name: string;
  email: string;
  password: string;
}): Promise<LoginResult> {
  return api<LoginResult>('/auth/register', {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export async function getMe(token?: string): Promise<User> {
  return api<User>('/auth/me', { token });
}

export async function logout(): Promise<void> {
  return api<void>('/auth/logout', { method: 'POST' });
}

/* ─── MFA ───────────────────────────────────────────────────────────── */

export async function mfaLogin(mfaToken: string, code: string): Promise<LoginResult> {
  return api<LoginResult>('/auth/mfa/login', {
    method: 'POST',
    body: JSON.stringify({ mfa_token: mfaToken, code }),
  });
}

// MFA setup/verify/enable/disable are only ever called from client-side
// islands (profile settings) — cookie auth via credentials:'include' is
// enough, no token param needed.

export async function mfaSetup(): Promise<MFASetupResult> {
  return api<MFASetupResult>('/auth/mfa/setup', { method: 'POST' });
}

export async function mfaVerify(code: string): Promise<{ valid: boolean }> {
  return api<{ valid: boolean }>('/auth/mfa/verify', {
    method: 'POST',
    body: JSON.stringify({ code }),
  });
}

export async function mfaEnable(secret: string, code: string): Promise<void> {
  return api<void>('/auth/mfa/enable', {
    method: 'POST',
    body: JSON.stringify({ secret, code }),
  });
}

export async function mfaDisable(code: string): Promise<void> {
  return api<void>('/auth/mfa/disable', {
    method: 'POST',
    body: JSON.stringify({ code }),
  });
}

/* ─── Dashboard ─────────────────────────────────────────────────────── */

export async function getSummary(token?: string): Promise<Summary> {
  return api<Summary>('/dashboard/summary', { token });
}

/* ─── Tickets ───────────────────────────────────────────────────────── */

export async function getTickets(token?: string): Promise<Ticket[]> {
  return api<Ticket[]>('/tickets', { token });
}

export async function getTicket(id: string, token?: string): Promise<Ticket> {
  return api<Ticket>(`/tickets/${id}`, { token });
}

export async function createTicket(data: { title: string; description: string; priority: string; category: string }, token?: string): Promise<Ticket> {
  return api<Ticket>('/tickets', { token, method: 'POST', body: JSON.stringify(data) });
}

export async function updateTicketStatus(id: string, status: string, token?: string): Promise<{ status: string }> {
  return api<{ status: string }>(`/tickets/${id}/status`, { token, method: 'PATCH', body: JSON.stringify({ status }) });
}

export async function assignTicket(id: string, assignedTo: string, token?: string): Promise<{ assigned_to: string | null }> {
  return api<{ assigned_to: string | null }>(`/tickets/${id}/assign`, { token, method: 'POST', body: JSON.stringify({ assigned_to: assignedTo }) });
}

export async function getTicketComments(id: string, token?: string): Promise<Comment[]> {
  return api<Comment[]>(`/tickets/${id}/comments`, { token });
}

export async function addTicketComment(id: string, content: string, isInternal = false, token?: string): Promise<Comment> {
  return api<Comment>(`/tickets/${id}/comments`, { token, method: 'POST', body: JSON.stringify({ content, is_internal: isInternal }) });
}

export async function getTicketHistory(id: string, token?: string): Promise<HistoryEntry[]> {
  return api<HistoryEntry[]>(`/tickets/${id}/history`, { token });
}

export async function getTicketAttachments(id: string, token?: string): Promise<Attachment[]> {
  return api<Attachment[]>(`/tickets/${id}/attachments`, { token });
}

export async function uploadTicketAttachment(id: string, file: File, token?: string): Promise<Attachment> {
  const formData = new FormData();
  formData.append('file', file);
  return api<Attachment>(`/tickets/${id}/attachments`, { token, method: 'POST', body: formData });
}

/* ─── Articles ──────────────────────────────────────────────────────── */

export async function getArticles(query = '', token?: string): Promise<Article[]> {
  return api<Article[]>(`/articles?q=${encodeURIComponent(query)}`, { token });
}

export async function getArticle(id: string, token?: string): Promise<Article> {
  return api<Article>(`/articles/${id}`, { token });
}

export async function createArticle(data: Partial<Article>, token?: string): Promise<Article> {
  return api<Article>('/articles', { token, method: 'POST', body: JSON.stringify(data) });
}

/* ─── Admin ─────────────────────────────────────────────────────────── */

export async function getUsers(token?: string): Promise<User[]> {
  return api<User[]>('/admin/users', { token });
}

export async function createUser(data: { name: string; email: string; role: string; password: string }, token?: string): Promise<User> {
  return api<User>('/admin/users', { token, method: 'POST', body: JSON.stringify(data) });
}

export async function updateUser(id: string, data: Partial<Pick<User, 'name' | 'role' | 'active'>>, token?: string): Promise<User> {
  return api<User>(`/admin/users/${id}`, { token, method: 'PUT', body: JSON.stringify(data) });
}

export async function getSLAPolicies(token?: string): Promise<SLAPolicy[]> {
  return api<SLAPolicy[]>('/admin/sla', { token });
}

export async function upsertSLAPolicy(priority: string, data: Omit<SLAPolicy, 'priority'>, token?: string): Promise<SLAPolicy> {
  return api<SLAPolicy>(`/admin/sla/${priority}`, { token, method: 'PUT', body: JSON.stringify(data) });
}

export async function getAutomations(token?: string): Promise<AutomationRule[]> {
  return api<AutomationRule[]>('/admin/automations', { token });
}

export async function createAutomation(data: AutomationRule, token?: string): Promise<AutomationRule> {
  return api<AutomationRule>('/admin/automations', { token, method: 'POST', body: JSON.stringify(data) });
}

export async function updateAutomation(id: string, data: AutomationRule, token?: string): Promise<AutomationRule> {
  return api<AutomationRule>(`/admin/automations/${id}`, { token, method: 'PUT', body: JSON.stringify(data) });
}

export async function deleteAutomation(id: string, token?: string): Promise<void> {
  return api<void>(`/admin/automations/${id}`, { token, method: 'DELETE' });
}

export async function getAnalyticsVolume(token?: string): Promise<VolumePoint[]> {
  return api<VolumePoint[]>('/admin/analytics/volume', { token });
}

export async function getAnalyticsStatus(token?: string): Promise<StatusDistribution[]> {
  return api<StatusDistribution[]>('/admin/analytics/status-distribution', { token });
}

export async function getAnalyticsResolution(token?: string): Promise<ResolutionMetrics[]> {
  return api<ResolutionMetrics[]>('/admin/analytics/resolution-time', { token });
}

/* ─── Notifications ─────────────────────────────────────────────────── */

export async function getNotifications(token?: string): Promise<Notification[]> {
  return api<Notification[]>('/notifications', { token });
}

export async function markNotificationRead(id: string, token?: string): Promise<void> {
  return api<void>(`/notifications/${id}/read`, { token, method: 'POST' });
}

export async function markAllNotificationsRead(token?: string): Promise<void> {
  return api<void>('/notifications/read-all', { token, method: 'POST' });
}
