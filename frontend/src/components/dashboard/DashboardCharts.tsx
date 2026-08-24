import { useId, useMemo } from 'react';
import { motion } from 'motion/react';
import { AreaChart, Area, XAxis, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { ArrowDownRight, ArrowUpRight, Minus } from 'lucide-react';
import { Badge, statusBadgeVariant, priorityBadgeVariant, statusChartColor } from '@/components/ui/badge';
import type { Summary, Ticket } from '@/lib/api';
import type { DailyPoint, StatusSlice } from '@/lib/dashboard-metrics';
import { weekOverWeekDelta } from '@/lib/dashboard-metrics';

interface Props {
  summary: Summary;
  totalVolume: DailyPoint[];
  openVolume: DailyPoint[];
  resolvedVolume: DailyPoint[];
  slaVolume: DailyPoint[];
  distribution: StatusSlice[];
  myTickets: Ticket[];
}

function Sparkline({ points, color }: { points: DailyPoint[]; color: string }) {
  const id = useId();
  const max = Math.max(1, ...points.map((p) => p.count));
  const w = 100;
  const h = 28;
  const step = points.length > 1 ? w / (points.length - 1) : w;
  const coords = points.map((p, i) => [i * step, h - (p.count / max) * (h - 4) - 2] as const);
  const path = coords.map(([x, y], i) => `${i === 0 ? 'M' : 'L'}${x.toFixed(1)},${y.toFixed(1)}`).join(' ');
  const area = `${path} L${w},${h} L0,${h} Z`;
  return (
    <svg viewBox={`0 0 ${w} ${h}`} className="w-full h-7" preserveAspectRatio="none" aria-hidden="true">
      <defs>
        <linearGradient id={id} x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor={color} stopOpacity="0.35" />
          <stop offset="100%" stopColor={color} stopOpacity="0" />
        </linearGradient>
      </defs>
      <path d={area} fill={`url(#${id})`} stroke="none" />
      <path d={path} fill="none" stroke={color} strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

/**
 * `goodDirection` says which way is good news for this particular stat — an
 * uptick in Resolved is good, an uptick in SLA Breached isn't. `neutral`
 * skips the judgment call entirely (more total tickets isn't inherently
 * good or bad) and just reports the number.
 */
function DeltaChip({ value, goodDirection }: { value: number; goodDirection: 'up' | 'down' | 'neutral' }) {
  if (value === 0) {
    return (
      <span className="inline-flex items-center gap-0.5 text-xs font-medium text-muted-foreground">
        <Minus className="size-3" /> flat
      </span>
    );
  }
  const up = value > 0;
  const isGood = goodDirection === 'neutral' ? null : up === (goodDirection === 'up');
  const colorClass = isGood === null ? 'text-muted-foreground' : isGood ? 'text-success' : 'text-warning';
  return (
    <span className={`inline-flex items-center gap-0.5 text-xs font-medium ${colorClass}`}>
      {up ? <ArrowUpRight className="size-3" /> : <ArrowDownRight className="size-3" />}
      {Math.abs(value)}%
    </span>
  );
}

const stagger = {
  hidden: {},
  show: { transition: { staggerChildren: 0.06 } },
};
const item = {
  hidden: { opacity: 0, y: 8 },
  show: { opacity: 1, y: 0, transition: { duration: 0.25, ease: 'easeOut' as const } },
};

export function DashboardCharts({
  summary,
  totalVolume,
  openVolume,
  resolvedVolume,
  slaVolume,
  distribution,
  myTickets,
}: Props) {
  const chartVolume = useMemo(
    () => totalVolume.map((v) => ({ ...v, label: new Date(v.date).toLocaleDateString(undefined, { month: 'short', day: 'numeric' }) })),
    [totalVolume],
  );
  const totalCount = distribution.reduce((s, x) => s + x.count, 0);
  const slaActive = summary.sla_breached > 0;

  const stats = [
    { label: 'Total', value: summary.total_tickets, points: totalVolume, color: 'var(--primary)', tone: 'text-foreground', goodDirection: 'neutral' as const },
    { label: 'Open', value: summary.open_tickets, points: openVolume, color: 'var(--warning)', tone: 'text-warning', goodDirection: 'down' as const },
    { label: 'Resolved', value: summary.resolved_tickets, points: resolvedVolume, color: 'var(--success)', tone: 'text-success', goodDirection: 'up' as const },
    { label: 'SLA Breached', value: summary.sla_breached, points: slaVolume, color: 'var(--destructive)', tone: 'text-destructive', goodDirection: 'down' as const },
  ];

  return (
    <motion.div variants={stagger} initial="hidden" animate="show" className="flex flex-col gap-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((s) => (
          <motion.div
            key={s.label}
            variants={item}
            className={`relative overflow-hidden rounded-xl border bg-card p-5 flex flex-col gap-1 ${
              s.label === 'SLA Breached' && slaActive ? 'border-destructive/40' : ''
            }`}
          >
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wide font-heading">{s.label}</span>
              {s.label === 'SLA Breached' && slaActive && (
                <span className="relative flex size-2" aria-hidden="true">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-destructive opacity-75" />
                  <span className="relative inline-flex size-2 rounded-full bg-destructive" />
                </span>
              )}
            </div>
            <div className="flex items-end justify-between gap-2">
              <span className={`font-heading text-4xl font-bold tabular-nums ${s.tone}`}>{s.value}</span>
              <DeltaChip value={weekOverWeekDelta(s.points)} goodDirection={s.goodDirection} />
            </div>
            <Sparkline points={s.points} color={s.color} />
          </motion.div>
        ))}
      </div>

      <motion.div variants={item} className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="lg:col-span-2 rounded-xl border bg-card p-5">
          <h2 className="font-heading font-semibold mb-4">Ticket volume — last 14 days</h2>
          {chartVolume.some((v) => v.count > 0) ? (
            <ResponsiveContainer width="100%" height={220}>
              <AreaChart data={chartVolume} margin={{ top: 5, right: 12, left: 0, bottom: 0 }}>
                <defs>
                  <linearGradient id="dash-volume" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="var(--primary)" stopOpacity="0.25" />
                    <stop offset="100%" stopColor="var(--primary)" stopOpacity="0" />
                  </linearGradient>
                </defs>
                <XAxis dataKey="label" fontSize={11} stroke="var(--muted-foreground)" tickLine={false} axisLine={false} />
                <Tooltip
                  contentStyle={{ background: 'var(--popover)', border: '1px solid var(--border)', borderRadius: 6, fontSize: 12 }}
                  labelFormatter={(_, payload) => (payload?.[0]?.payload ? new Date(payload[0].payload.date).toLocaleDateString() : '')}
                />
                <Area type="monotone" dataKey="count" stroke="var(--primary)" strokeWidth={2} fill="url(#dash-volume)" />
              </AreaChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-center text-muted-foreground py-10 text-sm">No tickets created in the last 14 days.</p>
          )}
        </div>

        <div className="rounded-xl border bg-card p-5">
          <h2 className="font-heading font-semibold mb-4">Status distribution</h2>
          {totalCount > 0 ? (
            <div className="flex flex-col items-center gap-4">
              <ResponsiveContainer width={140} height={140}>
                <PieChart>
                  <Pie data={distribution} dataKey="count" nameKey="status" innerRadius={38} outerRadius={62} paddingAngle={2}>
                    {distribution.map((s) => (
                      <Cell key={s.status} fill={statusChartColor(s.status)} />
                    ))}
                  </Pie>
                  <Tooltip contentStyle={{ background: 'var(--popover)', border: '1px solid var(--border)', borderRadius: 6, fontSize: 12 }} />
                </PieChart>
              </ResponsiveContainer>
              <div className="w-full flex flex-col gap-1.5">
                {distribution.filter((s) => s.count > 0).map((s) => (
                  <div key={s.status} className="flex items-center gap-2 text-xs">
                    <span className="size-2 rounded-full shrink-0" style={{ background: statusChartColor(s.status) }} aria-hidden="true" />
                    <span className="font-medium capitalize flex-1">{s.status.replace('_', ' ')}</span>
                    <span className="text-muted-foreground tabular-nums">{s.count}</span>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <p className="text-center text-muted-foreground py-10 text-sm">No tickets yet.</p>
          )}
        </div>
      </motion.div>

      <motion.div variants={item} className="rounded-xl border bg-card p-5">
        <h2 className="font-heading font-semibold mb-4">My tickets</h2>
        {myTickets.length > 0 ? (
          <ul className="flex flex-col divide-y">
            {myTickets.map((t) => (
              <li key={t.id}>
                <a
                  href={`/tickets/${t.id}`}
                  className={`flex items-center gap-3 py-3 no-underline text-foreground hover:bg-muted/50 -mx-2 px-2 rounded-md transition-colors ${
                    t.sla_breached ? 'border-l-2 border-destructive pl-3 -ml-px' : ''
                  }`}
                >
                  <span className="flex-1 min-w-0 truncate text-sm font-medium">{t.title}</span>
                  <Badge variant={priorityBadgeVariant(t.priority)} dot className="hidden sm:inline-flex">
                    {t.priority}
                  </Badge>
                  <Badge variant={statusBadgeVariant(t.status)} dot>
                    {t.status.replace('_', ' ')}
                  </Badge>
                </a>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-center text-muted-foreground py-6 text-sm">No tickets assigned to you.</p>
        )}
      </motion.div>
    </motion.div>
  );
}
