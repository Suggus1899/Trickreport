import { useMemo, useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { createUser, updateUser, type User } from '@/lib/api';
import { toast } from '@/lib/toast';

const ROLES = [
  { value: 'end_user', label: 'End User' },
  { value: 'agent', label: 'Agent' },
  { value: 'admin', label: 'Admin' },
];

// A role isn't an alarm, so it never gets the destructive/coral treatment —
// that color is reserved for SLA breach and critical priority.
const ROLE_VARIANT: Record<string, 'default' | 'secondary' | 'outline'> = {
  admin: 'default',
  agent: 'secondary',
  end_user: 'outline',
};

export function UsersManager({ initialUsers }: { initialUsers: User[] }) {
  const [users, setUsers] = useState(initialUsers);
  const [createError, setCreateError] = useState('');
  const [creating, setCreating] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [busyId, setBusyId] = useState<string | null>(null);
  const [rowError, setRowError] = useState('');

  const [search, setSearch] = useState('');
  const [roleFilter, setRoleFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState('');

  const editingUser = users.find((u) => u.id === editingId) ?? null;

  const filteredUsers = useMemo(() => {
    const q = search.trim().toLowerCase();
    return users.filter((u) => {
      const matchQ = !q || u.name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q);
      const matchR = !roleFilter || u.role === roleFilter;
      const matchS = !statusFilter || (statusFilter === 'active' ? u.active : !u.active);
      return matchQ && matchR && matchS;
    });
  }, [users, search, roleFilter, statusFilter]);

  async function handleCreate(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const data = new FormData(form);
    setCreateError('');
    setCreating(true);
    try {
      const created = await createUser({
        name: data.get('name')?.toString() || '',
        email: data.get('email')?.toString() || '',
        role: data.get('role')?.toString() || 'end_user',
        password: data.get('password')?.toString() || '',
      });
      setUsers((prev) => [...prev, created]);
      form.reset();
      toast(`User ${created.name} created`, 'success');
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : 'Failed to create user');
    } finally {
      setCreating(false);
    }
  }

  async function handleToggleActive(u: User) {
    if (!confirm(`${u.active ? 'Deactivate' : 'Reactivate'} user "${u.name}"?`)) return;
    setBusyId(u.id);
    setRowError('');
    try {
      const updated = await updateUser(u.id, { active: !u.active });
      setUsers((prev) => prev.map((x) => (x.id === u.id ? updated : x)));
      toast(`${u.name} ${updated.active ? 'reactivated' : 'deactivated'}`, 'success');
    } catch (err) {
      setRowError(err instanceof Error ? err.message : 'Failed to update user');
    } finally {
      setBusyId(null);
    }
  }

  async function handleRoleChange(u: User, role: string) {
    if (role === u.role) return;
    if (!confirm(`Change "${u.name}" role to ${role.replace('_', ' ')}?`)) return;
    setBusyId(u.id);
    setRowError('');
    try {
      const updated = await updateUser(u.id, { role: role as User['role'] });
      setUsers((prev) => prev.map((x) => (x.id === u.id ? updated : x)));
      toast(`${u.name} is now ${updated.role.replace('_', ' ')}`, 'success');
    } catch (err) {
      setRowError(err instanceof Error ? err.message : 'Failed to update role');
    } finally {
      setBusyId(null);
    }
  }

  async function handleEditSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!editingUser) return;
    const data = new FormData(e.currentTarget);
    setBusyId(editingUser.id);
    setRowError('');
    try {
      const updated = await updateUser(editingUser.id, {
        name: data.get('name')?.toString() || editingUser.name,
        role: (data.get('role')?.toString() as User['role']) || editingUser.role,
      });
      setUsers((prev) => prev.map((x) => (x.id === editingUser.id ? updated : x)));
      setEditingId(null);
      toast(`${updated.name} updated`, 'success');
    } catch (err) {
      setRowError(err instanceof Error ? err.message : 'Failed to update user');
    } finally {
      setBusyId(null);
    }
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div className="lg:col-span-1 rounded-xl border bg-card p-5 h-fit">
        <h2 className="font-heading text-lg font-semibold mb-4">Create user</h2>
        {createError && <p className="text-sm text-destructive mb-3">{createError}</p>}
        <form onSubmit={handleCreate} className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <Label htmlFor="name">Name</Label>
            <Input id="name" name="name" required />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" name="email" type="email" required />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="role">Role</Label>
            <select id="role" name="role" className="h-9 rounded-md border bg-background px-3 text-sm">
              {ROLES.map((r) => (
                <option key={r.value} value={r.value}>
                  {r.label}
                </option>
              ))}
            </select>
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="password">Password</Label>
            <Input id="password" name="password" type="password" required />
          </div>
          <div className="flex justify-end pt-2">
            <Button type="submit" disabled={creating}>
              {creating ? 'Creating…' : 'Create User'}
            </Button>
          </div>
        </form>
      </div>

      <div className="lg:col-span-2 rounded-xl border bg-card overflow-hidden">
        <div className="p-5 pb-4 flex items-center gap-3 flex-wrap">
          <div className="flex flex-col gap-1 flex-1 min-w-48">
            <Label htmlFor="filter-search" className="text-xs">Search</Label>
            <Input id="filter-search" placeholder="Name or email…" value={search} onChange={(e) => setSearch(e.target.value)} />
          </div>
          <div className="flex flex-col gap-1">
            <Label htmlFor="filter-role" className="text-xs">Role</Label>
            <select id="filter-role" className="h-9 rounded-md border bg-background px-3 text-sm" value={roleFilter} onChange={(e) => setRoleFilter(e.target.value)}>
              <option value="">All roles</option>
              {ROLES.map((r) => (
                <option key={r.value} value={r.value}>{r.label}</option>
              ))}
            </select>
          </div>
          <div className="flex flex-col gap-1">
            <Label htmlFor="filter-status" className="text-xs">Status</Label>
            <select id="filter-status" className="h-9 rounded-md border bg-background px-3 text-sm" value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
              <option value="">All</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>
        </div>

        {rowError && <p className="px-5 text-sm text-destructive mb-2">{rowError}</p>}

        <table className="w-full text-left">
          <thead className="bg-muted/50 text-xs uppercase text-muted-foreground">
            <tr>
              <th className="px-5 py-3 font-semibold">Name</th>
              <th className="px-5 py-3 font-semibold">Email</th>
              <th className="px-5 py-3 font-semibold">Role</th>
              <th className="px-5 py-3 font-semibold">Status</th>
              <th className="px-5 py-3 font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {filteredUsers.length > 0 ? (
              filteredUsers.map((u) => (
                <tr key={u.id} className="hover:bg-muted/30 transition-colors">
                  <td className="px-5 py-3 font-semibold">{u.name}</td>
                  <td className="px-5 py-3 text-sm text-muted-foreground">{u.email}</td>
                  <td className="px-5 py-3">
                    <Badge variant={ROLE_VARIANT[u.role] || 'default'}>{u.role.replace('_', ' ')}</Badge>
                  </td>
                  <td className="px-5 py-3">
                    <Badge variant={u.active ? 'status-resolved' : 'status-closed'} dot>{u.active ? 'active' : 'inactive'}</Badge>
                  </td>
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-2 flex-wrap">
                      <Button size="sm" variant="outline" disabled={busyId === u.id} onClick={() => setEditingId(u.id)}>
                        Edit
                      </Button>
                      <Button size="sm" variant="outline" disabled={busyId === u.id} onClick={() => handleToggleActive(u)}>
                        {u.active ? 'Deactivate' : 'Reactivate'}
                      </Button>
                      <select
                        className="h-8 rounded-md border bg-background px-2 text-xs"
                        value={u.role}
                        disabled={busyId === u.id}
                        onChange={(e) => handleRoleChange(u, e.target.value)}
                      >
                        {ROLES.map((r) => (
                          <option key={r.value} value={r.value}>{r.label}</option>
                        ))}
                      </select>
                    </div>
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={5} className="px-5 py-12 text-center text-muted-foreground">
                  No users found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {editingUser && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/40 p-4" onClick={() => setEditingId(null)}>
          <div className="w-full max-w-md rounded-xl border bg-card p-6 shadow-lg" onClick={(e) => e.stopPropagation()}>
            <h2 className="font-heading text-lg font-semibold mb-4">Edit user</h2>
            <form onSubmit={handleEditSubmit} className="flex flex-col gap-4">
              <div className="flex flex-col gap-2">
                <Label htmlFor="edit-name">Name</Label>
                <Input id="edit-name" name="name" defaultValue={editingUser.name} required />
              </div>
              <div className="flex flex-col gap-2">
                <Label htmlFor="edit-role">Role</Label>
                <select id="edit-role" name="role" defaultValue={editingUser.role} className="h-9 rounded-md border bg-background px-3 text-sm">
                  {ROLES.map((r) => (
                    <option key={r.value} value={r.value}>{r.label}</option>
                  ))}
                </select>
              </div>
              <div className="flex justify-end gap-3 pt-2">
                <Button type="button" variant="ghost" onClick={() => setEditingId(null)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={busyId === editingUser.id}>
                  {busyId === editingUser.id ? 'Saving…' : 'Save'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
