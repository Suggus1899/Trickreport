import { useMemo } from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import type { VolumePoint, StatusDistribution, ResolutionMetrics } from '@/lib/api';

const CHART_COLORS = ['#6d5dfc', '#4cd964', '#ffcc00', '#ff5a5a', '#34a853', '#8b7dff'];

const PRIORITY_VARIANT: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
  low: 'outline',
  medium: 'secondary',
  high: 'secondary',
  critical: 'destructive',
};

interface Props {
  volume: VolumePoint[];
  status: StatusDistribution[];
  resolution: ResolutionMetrics[];
}

function toCsvValue(v: unknown): string {
  if (v === undefined || v === null) return '';
  const s = String(v);
  return /[",\n]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s;
}

export function AnalyticsCharts({ volume, status, resolution }: Props) {
  const chartVolume = useMemo(
    () => volume.map((v) => ({ ...v, label: new Date(v.date).toLocaleDateString(undefined, { month: 'short', day: 'numeric' }) })),
    [volume],
  );
  const total = status.reduce((s, x) => s + x.count, 0);

  function handleExport() {
    const rows: Record<string, unknown>[] = [];
    volume.forEach((v) => rows.push({ section: 'volume', date: v.date, count: v.count }));
    status.forEach((s) => rows.push({ section: 'status', status: s.status, count: s.count }));
    resolution.forEach((r) => rows.push({ section: 'resolution', priority: r.priority, avg_hours: r.avg_hours }));
    if (!rows.length) return;
    const headers = ['section', 'date', 'count', 'status', 'priority', 'avg_hours'];
    const lines = [headers.join(','), ...rows.map((row) => headers.map((h) => toCsvValue(row[h])).join(','))];
    const blob = new Blob([lines.join('\n')], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `analytics-${new Date().toISOString().slice(0, 10)}.csv`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex justify-end">
        <Button variant="outline" size="sm" onClick={handleExport}>Export CSV</Button>
      </div>

      <div className="rounded-xl border bg-card p-5">
        <h2 className="text-xl font-semibold mb-4">Ticket volume (trend)</h2>
        {chartVolume.length > 0 ? (
          <ResponsiveContainer width="100%" height={260}>
            <LineChart data={chartVolume} margin={{ top: 5, right: 20, left: 0, bottom: 5 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
              <XAxis dataKey="label" fontSize={11} stroke="var(--muted-foreground)" />
              <YAxis fontSize={11} stroke="var(--muted-foreground)" allowDecimals={false} />
              <Tooltip
                contentStyle={{ background: 'var(--popover)', border: '1px solid var(--border)', borderRadius: 8, fontSize: 12 }}
                labelFormatter={(_, payload) => (payload?.[0]?.payload ? new Date(payload[0].payload.date).toLocaleDateString() : '')}
              />
              <Line type="monotone" dataKey="count" stroke={CHART_COLORS[0]} strokeWidth={2.5} dot={{ r: 3 }} />
            </LineChart>
          </ResponsiveContainer>
        ) : (
          <p className="text-center text-muted-foreground py-6">No volume data available.</p>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="rounded-xl border bg-card p-5">
          <h2 className="text-xl font-semibold mb-4">Status distribution</h2>
          {status.length > 0 ? (
            <div className="flex flex-col md:flex-row items-center gap-4">
              <ResponsiveContainer width={200} height={200}>
                <PieChart>
                  <Pie data={status} dataKey="count" nameKey="status" innerRadius={48} outerRadius={80} paddingAngle={1}>
                    {status.map((_, i) => (
                      <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip contentStyle={{ background: 'var(--popover)', border: '1px solid var(--border)', borderRadius: 8, fontSize: 12 }} />
                </PieChart>
              </ResponsiveContainer>
              <div className="flex flex-col gap-2 flex-1">
                {status.map((s, i) => (
                  <div key={s.status} className="flex items-center gap-2">
                    <span className="w-2.5 h-2.5 rounded-full shrink-0" style={{ background: CHART_COLORS[i % CHART_COLORS.length] }} />
                    <span className="text-sm font-semibold capitalize">{s.status.replace('_', ' ')}</span>
                    <span className="text-xs text-muted-foreground">
                      {s.count} · {total ? ((s.count / total) * 100).toFixed(1) : 0}%
                    </span>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <p className="text-center text-muted-foreground py-6">No data available.</p>
          )}
        </div>

        <div className="rounded-xl border bg-card overflow-hidden">
          <div className="p-5 pb-0">
            <h2 className="text-xl font-semibold mb-4">Resolution time by priority</h2>
          </div>
          <table className="w-full text-left">
            <thead className="bg-muted/50 text-xs uppercase text-muted-foreground">
              <tr>
                <th className="px-5 py-3 font-semibold">Priority</th>
                <th className="px-5 py-3 font-semibold">Avg hours</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {resolution.length > 0 ? (
                resolution.map((r) => (
                  <tr key={r.priority} className="hover:bg-muted/30 transition-colors">
                    <td className="px-5 py-3">
                      <Badge variant={PRIORITY_VARIANT[r.priority] || 'default'}>{r.priority}</Badge>
                    </td>
                    <td className="px-5 py-3 font-semibold">{r.avg_hours.toFixed(1)}h</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={2} className="px-5 py-12 text-center text-muted-foreground">No data available.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
