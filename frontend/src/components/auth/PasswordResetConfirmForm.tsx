import { useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { confirmPasswordReset, ApiNetworkError, ApiValidationError } from '@/lib/api';

const MIN_PASSWORD_LENGTH = 12;

export function PasswordResetConfirmForm({ token }: { token: string }) {
  const [password, setPassword] = useState('');
  const [confirmation, setConfirmation] = useState('');
  const [done, setDone] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (password !== confirmation) {
      setError('The two passwords do not match.');
      return;
    }
    if (password.length < MIN_PASSWORD_LENGTH) {
      setError(`Use at least ${MIN_PASSWORD_LENGTH} characters.`);
      return;
    }

    setError('');
    setLoading(true);
    try {
      await confirmPasswordReset(token, password);
      setDone(true);
    } catch (err) {
      if (err instanceof ApiNetworkError) {
        setError('Could not reach the server. Check your connection and try again.');
      } else if (err instanceof ApiValidationError) {
        // The token is single-use and expires, so this is the common failure.
        setError('This reset link is no longer valid. Request a new one.');
      } else {
        setError('Something went wrong. Try again in a moment.');
      }
    } finally {
      setLoading(false);
    }
  }

  if (done) {
    return (
      <div className="flex flex-col gap-5" role="status">
        <p className="text-sm">Your password has been changed.</p>
        <a
          href="/login"
          className="inline-flex h-11 w-full items-center justify-center rounded-md bg-primary px-8 text-sm font-medium text-primary-foreground no-underline hover:bg-primary/90"
        >
          Sign in
        </a>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-5">
      {error && (
        <p className="text-sm text-destructive" role="alert">
          {error}
        </p>
      )}
      <div className="flex flex-col gap-2">
        <Label htmlFor="password">New password</Label>
        <Input
          id="password"
          type="password"
          autoComplete="new-password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength={MIN_PASSWORD_LENGTH}
          autoFocus
        />
        <p className="text-xs text-muted-foreground">At least {MIN_PASSWORD_LENGTH} characters, with uppercase, lowercase, a number and a symbol.</p>
      </div>
      <div className="flex flex-col gap-2">
        <Label htmlFor="confirmation">Confirm new password</Label>
        <Input
          id="confirmation"
          type="password"
          autoComplete="new-password"
          value={confirmation}
          onChange={(e) => setConfirmation(e.target.value)}
          required
        />
      </div>
      <Button type="submit" size="lg" disabled={loading} className="w-full">
        {loading ? 'Saving…' : 'Change password'}
      </Button>
    </form>
  );
}
