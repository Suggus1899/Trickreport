import { useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { requestPasswordReset, ApiNetworkError } from '@/lib/api';

export function PasswordResetRequestForm() {
  const [email, setEmail] = useState('');
  const [sent, setSent] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await requestPasswordReset(email);
      setSent(true);
    } catch (err) {
      // The API answers 202 even for unknown addresses, so anything that
      // lands here is a transport or server problem worth showing.
      setError(
        err instanceof ApiNetworkError
          ? 'Could not reach the server. Check your connection and try again.'
          : 'Something went wrong. Try again in a moment.',
      );
    } finally {
      setLoading(false);
    }
  }

  if (sent) {
    return (
      <div className="flex flex-col gap-4" role="status">
        <p className="text-sm">
          If an account exists for <span className="font-semibold">{email}</span>, a reset link is on
          its way. The link expires shortly, so use it soon.
        </p>
        <p className="text-sm text-muted-foreground">
          Nothing arrived? Check your spam folder, or{' '}
          <button
            type="button"
            className="font-semibold text-primary hover:underline"
            onClick={() => setSent(false)}
          >
            try another address
          </button>
          .
        </p>
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
        <Label htmlFor="email">Email</Label>
        <Input
          id="email"
          type="email"
          placeholder="you@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          autoFocus
        />
        <p className="text-xs text-muted-foreground">We'll send a link to reset your password.</p>
      </div>
      <Button type="submit" size="lg" disabled={loading} className="w-full">
        {loading ? 'Sending…' : 'Send reset link'}
      </Button>
    </form>
  );
}
