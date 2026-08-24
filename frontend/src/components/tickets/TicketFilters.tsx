import { useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

const STATUS_OPTIONS = [
  { value: 'open', label: 'Open' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'waiting_client', label: 'Waiting Client' },
  { value: 'resolved', label: 'Resolved' },
  { value: 'closed', label: 'Closed' },
];
const PRIORITY_OPTIONS = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
  { value: 'critical', label: 'Critical' },
];

interface Props {
  initialSearch: string;
  initialStatus: string;
  initialPriority: string;
  hasFilters: boolean;
}

export function TicketFilters({ initialSearch, initialStatus, initialPriority, hasFilters }: Props) {
  const [search, setSearch] = useState(initialSearch);
  const [status, setStatus] = useState(initialStatus);
  const [priority, setPriority] = useState(initialPriority);

  function apply(e: FormEvent) {
    e.preventDefault();
    const params = new URLSearchParams();
    if (search) params.set('q', search);
    if (status) params.set('status', status);
    if (priority) params.set('priority', priority);
    window.location.href = params.toString() ? `/tickets?${params}` : '/tickets';
  }

  return (
    <form onSubmit={apply} className="flex flex-col gap-3 sm:flex-row sm:items-end">
      <div className="flex flex-col gap-2 flex-1">
        <Label htmlFor="ticket-search">Search</Label>
        <Input
          id="ticket-search"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search title or description..."
        />
      </div>
      <div className="flex flex-col gap-2 w-full sm:w-44">
        <Label>Status</Label>
        <Select value={status || 'all'} onValueChange={(v) => setStatus(v === 'all' ? '' : v ?? '')}>
          <SelectTrigger className="w-full">
            <SelectValue>
              {(value: string) => (value === 'all' ? 'All' : STATUS_OPTIONS.find((o) => o.value === value)?.label ?? value)}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All</SelectItem>
            {STATUS_OPTIONS.map((o) => (
              <SelectItem key={o.value} value={o.value}>
                {o.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="flex flex-col gap-2 w-full sm:w-44">
        <Label>Priority</Label>
        <Select value={priority || 'all'} onValueChange={(v) => setPriority(v === 'all' ? '' : v ?? '')}>
          <SelectTrigger className="w-full">
            <SelectValue>
              {(value: string) => (value === 'all' ? 'All' : PRIORITY_OPTIONS.find((o) => o.value === value)?.label ?? value)}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All</SelectItem>
            {PRIORITY_OPTIONS.map((o) => (
              <SelectItem key={o.value} value={o.value}>
                {o.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <Button type="submit">Filter</Button>
      {hasFilters && (
        <a href="/tickets" className="text-sm font-semibold text-muted-foreground hover:text-foreground px-2 py-2">
          Clear
        </a>
      )}
    </form>
  );
}
