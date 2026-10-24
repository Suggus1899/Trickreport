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
export interface Notification { id: string; title: string; body: string; type?: string; read: boolean; ref_id?: string; ref_type?: string; created_at: string }
export interface CatalogCategory { id: string; tenant_id: string; name: string; description: string; created_at: string; updated_at: string }
export interface CatalogService { id: string; tenant_id: string; category_id?: string | null; category_name?: string; name: string; description: string; requires_approval: boolean; active: boolean; created_at: string; updated_at: string }
export interface Approval { id: string; tenant_id: string; target_type: string; target_id: string; requester_id: string; approver_id?: string | null; status: string; comment: string; decided_by?: string | null; decided_at?: string | null; created_at: string; updated_at: string }
export interface ApprovalEvent { id: string; approval_id: string; actor_id: string; from_status: string; to_status: string; comment: string; created_at: string }
export interface CIType { id: string; tenant_id: string; name: string; description: string; created_at: string; updated_at: string }
export interface Asset { id: string; tenant_id: string; ci_type_id: string; ci_type_name?: string; name: string; serial_number: string; status: string; attributes: Record<string, unknown>; owner_id?: string | null; owner_name?: string; created_at: string; updated_at: string }
export interface CIRelationship { id: string; tenant_id: string; source_ci_id: string; target_ci_id: string; relationship_type: string; created_at: string }
export interface Problem { id: string; tenant_id: string; title: string; description: string; status: string; root_cause: string; workaround: string; created_by: string; created_at: string; updated_at: string }
export interface Change { id: string; tenant_id: string; title: string; description: string; problem_id?: string | null; status: string; risk: string; planned_start?: string | null; planned_end?: string | null; implemented_at?: string | null; closed_at?: string | null; created_by: string; created_at: string; updated_at: string }
export interface ITSMReport { changes_by_status: { status: string; count: number }[]; problems_by_status: { status: string; count: number }[]; assets_by_type: { type: string; count: number }[]; approvals_by_status: { status: string; count: number }[]; pending_approvals: number; active_services: number }
export const problemTransitions: Record<string, string[]> = { open: ['investigating', 'resolved', 'closed'], investigating: ['resolved', 'closed', 'open'], resolved: ['closed', 'investigating'], closed: ['investigating'] };
export const changeTransitions: Record<string, string[]> = { draft: ['planned', 'rejected'], planned: ['awaiting_approval', 'rejected'], awaiting_approval: ['approved', 'rejected', 'planned'], approved: ['in_progress', 'rejected'], in_progress: ['completed', 'rejected'], completed: ['closed'], closed: [], rejected: ['draft'] };
export interface Charts { volume: { date: string; count: number }[]; status_distribution: { status: string; count: number }[]; resolution_time: { priority: string; avg_hours: number }[]; sla_compliance: { total: number; breached: number; compliant: number } }
export const transitions: Record<string, string[]> = { open: ['in_progress', 'waiting_client', 'resolved', 'closed'], in_progress: ['waiting_client', 'resolved', 'closed'], waiting_client: ['in_progress', 'resolved', 'closed'], resolved: ['closed', 'in_progress'], closed: ['in_progress'] };
export const priorities = ['low', 'medium', 'high', 'critical'] as const;
