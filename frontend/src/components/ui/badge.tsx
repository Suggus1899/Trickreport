import { mergeProps } from "@base-ui/react/merge-props"
import { useRender } from "@base-ui/react/use-render"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const badgeVariants = cva(
  "group/badge inline-flex h-5 w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-4xl border border-transparent px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-all focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&>svg]:pointer-events-none [&>svg]:size-3!",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground [a]:hover:bg-primary/80",
        secondary:
          "bg-secondary text-secondary-foreground [a]:hover:bg-secondary/80",
        destructive:
          "bg-destructive/10 text-destructive focus-visible:ring-destructive/20 dark:bg-destructive/20 dark:focus-visible:ring-destructive/40 [a]:hover:bg-destructive/20",
        outline:
          "border-border text-foreground [a]:hover:bg-muted [a]:hover:text-muted-foreground",
        ghost:
          "hover:bg-muted hover:text-muted-foreground dark:hover:bg-muted/50",
        link: "text-primary underline-offset-4 hover:underline",
        // Estado de ticket. El color es la señal: solo "critical"/SLA breach
        // usa el acento coral (destructive) — el resto son semánticos, nunca
        // decorativos.
        "status-open": "bg-primary/10 text-primary dark:bg-primary/15",
        "status-in-progress": "bg-warning/10 text-warning dark:bg-warning/15",
        "status-waiting-client": "bg-muted text-muted-foreground",
        "status-resolved": "bg-success/10 text-success dark:bg-success/15",
        "status-closed": "bg-muted text-muted-foreground opacity-70",
        "priority-low": "bg-muted text-muted-foreground",
        "priority-medium": "bg-primary/10 text-primary dark:bg-primary/15",
        "priority-high": "bg-warning/10 text-warning dark:bg-warning/15",
        "priority-critical":
          "bg-destructive/10 text-destructive dark:bg-destructive/20",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

/** Ticket status -> badge variant. Keep in sync with the backend's status enum. */
function statusBadgeVariant(status: string): VariantProps<typeof badgeVariants>["variant"] {
  switch (status) {
    case "open":
      return "status-open"
    case "in_progress":
      return "status-in-progress"
    case "waiting_client":
      return "status-waiting-client"
    case "resolved":
      return "status-resolved"
    case "closed":
      return "status-closed"
    default:
      return "outline"
  }
}

/** Ticket priority -> badge variant. Keep in sync with the backend's priority enum. */
function priorityBadgeVariant(priority: string): VariantProps<typeof badgeVariants>["variant"] {
  switch (priority) {
    case "low":
      return "priority-low"
    case "medium":
      return "priority-medium"
    case "high":
      return "priority-high"
    case "critical":
      return "priority-critical"
    default:
      return "outline"
  }
}

/**
 * Ticket status -> raw color for chart libraries that need a literal fill
 * (recharts Cell/Line strokes can't consume Tailwind classes). Kept in sync
 * with statusBadgeVariant above so a chart and a badge always agree on what
 * color a given status means. "closed" gets its own neutral rather than
 * reusing muted-foreground so it reads as distinct in a legend.
 */
function statusChartColor(status: string): string {
  switch (status) {
    case "open":
      return "var(--primary)"
    case "in_progress":
      return "var(--warning)"
    case "waiting_client":
      return "var(--muted-foreground)"
    case "resolved":
      return "var(--success)"
    case "closed":
      return "color-mix(in oklch, var(--muted-foreground), var(--foreground) 35%)"
    default:
      return "var(--muted-foreground)"
  }
}

function Badge({
  className,
  variant = "default",
  dot = false,
  render,
  children,
  ...props
}: useRender.ComponentProps<"span"> & VariantProps<typeof badgeVariants> & { dot?: boolean }) {
  return useRender({
    defaultTagName: "span",
    props: mergeProps<"span">(
      {
        className: cn(badgeVariants({ variant }), className),
      },
      {
        ...props,
        children: dot ? (
          <>
            <span aria-hidden="true" className="size-1.5 shrink-0 rounded-full bg-current" />
            {children}
          </>
        ) : (
          children
        ),
      }
    ),
    render,
    state: {
      slot: "badge",
      variant,
    },
  })
}

export { Badge, badgeVariants, statusBadgeVariant, priorityBadgeVariant, statusChartColor }
