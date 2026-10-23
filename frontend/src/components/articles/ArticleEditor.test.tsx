// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { ArticleEditor } from './ArticleEditor';
import type { Article } from '@/lib/api';

const article: Article = { id: 'article-a', tenant_id: 'tenant-a', title: 'Printer setup', content: '</textarea><div id="content-marker">Text</div>', category: 'hardware', tags: ['printer'], published: false, created_by: 'user-a', created_at: '2026-10-20T10:00:00Z', updated_at: '2026-10-20T10:00:00Z' };

beforeEach(() => { document.cookie = 'trickreport_csrf=csrf; path=/'; });
afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.unstubAllGlobals(); });

describe('knowledge base editor', () => {
  it('keeps article content inert and restricts deletion to an admin', () => {
    render(<ArticleEditor article={article} role="agent" />);
    expect(screen.getByLabelText('Content')).toHaveValue(article.content);
    expect(document.querySelector('#content-marker')).toBeNull();
    expect(screen.queryByRole('button', { name: 'Delete' })).toBeNull();
  });

  it('does not delete when the admin cancels confirmation', () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false);
    const fetch = vi.fn();
    vi.stubGlobal('fetch', fetch);
    render(<ArticleEditor article={article} role="admin" />);
    fireEvent.click(screen.getByRole('button', { name: 'Delete' }));
    expect(fetch).not.toHaveBeenCalled();
    expect(screen.getByLabelText('Title')).toHaveValue('Printer setup');
  });

  it('retains the edited draft after a rejected publish request', async () => {
    let payload: Record<string, unknown> | undefined;
    vi.stubGlobal('fetch', async (_url: string, init: RequestInit) => {
      payload = JSON.parse(String(init.body));
      return new Response('{"error":"Article is not available"}', { status: 403 });
    });
    render(<ArticleEditor article={article} role="agent" />);
    fireEvent.change(screen.getByLabelText('Title'), { target: { value: ' New instructions ' } });
    fireEvent.change(screen.getByLabelText('Tags'), { target: { value: ' printer, setup, ' } });
    fireEvent.click(screen.getByRole('checkbox'));
    fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('Article is not available');
    expect(payload).toMatchObject({ title: 'New instructions', content: article.content, tags: ['printer', 'setup'], published: true });
    expect(screen.getByLabelText('Title')).toHaveValue(' New instructions ');
    expect(screen.getByRole('checkbox')).toBeChecked();
  });

  it('rejects a whitespace-only title before submitting', async () => {
    const fetch = vi.fn(async () => new Response('{"error":"unexpected request"}', { status: 422 }));
    vi.stubGlobal('fetch', fetch);
    render(<ArticleEditor article={article} role="agent" />);
    fireEvent.change(screen.getByLabelText('Title'), { target: { value: '   ' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));
    expect(await screen.findByRole('alert')).toHaveTextContent(/title.*3/i);
    expect(fetch).not.toHaveBeenCalled();
  });
});
