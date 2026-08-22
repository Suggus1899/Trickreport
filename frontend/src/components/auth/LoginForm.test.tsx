// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { LoginForm } from './LoginForm';

const loginMock = vi.fn();
const mfaLoginMock = vi.fn();

vi.mock('@/lib/api', () => ({
  login: (...args: unknown[]) => loginMock(...args),
  mfaLogin: (...args: unknown[]) => mfaLoginMock(...args),
}));

beforeEach(() => {
  loginMock.mockReset();
  mfaLoginMock.mockReset();
  // @ts-expect-error jsdom doesn't implement navigation
  delete window.location;
  // @ts-expect-error minimal stub, only `.href` is used by the component
  window.location = { href: '' };
});

describe('LoginForm', () => {
  it('navigates to the dashboard on a normal successful login', async () => {
    loginMock.mockResolvedValue({ token: 'jwt', user: { id: '1' } });
    render(<LoginForm />);

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'a@b.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

    await waitFor(() => expect(window.location.href).toBe('/dashboard'));
    expect(mfaLoginMock).not.toHaveBeenCalled();
  });

  it('switches to the TOTP step when the server requires MFA, without navigating away', async () => {
    loginMock.mockResolvedValue({ requires_mfa: true, mfa_token: 'temp-token' });
    render(<LoginForm />);

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'a@b.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

    expect(await screen.findByLabelText(/authentication code/i)).toBeInTheDocument();
    expect(window.location.href).toBe('');
  });

  it('submits the TOTP code against the mfa_token from step one, then navigates', async () => {
    loginMock.mockResolvedValue({ requires_mfa: true, mfa_token: 'temp-token' });
    mfaLoginMock.mockResolvedValue({ token: 'jwt', user: { id: '1' } });
    render(<LoginForm />);

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'a@b.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

    const codeInput = await screen.findByLabelText(/authentication code/i);
    fireEvent.change(codeInput, { target: { value: '123456' } });
    fireEvent.click(screen.getByRole('button', { name: /verify/i }));

    await waitFor(() => expect(mfaLoginMock).toHaveBeenCalledWith('temp-token', '123456'));
    await waitFor(() => expect(window.location.href).toBe('/dashboard'));
  });

  it('shows an error and stays on the TOTP step when the code is rejected', async () => {
    loginMock.mockResolvedValue({ requires_mfa: true, mfa_token: 'temp-token' });
    mfaLoginMock.mockRejectedValue(new Error('invalid mfa code or token'));
    render(<LoginForm />);

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'a@b.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

    const codeInput = await screen.findByLabelText(/authentication code/i);
    fireEvent.change(codeInput, { target: { value: '000000' } });
    fireEvent.click(screen.getByRole('button', { name: /verify/i }));

    expect(await screen.findByRole('alert')).toHaveTextContent('invalid mfa code or token');
    expect(window.location.href).toBe('');
  });
});
