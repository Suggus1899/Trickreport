// @vitest-environment jsdom
import { afterEach, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { RegisterForm } from './RegisterForm';

afterEach(() => { cleanup(); vi.unstubAllGlobals(); });
it('explains the 12-character policy before submitting a short registration password', async () => {
  const fetch = vi.fn(async () => { throw new TypeError('Unexpected request'); });
  vi.stubGlobal('fetch', fetch);
  document.cookie = 'trickreport_csrf=csrf; path=/';
  render(<RegisterForm />);
  fireEvent.change(screen.getByLabelText('Full name'), { target: { value: 'Reader' } });
  fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'reader@test.test' } });
  fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'Weakpass1!' } });
  fireEvent.change(screen.getByLabelText('Confirm password'), { target: { value: 'Weakpass1!' } });
  fireEvent.submit(screen.getByRole('button', { name: /create account/i }).closest('form')!);
  expect(await screen.findByRole('alert')).toHaveTextContent(/12/);
  expect(fetch).not.toHaveBeenCalled();
});
