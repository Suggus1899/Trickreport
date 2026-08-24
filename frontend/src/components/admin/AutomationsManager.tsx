import { useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { createAutomation, updateAutomation, deleteAutomation, type AutomationRule } from '@/lib/api';
import { toast } from '@/lib/toast';

// Matches the backend's validator.oneof list exactly (handler/automation.go).
const TRIGGER_TYPES = [
  { value: 'ticket_created', label: 'Ticket Created' },
  { value: 'status_changed', label: 'Status Changed' },
  { value: 'sla_breach', label: 'SLA Breach' },
  { value: 'priority_changed', label: 'Priority Changed' },
  { value: 'escalation', label: 'SLA Escalation' },
];

const CONDITION_FIELDS = [
  { value: 'status', label: 'Status' },
  { value: 'priority', label: 'Priority' },
  { value: 'category', label: 'Category' },
];

const CONDITION_OPERATORS = [
  { value: 'equals', label: 'equals' },
  { value: 'not_equals', label: 'not equals' },
  { value: 'contains', label: 'contains' },
];

// Matches the backend's action executor exactly (application/automation/engine.go).
// "add_tag" is deliberately not offered here: the tickets table has no tags
// column, so the backend's AddTag is a permanent no-op (automation_executor.go).
const ACTION_TYPES = [
  { value: 'set_status', label: 'Change Status' },
  { value: 'assign_to', label: 'Assign' },
  { value: 'send_email', label: 'Send Email' },
];

const ACTION_FIELDS: Record<string, { key: string; label: string; placeholder: string }[]> = {
  set_status: [{ key: 'status', label: 'Status', placeholder: 'resolved' }],
  assign_to: [{ key: 'user_id', label: 'Assignee (user ID)', placeholder: '00000000-0000-0000-0000-000000000000' }],
  send_email: [
    { key: 'to', label: 'To', placeholder: 'manager@example.com' },
    { key: 'subject', label: 'Subject', placeholder: 'Ticket needs attention' },
  ],
};

interface Clause {
  field: string;
  op: string;
  value: string;
}

interface ActionRow {
  type: string;
  [key: string]: string;
}

const EMPTY_CLAUSE: Clause = { field: 'status', op: 'equals', value: '' };
const EMPTY_ACTION: ActionRow = { type: 'set_status', status: '' };

function conditionsToClauses(conditions: Record<string, unknown>): Clause[] {
  if (Array.isArray((conditions as { conditions?: unknown }).conditions)) {
    return (conditions as { conditions: Clause[] }).conditions;
  }
  if (Object.keys(conditions).length) {
    return Object.entries(conditions).map(([field, value]) => {
      if (value && typeof value === 'object' && !Array.isArray(value)) {
        const v = value as { op?: string; value?: unknown };
        return { field, op: v.op || 'equals', value: String(v.value ?? '') };
      }
      return { field, op: 'equals', value: String(value ?? '') };
    });
  }
  return [{ ...EMPTY_CLAUSE }];
}

function clausesToConditions(clauses: Clause[]): Record<string, unknown> {
  if (clauses.length === 0) return {};
  if (clauses.length === 1 && clauses[0].op === 'equals') {
    return { [clauses[0].field]: clauses[0].value };
  }
  return { conditions: clauses };
}

export function AutomationsManager({ initialAutomations }: { initialAutomations: AutomationRule[] }) {
  const [automations, setAutomations] = useState(initialAutomations);
  const [editId, setEditId] = useState<string | null>(null);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [triggerType, setTriggerType] = useState(TRIGGER_TYPES[0].value);
  const [isActive, setIsActive] = useState(true);
  const [clauses, setClauses] = useState<Clause[]>([{ ...EMPTY_CLAUSE }]);
  const [actions, setActions] = useState<ActionRow[]>([{ ...EMPTY_ACTION }]);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  function resetBuilder() {
    setEditId(null);
    setName('');
    setDescription('');
    setTriggerType(TRIGGER_TYPES[0].value);
    setIsActive(true);
    setClauses([{ ...EMPTY_CLAUSE }]);
    setActions([{ ...EMPTY_ACTION }]);
    setError('');
  }

  function loadRule(rule: AutomationRule, duplicate: boolean) {
    setEditId(duplicate ? null : rule.id ?? null);
    setName(duplicate ? `${rule.name} (copy)` : rule.name);
    setDescription(rule.description || '');
    setTriggerType(rule.trigger_type);
    setIsActive(rule.is_active);
    setClauses(conditionsToClauses(rule.conditions || {}));
    setActions(rule.actions && rule.actions.length ? (rule.actions as ActionRow[]) : [{ ...EMPTY_ACTION }]);
    setError('');
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');
    setSaving(true);
    const payload: AutomationRule = {
      name,
      description,
      trigger_type: triggerType,
      conditions: clausesToConditions(clauses),
      actions,
      is_active: isActive,
    };
    try {
      if (editId) {
        const updated = await updateAutomation(editId, payload);
        setAutomations((prev) => prev.map((a) => (a.id === editId ? updated : a)));
        toast(`Automation "${updated.name}" updated`, 'success');
      } else {
        const created = await createAutomation(payload);
        setAutomations((prev) => [...prev, created]);
        toast(`Automation "${created.name}" created`, 'success');
      }
      resetBuilder();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save automation');
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(rule: AutomationRule) {
    if (!rule.id || !confirm(`Delete automation "${rule.name}"? This cannot be undone.`)) return;
    try {
      await deleteAutomation(rule.id);
      setAutomations((prev) => prev.filter((a) => a.id !== rule.id));
      toast(`Automation "${rule.name}" deleted`, 'success');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete automation');
    }
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div className="lg:col-span-1 rounded-xl border bg-card p-5 h-fit">
        <h2 className="font-heading text-lg font-semibold mb-1">{editId ? 'Edit automation' : 'Create automation'}</h2>
        <p className="text-xs text-muted-foreground mb-4">
          {editId ? 'Changes will update the existing rule.' : 'Build conditions and actions visually.'}
        </p>
        {error && <p className="text-sm text-destructive mb-3">{error}</p>}
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <Label htmlFor="auto-name">Name</Label>
            <Input id="auto-name" value={name} onChange={(e) => setName(e.target.value)} required />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="auto-desc">Description</Label>
            <Input id="auto-desc" value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="auto-trigger">Trigger type</Label>
            <select
              id="auto-trigger"
              className="h-9 rounded-md border bg-background px-3 text-sm"
              value={triggerType}
              onChange={(e) => setTriggerType(e.target.value)}
            >
              {TRIGGER_TYPES.map((t) => (
                <option key={t.value} value={t.value}>{t.label}</option>
              ))}
            </select>
          </div>

          <div className="flex flex-col gap-2">
            <Label>Conditions</Label>
            {clauses.map((c, i) => (
              <div key={i} className="flex items-center gap-2 flex-wrap">
                <select
                  className="h-8 rounded-md border bg-background px-2 text-xs"
                  value={c.field}
                  onChange={(e) => setClauses((prev) => prev.map((x, idx) => (idx === i ? { ...x, field: e.target.value } : x)))}
                >
                  {CONDITION_FIELDS.map((f) => <option key={f.value} value={f.value}>{f.label}</option>)}
                </select>
                <select
                  className="h-8 rounded-md border bg-background px-2 text-xs"
                  value={c.op}
                  onChange={(e) => setClauses((prev) => prev.map((x, idx) => (idx === i ? { ...x, op: e.target.value } : x)))}
                >
                  {CONDITION_OPERATORS.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
                </select>
                <input
                  type="text"
                  placeholder="value"
                  className="h-8 flex-1 min-w-20 rounded-md border bg-background px-2 text-xs"
                  value={c.value}
                  onChange={(e) => setClauses((prev) => prev.map((x, idx) => (idx === i ? { ...x, value: e.target.value } : x)))}
                />
                <button
                  type="button"
                  className="h-8 w-8 rounded-md border text-destructive text-sm"
                  onClick={() => setClauses((prev) => prev.filter((_, idx) => idx !== i))}
                >
                  ×
                </button>
              </div>
            ))}
            <Button type="button" variant="outline" size="sm" className="self-start" onClick={() => setClauses((prev) => [...prev, { ...EMPTY_CLAUSE }])}>
              + Add condition
            </Button>
          </div>

          <div className="flex flex-col gap-2">
            <Label>Actions</Label>
            {actions.map((a, i) => (
              <div key={i} className="flex flex-col gap-2 p-3 rounded-lg bg-muted/50">
                <div className="flex items-center gap-2">
                  <select
                    className="h-8 flex-1 rounded-md border bg-background px-2 text-xs"
                    value={a.type}
                    onChange={(e) =>
                      setActions((prev) => prev.map((x, idx) => (idx === i ? { type: e.target.value } : x)))
                    }
                  >
                    {ACTION_TYPES.map((t) => <option key={t.value} value={t.value}>{t.label}</option>)}
                  </select>
                  <button
                    type="button"
                    className="h-8 w-8 rounded-md border text-destructive text-sm"
                    onClick={() => setActions((prev) => prev.filter((_, idx) => idx !== i))}
                  >
                    ×
                  </button>
                </div>
                {(ACTION_FIELDS[a.type] || []).map((f) => (
                  <div key={f.key} className="flex flex-col gap-1">
                    <label className="text-xs font-semibold text-muted-foreground ml-1">{f.label}</label>
                    <input
                      type="text"
                      placeholder={f.placeholder}
                      className="h-8 rounded-md border bg-background px-2 text-xs"
                      value={a[f.key] || ''}
                      onChange={(e) =>
                        setActions((prev) => prev.map((x, idx) => (idx === i ? { ...x, [f.key]: e.target.value } : x)))
                      }
                    />
                  </div>
                ))}
              </div>
            ))}
            <Button type="button" variant="outline" size="sm" className="self-start" onClick={() => setActions((prev) => [...prev, { ...EMPTY_ACTION }])}>
              + Add action
            </Button>
          </div>

          <label className="flex items-center gap-2 text-sm font-semibold">
            <input type="checkbox" className="w-auto" checked={isActive} onChange={(e) => setIsActive(e.target.checked)} />
            Active
          </label>

          <div className="flex flex-wrap gap-3 pt-2">
            <Button type="submit" disabled={saving}>
              {saving ? 'Saving…' : editId ? 'Update' : 'Create'}
            </Button>
            <Button type="button" variant="ghost" onClick={resetBuilder}>
              Reset
            </Button>
          </div>
        </form>
      </div>

      <div className="lg:col-span-2 rounded-xl border bg-card overflow-hidden">
        <table className="w-full text-left">
          <thead className="bg-muted/50 text-xs uppercase text-muted-foreground">
            <tr>
              <th className="px-5 py-3 font-semibold">Name</th>
              <th className="px-5 py-3 font-semibold">Trigger</th>
              <th className="px-5 py-3 font-semibold">Status</th>
              <th className="px-5 py-3 font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {automations.length > 0 ? (
              automations.map((rule) => (
                <tr key={rule.id} className="hover:bg-muted/30 transition-colors">
                  <td className="px-5 py-3">
                    <p className="font-semibold">{rule.name}</p>
                    {rule.description && <p className="text-xs text-muted-foreground">{rule.description}</p>}
                  </td>
                  <td className="px-5 py-3 text-sm text-muted-foreground">{rule.trigger_type.replace('_', ' ')}</td>
                  <td className="px-5 py-3">
                    <Badge variant={rule.is_active ? 'status-resolved' : 'status-closed'} dot>{rule.is_active ? 'active' : 'inactive'}</Badge>
                  </td>
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-2 flex-wrap">
                      <Button size="sm" variant="outline" onClick={() => loadRule(rule, false)}>Edit</Button>
                      <Button size="sm" variant="outline" onClick={() => loadRule(rule, true)}>Duplicate</Button>
                      <Button size="sm" variant="destructive" onClick={() => handleDelete(rule)}>Delete</Button>
                    </div>
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={4} className="px-5 py-12 text-center text-muted-foreground">
                  No automations configured yet.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
