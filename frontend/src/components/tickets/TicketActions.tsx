import { useState, type FormEvent } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { Paperclip } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge, statusBadgeVariant, priorityBadgeVariant } from '@/components/ui/badge';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import {
  updateTicketStatus,
  assignTicket,
  addTicketComment,
  uploadTicketAttachment,
  type Ticket,
  type Comment,
  type Attachment,
  type User,
} from '@/lib/api';
import { PUBLIC_API_URL } from '@/lib/config';
import { toast } from '@/lib/toast';

export const STATUS_TRANSITIONS: Record<string, string[]> = {
  open: ['in_progress', 'waiting_client', 'resolved', 'closed'],
  in_progress: ['waiting_client', 'resolved', 'closed'],
  waiting_client: ['in_progress', 'resolved', 'closed'],
  resolved: ['closed', 'in_progress'],
  closed: ['in_progress'],
};

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

interface Props {
  ticketId: string;
  initialTicket: Ticket;
  initialComments: Comment[];
  initialAttachments: Attachment[];
  users: User[];
  canManage: boolean;
}

export function TicketActions({ ticketId, initialTicket, initialComments, initialAttachments, users, canManage }: Props) {
  const [ticket, setTicket] = useState(initialTicket);
  const [comments, setComments] = useState(initialComments);
  const [attachments, setAttachments] = useState(initialAttachments);

  const [statusValue, setStatusValue] = useState(ticket.status);
  const [statusLoading, setStatusLoading] = useState(false);
  const [statusError, setStatusError] = useState('');

  const [assignValue, setAssignValue] = useState(ticket.assigned_to || 'unassigned');
  const [assignLoading, setAssignLoading] = useState(false);
  const [assignError, setAssignError] = useState('');

  const [commentText, setCommentText] = useState('');
  const [commentInternal, setCommentInternal] = useState(false);
  const [commentLoading, setCommentLoading] = useState(false);
  const [commentError, setCommentError] = useState('');

  const [uploadLoading, setUploadLoading] = useState(false);
  const [uploadError, setUploadError] = useState('');

  const allowedTransitions = STATUS_TRANSITIONS[ticket.status] || [];
  const assignableUsers = users.filter((u) => u.role === 'agent' || u.role === 'admin');

  async function handleStatusSubmit(e: FormEvent) {
    e.preventDefault();
    setStatusError('');
    setStatusLoading(true);
    try {
      await updateTicketStatus(ticketId, statusValue);
      setTicket((t) => ({ ...t, status: statusValue }));
      toast(`Status updated to ${statusValue.replace('_', ' ')}`, 'success');
    } catch (err) {
      setStatusError(err instanceof Error ? err.message : 'Failed to update status');
    } finally {
      setStatusLoading(false);
    }
  }

  async function handleAssignSubmit(e: FormEvent) {
    e.preventDefault();
    setAssignError('');
    setAssignLoading(true);
    const newAssignee = assignValue === 'unassigned' ? '' : assignValue;
    try {
      await assignTicket(ticketId, newAssignee);
      const assignee = assignableUsers.find((u) => u.id === newAssignee);
      setTicket((t) => ({ ...t, assigned_to: newAssignee || undefined, assignee_name: assignee?.name }));
      toast(assignee ? `Assigned to ${assignee.name}` : 'Ticket unassigned', 'success');
    } catch (err) {
      setAssignError(err instanceof Error ? err.message : 'Failed to update assignment');
    } finally {
      setAssignLoading(false);
    }
  }

  async function handleCommentSubmit(e: FormEvent) {
    e.preventDefault();
    setCommentError('');
    setCommentLoading(true);
    try {
      const comment = await addTicketComment(ticketId, commentText, commentInternal);
      setComments((prev) => [...prev, comment]);
      setCommentText('');
      setCommentInternal(false);
      toast('Comment posted', 'success');
    } catch (err) {
      setCommentError(err instanceof Error ? err.message : 'Failed to add comment');
    } finally {
      setCommentLoading(false);
    }
  }

  async function handleUpload(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const fileInput = form.elements.namedItem('file') as HTMLInputElement;
    const file = fileInput.files?.[0];
    if (!file) return;
    setUploadError('');
    setUploadLoading(true);
    try {
      const attachment = await uploadTicketAttachment(ticketId, file);
      setAttachments((prev) => [...prev, attachment]);
      form.reset();
      toast(`Uploaded ${attachment.filename}`, 'success');
    } catch (err) {
      setUploadError(err instanceof Error ? err.message : 'Failed to upload attachment');
    } finally {
      setUploadLoading(false);
    }
  }

  return (
    <div className="flex flex-col gap-6">
      {/* Status/priority header — updates in place after status/assign changes */}
      <div className="flex flex-wrap items-center gap-3">
        <Badge variant={statusBadgeVariant(ticket.status)} dot>{ticket.status.replace('_', ' ')}</Badge>
        <Badge variant={priorityBadgeVariant(ticket.priority)} dot>{ticket.priority}</Badge>
        {ticket.sla_breached && <Badge variant="priority-critical" dot>SLA Breached</Badge>}
        <span className="text-sm text-muted-foreground ml-auto">
          {ticket.assignee_name ? `Assigned to ${ticket.assignee_name}` : 'Unassigned'}
        </span>
      </div>

      {canManage && (
        <div className="grid md:grid-cols-2 gap-6">
          <form onSubmit={handleStatusSubmit} className="rounded-xl border bg-card p-5 flex flex-col gap-3">
            <h2 className="text-sm font-semibold">Change status</h2>
            {statusError && <p className="text-sm text-destructive">{statusError}</p>}
            <Select value={statusValue} onValueChange={(value) => setStatusValue(value ?? ticket.status)}>
              <SelectTrigger className="w-full">
                <SelectValue>
                  {(value: string) => (value === ticket.status ? `${value.replace('_', ' ')} (current)` : value.replace('_', ' '))}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                <SelectItem value={ticket.status}>{ticket.status.replace('_', ' ')} (current)</SelectItem>
                {allowedTransitions.map((s) => (
                  <SelectItem key={s} value={s}>
                    {s.replace('_', ' ')}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <div className="flex justify-end">
              <Button type="submit" size="sm" disabled={statusLoading || statusValue === ticket.status}>
                {statusLoading ? 'Updating…' : 'Update status'}
              </Button>
            </div>
          </form>

          <form onSubmit={handleAssignSubmit} className="rounded-xl border bg-card p-5 flex flex-col gap-3">
            <h2 className="text-sm font-semibold">Assignment</h2>
            {assignError && <p className="text-sm text-destructive">{assignError}</p>}
            <Select value={assignValue} onValueChange={(value) => setAssignValue(value ?? 'unassigned')}>
              <SelectTrigger className="w-full">
                <SelectValue>
                  {(value: string) =>
                    value === 'unassigned' ? 'Unassigned' : assignableUsers.find((u) => u.id === value)?.name ?? value
                  }
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="unassigned">Unassigned</SelectItem>
                {assignableUsers.map((u) => (
                  <SelectItem key={u.id} value={u.id}>
                    {u.name} ({u.role})
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <div className="flex justify-end">
              <Button type="submit" size="sm" disabled={assignLoading}>
                {assignLoading ? 'Assigning…' : 'Assign'}
              </Button>
            </div>
          </form>
        </div>
      )}

      {/* Comments */}
      <div id="comments" className="rounded-xl border bg-card p-5">
        <h2 className="font-heading text-lg font-semibold mb-4">Comments ({comments.length})</h2>
        {comments.length > 0 ? (
          <div className="flex flex-col gap-4 mb-6">
            <AnimatePresence initial={false}>
              {comments.map((comment) => (
                <motion.div
                  key={comment.id}
                  initial={{ opacity: 0, y: 8 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.2, ease: 'easeOut' }}
                  className={`rounded-lg bg-muted/50 p-4 ${comment.is_internal ? 'border-l-2 border-warning' : ''}`}
                >
                  <div className="flex items-center justify-between mb-1.5">
                    <span className="font-semibold text-sm">{comment.user_name || comment.user_id}</span>
                    <span className="text-xs text-muted-foreground">{new Date(comment.created_at).toLocaleString()}</span>
                  </div>
                  <p className="text-sm whitespace-pre-wrap">{comment.content}</p>
                  {comment.is_internal && <span className="text-xs font-semibold text-warning mt-2 inline-block">Internal</span>}
                </motion.div>
              ))}
            </AnimatePresence>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground mb-6">No comments yet.</p>
        )}

        <h3 className="font-semibold text-sm mb-3">Add comment</h3>
        <form onSubmit={handleCommentSubmit} className="flex flex-col gap-3">
          {commentError && <p className="text-sm text-destructive">{commentError}</p>}
          <Textarea
            placeholder="Write a comment..."
            required
            rows={3}
            value={commentText}
            onChange={(e) => setCommentText(e.target.value)}
          />
          <label className="flex items-center gap-2 text-sm text-muted-foreground">
            <input
              type="checkbox"
              className="w-auto"
              checked={commentInternal}
              onChange={(e) => setCommentInternal(e.target.checked)}
            />
            Internal note (agents/admins only)
          </label>
          <div className="flex justify-end">
            <Button type="submit" size="sm" disabled={commentLoading || !commentText.trim()}>
              {commentLoading ? 'Posting…' : 'Post comment'}
            </Button>
          </div>
        </form>
      </div>

      {/* Attachments */}
      <div id="attachments" className="rounded-xl border bg-card p-5">
        <h2 className="font-heading text-lg font-semibold mb-4 flex items-center gap-2">
          <Paperclip className="size-4 text-muted-foreground" aria-hidden="true" />
          Attachments ({attachments.length})
        </h2>
        {attachments.length > 0 ? (
          <ul className="flex flex-col gap-3 mb-6">
            {attachments.map((att) => (
              <li key={att.id} className="rounded-lg bg-muted/50 p-4 flex items-center justify-between gap-3">
                <div className="min-w-0">
                  <p className="font-semibold text-sm truncate">{att.filename}</p>
                  <p className="text-xs text-muted-foreground">
                    {formatBytes(att.file_size)} · {att.content_type} · {att.user_name || att.user_id} ·{' '}
                    {new Date(att.created_at).toLocaleString()}
                  </p>
                </div>
                <a
                  href={`${PUBLIC_API_URL}/api/v1/tickets/${ticketId}/attachments/${att.id}`}
                  download={att.filename}
                  className="shrink-0 text-sm font-semibold text-primary hover:underline"
                >
                  Download
                </a>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-sm text-muted-foreground mb-6">No attachments.</p>
        )}

        <h3 className="font-semibold text-sm mb-3">Upload attachment</h3>
        <form onSubmit={handleUpload} className="flex flex-col gap-3">
          {uploadError && <p className="text-sm text-destructive">{uploadError}</p>}
          <input type="file" name="file" required aria-label="Choose file to upload" className="text-sm" />
          <div className="flex justify-end">
            <Button type="submit" size="sm" disabled={uploadLoading}>
              {uploadLoading ? 'Uploading…' : 'Upload'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
