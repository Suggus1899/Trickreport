export type Product = 'it' | 'msp' | 'customer-service';
export type Locale = 'es' | 'en';
export interface User { id: string; tenant_id: string; name: string; email: string; role: 'admin' | 'agent' | 'end_user'; active?: boolean; mfa_enabled?: boolean }
export interface Ticket { id: string; title: string; description: string; status: string; priority: string; category: string; created_by: string; creator_name?: string; assigned_to?: string | null; assignee_name?: string; created_at: string; updated_at: string; sla_breached: boolean; version?: number }
export interface Article { id: string; title: string; content: string; category: string; tags: string[]; published: boolean; author_name?: string; updated_at: string }
export interface Comment { id: string; content: string; user_name?: string; is_internal: boolean; created_at: string }
export interface Attachment { id: string; filename: string; file_size: number; content_type: string }
export interface History { id: string; field: string; old_value?: string; new_value?: string; user_name?: string; created_at: string }
export interface Summary { total_tickets: number; open_tickets: number; resolved_tickets: number; sla_breached: number }
export interface SLA { priority: string; response_time_minutes: number; resolution_time_minutes: number; escalation_minutes: number }
export interface Automation { id?: string; name: string; description: string; trigger_type: string; conditions: Record<string, unknown>; actions: unknown[]; is_active: boolean }
export interface Notification { id: string; title: string; body: string; read: boolean; ref_id?: string; ref_type?: string; created_at: string }
export interface Charts { volume: { date: string; count: number }[]; status_distribution: { status: string; count: number }[]; resolution_time: { priority: string; avg_hours: number }[]; sla_compliance: { total: number; breached: number; compliant: number } }
export const transitions: Record<string, string[]> = { open: ['in_progress', 'waiting_client', 'resolved', 'closed'], in_progress: ['waiting_client', 'resolved', 'closed'], waiting_client: ['in_progress', 'resolved', 'closed'], resolved: ['closed', 'in_progress'], closed: ['in_progress'] };
export const priorities = ['low', 'medium', 'high', 'critical'] as const;
