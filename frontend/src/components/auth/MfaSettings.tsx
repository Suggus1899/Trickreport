import { useEffect, useRef, useState, type FormEvent } from 'react';
import QRCode from 'qrcode';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { mfaSetup, mfaEnable, mfaDisable } from '@/lib/api';

type Step = 'idle' | 'setup' | 'disable';

export function MfaSettings({ initialEnabled }: { initialEnabled: boolean }) {
  const [enabled, setEnabled] = useState(initialEnabled);
  const [step, setStep] = useState<Step>('idle');
  const [secret, setSecret] = useState('');
  const [qrDataUrl, setQrDataUrl] = useState('');
  const [code, setCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (step === 'setup' && secret && canvasRef.current) {
      QRCode.toCanvas(canvasRef.current, qrDataUrl, { width: 200 }).catch(() => {});
    }
  }, [step, secret, qrDataUrl]);

  async function startSetup() {
    setError('');
    setLoading(true);
    try {
      const result = await mfaSetup();
      setSecret(result.secret);
      setQrDataUrl(result.qr_url);
      setCode('');
      setStep('setup');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to start MFA setup');
    } finally {
      setLoading(false);
    }
  }

  async function confirmEnable(e: FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await mfaEnable(secret, code);
      setEnabled(true);
      setStep('idle');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Invalid code');
    } finally {
      setLoading(false);
    }
  }

  async function confirmDisable(e: FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await mfaDisable(code);
      setEnabled(false);
      setStep('idle');
      setCode('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Invalid code');
    } finally {
      setLoading(false);
    }
  }

  if (step === 'setup') {
    return (
      <form onSubmit={confirmEnable} className="flex flex-col gap-4 mt-4">
        {error && <p className="text-sm text-destructive" role="alert">{error}</p>}
        <canvas ref={canvasRef} className="rounded-lg border self-start" />
        <p className="text-xs text-muted-foreground">
          Scan this with your authenticator app, or enter the secret manually:{' '}
          <code className="rounded bg-muted px-1.5 py-0.5 font-mono">{secret}</code>
        </p>
        <div className="flex flex-col gap-2 max-w-xs">
          <Label htmlFor="mfa-enable-code">Enter the 6-digit code to confirm</Label>
          <Input
            id="mfa-enable-code"
            inputMode="numeric"
            autoComplete="one-time-code"
            placeholder="123456"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            required
            autoFocus
          />
        </div>
        <div className="flex gap-2">
          <Button type="submit" disabled={loading}>
            {loading ? 'Confirming…' : 'Confirm and enable'}
          </Button>
          <Button type="button" variant="ghost" onClick={() => setStep('idle')}>
            Cancel
          </Button>
        </div>
      </form>
    );
  }

  if (step === 'disable') {
    return (
      <form onSubmit={confirmDisable} className="flex flex-col gap-4 mt-4 max-w-xs">
        {error && <p className="text-sm text-destructive" role="alert">{error}</p>}
        <div className="flex flex-col gap-2">
          <Label htmlFor="mfa-disable-code">Enter your current authentication code</Label>
          <Input
            id="mfa-disable-code"
            inputMode="numeric"
            autoComplete="one-time-code"
            placeholder="123456"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            required
            autoFocus
          />
        </div>
        <div className="flex gap-2">
          <Button type="submit" variant="destructive" disabled={loading}>
            {loading ? 'Disabling…' : 'Disable MFA'}
          </Button>
          <Button type="button" variant="ghost" onClick={() => setStep('idle')}>
            Cancel
          </Button>
        </div>
      </form>
    );
  }

  return (
    <div className="mt-4 flex items-center gap-3 flex-wrap">
      <p className="text-sm text-muted-foreground">
        {enabled ? 'MFA is currently enabled on your account.' : 'MFA is not enabled yet.'}
      </p>
      {enabled ? (
        <Button variant="destructive" size="sm" onClick={() => setStep('disable')}>
          Disable MFA
        </Button>
      ) : (
        <Button size="sm" onClick={startSetup} disabled={loading}>
          {loading ? 'Starting…' : 'Set up MFA'}
        </Button>
      )}
      {error && <p className="text-sm text-destructive w-full" role="alert">{error}</p>}
    </div>
  );
}
