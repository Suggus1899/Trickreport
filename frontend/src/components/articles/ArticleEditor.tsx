import { useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  createArticle,
  updateArticle,
  deleteArticle,
  ApiNetworkError,
  ApiValidationError,
  type Article,
} from '@/lib/api';

// Mirrors the backend validation on createArticleReq/updateArticleReq, so the
// form rejects what the API would reject anyway — without a round trip.
const LIMITS = {
  titleMin: 3,
  titleMax: 255,
  contentMin: 10,
  contentMax: 20000,
  categoryMax: 100,
} as const;

interface Props {
  /** Omitted when creating; present when editing an existing article. */
  article?: Article;
}

export function ArticleEditor({ article }: Props) {
  const isEdit = Boolean(article);

  const [title, setTitle] = useState(article?.title ?? '');
  const [category, setCategory] = useState(article?.category ?? '');
  const [tags, setTags] = useState((article?.tags ?? []).join(', '));
  const [content, setContent] = useState(article?.content ?? '');
  const [published, setPublished] = useState(article?.published ?? false);

  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);

  function describe(err: unknown): string {
    if (err instanceof ApiNetworkError) {
      return 'Could not reach the server. Your changes were not saved.';
    }
    if (err instanceof ApiValidationError) {
      return err.message;
    }
    return err instanceof Error ? err.message : 'Something went wrong.';
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');
    setSaving(true);

    const payload = {
      title: title.trim(),
      content: content.trim(),
      category: category.trim(),
      // Empty entries would fail the backend's `dive` validation.
      tags: tags
        .split(',')
        .map((t) => t.trim())
        .filter(Boolean),
      published,
    };

    try {
      const saved = isEdit ? await updateArticle(article!.id, payload) : await createArticle(payload);
      window.location.href = `/articles/${saved.id}`;
    } catch (err) {
      setError(describe(err));
      setSaving(false);
    }
  }

  async function handleDelete() {
    if (!article) return;
    if (!window.confirm(`Delete "${article.title}"? This cannot be undone.`)) return;

    setError('');
    setDeleting(true);
    try {
      await deleteArticle(article.id);
      window.location.href = '/articles';
    } catch (err) {
      setError(describe(err));
      setDeleting(false);
    }
  }

  const busy = saving || deleting;

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-6">
      {error && (
        <p
          className="rounded-lg border border-destructive/30 bg-destructive/5 p-4 text-sm text-destructive"
          role="alert"
        >
          {error}
        </p>
      )}

      <div className="flex flex-col gap-2">
        <Label htmlFor="title">Title</Label>
        <Input
          id="title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          minLength={LIMITS.titleMin}
          maxLength={LIMITS.titleMax}
          autoFocus={!isEdit}
        />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
        <div className="flex flex-col gap-2">
          <Label htmlFor="category">Category</Label>
          <Input
            id="category"
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            required
            maxLength={LIMITS.categoryMax}
            placeholder="billing"
          />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="tags">Tags</Label>
          <Input
            id="tags"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
            placeholder="refunds, invoices"
          />
          <p className="text-xs text-muted-foreground">Separate with commas.</p>
        </div>
      </div>

      <div className="flex flex-col gap-2">
        <Label htmlFor="content">Content</Label>
        <Textarea
          id="content"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          required
          minLength={LIMITS.contentMin}
          maxLength={LIMITS.contentMax}
          rows={16}
          className="font-mono text-sm"
        />
        <p className="text-xs text-muted-foreground">
          Line breaks are preserved. {content.length.toLocaleString()} /{' '}
          {LIMITS.contentMax.toLocaleString()}
        </p>
      </div>

      <label className="flex items-start gap-3 rounded-lg border p-4 cursor-pointer">
        <input
          type="checkbox"
          checked={published}
          onChange={(e) => setPublished(e.target.checked)}
          className="mt-0.5 size-4"
        />
        <span className="flex flex-col gap-0.5">
          <span className="text-sm font-semibold">Published</span>
          <span className="text-xs text-muted-foreground">
            Drafts stay visible to admins and agents only. End users see published articles.
          </span>
        </span>
      </label>

      <div className="flex flex-wrap items-center justify-end gap-3 pt-2 border-t">
        {isEdit && (
          <Button
            type="button"
            variant="outline"
            onClick={handleDelete}
            disabled={busy}
            className="mr-auto text-destructive hover:text-destructive"
          >
            {deleting ? 'Deleting…' : 'Delete'}
          </Button>
        )}
        <a
          href={isEdit ? `/articles/${article!.id}` : '/articles'}
          className="text-sm font-semibold text-muted-foreground hover:text-foreground px-4 py-2"
        >
          Cancel
        </a>
        <Button type="submit" disabled={busy}>
          {saving ? 'Saving…' : isEdit ? 'Save changes' : 'Create article'}
        </Button>
      </div>
    </form>
  );
}
