import { useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { login, mfaLogin } from '@/lib/api';

export function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [code, setCode] = useState('');
  const [mfaToken, setMfaToken] = useState<string | null>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleCredentials(e: FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const result = await login({ email, password });
      if (result.requires_mfa && result.mfa_token) {
        setMfaToken(result.mfa_token);
        return;
      }
      window.location.href = '/dashboard';
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  }

  async function handleMfaCode(e: FormEvent) {
    e.preventDefault();
    if (!mfaToken) return;
    setError('');
    setLoading(true);
    try {
      await mfaLogin(mfaToken, code);
      window.location.href = '/dashboard';
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Invalid code');
    } finally {
      setLoading(false);
    }
  }

  if (mfaToken) {
    return (
      <form onSubmit={handleMfaCode} className="flex flex-col gap-5">
        {error && <p className="text-sm text-destructive" role="alert">{error}</p>}
        <div className="flex flex-col gap-2">
          <Label htmlFor="mfa-code">Authentication code</Label>
          <Input
            id="mfa-code"
            inputMode="numeric"
            autoComplete="one-time-code"
            placeholder="123456"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            required
            autoFocus
          />
          <p className="text-xs text-muted-foreground">Enter the 6-digit code from your authenticator app.</p>
        </div>
        <Button type="submit" size="lg" disabled={loading} className="w-full">
          {loading ? 'Verifying…' : 'Verify'}
        </Button>
        <button
          type="button"
          className="text-sm text-muted-foreground hover:text-foreground underline-offset-4 hover:underline"
          onClick={() => {
            setMfaToken(null);
            setCode('');
            setError('');
          }}
        >
          Back to login
        </button>
      </form>
    );
  }

  return (
    <form onSubmit={handleCredentials} className="flex flex-col gap-5">
      {error && <p className="text-sm text-destructive" role="alert">{error}</p>}
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
      </div>
      <div className="flex flex-col gap-2">
        <Label htmlFor="password">Password</Label>
        <Input
          id="password"
          type="password"
          placeholder="Enter your password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </div>
      <div className="flex justify-end -mt-2">
        <a href="/password-reset" className="text-xs text-muted-foreground hover:text-foreground hover:underline">
          Forgot password?
        </a>
      </div>
      <Button type="submit" size="lg" disabled={loading} className="w-full">
        {loading ? 'Signing in…' : 'Sign In'}
      </Button>
    </form>
  );
}
