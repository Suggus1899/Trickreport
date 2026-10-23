'use client';
import { createContext, useContext, useEffect, useRef, useState, type ReactNode } from 'react';
import { Button } from '@heroui/react';
import type { Locale } from './contracts';
export { Button };
export const Language = createContext<Locale>('es');
export function useText() { const locale = useContext(Language); return (es: string, en: string) => locale === 'es' ? es : en; }
export function Field({ label, ...props }: React.InputHTMLAttributes<HTMLInputElement> & { label: string }) { return <label className="field"><span>{label}</span><input {...props} /></label>; }
export function Select({ label, children, ...props }: React.SelectHTMLAttributes<HTMLSelectElement> & { label: string }) { return <label className="field"><span>{label}</span><select {...props}>{children}</select></label>; }
export function Textarea({ label, ...props }: React.TextareaHTMLAttributes<HTMLTextAreaElement> & { label: string }) { return <label className="field"><span>{label}</span><textarea {...props} /></label>; }
export function Notice({ children, error = false }: { children: ReactNode; error?: boolean }) { return <div role={error ? 'alert' : 'status'} className={`notice ${error ? 'notice-error' : ''}`}>{children}</div>; }
export function Empty({ children }: { children: ReactNode }) { return <div className="empty">{children}</div>; }
export function Heading({ title, description, action }: { title: string; description?: string; action?: ReactNode }) { return <header className="page-heading"><div><p className="eyebrow">WORKSPACE</p><h1>{title}</h1>{description && <p>{description}</p>}</div>{action}</header>; }
export function Pill({ value }: { value: string }) { const t = useText(); const labels: Record<string, [string,string]> = { open:['Abierto','Open'], in_progress:['En curso','In progress'], waiting_client:['Esperando cliente','Waiting for customer'], resolved:['Resuelto','Resolved'], closed:['Cerrado','Closed'], low:['Baja','Low'], medium:['Media','Medium'], high:['Alta','High'], critical:['Crítica','Critical'], admin:['Administrador','Administrator'], agent:['Agente','Agent'], end_user:['Solicitante','Requester'] }; return <span className={`pill pill-${value}`}>{labels[value] ? t(...labels[value]) : value}</span>; }
export function useLoad<T>(load: () => Promise<T>, deps: React.DependencyList = []) {
  const loadRef = useRef(load);
  const [result, setResult] = useState<{ key: string; data?: T; error: string }>({ key: '', error: '' });
  const [revision, setRevision] = useState(0);
  const key = `${JSON.stringify(deps)}:${revision}`;
  useEffect(() => { loadRef.current = load; }, [load]);
  useEffect(() => {
    let active = true;
    loadRef.current().then(data => { if (active) setResult({ key, data, error: '' }); }).catch(error => {
      if (active) setResult({ key, error: error instanceof Error ? error.message : 'Error' });
    });
    return () => { active = false; };
  }, [key]);
  const loading = result.key !== key;
  return { data: loading ? undefined : result.data, error: loading ? '' : result.error, loading, reload: () => setRevision(value => value + 1) };
}
export function LoadState({ loading, error, retry, children }: { loading: boolean; error: string; retry: () => void; children: ReactNode }) { const t=useText(); if(loading) return <div className="loading" role="status"><span className="loading-dot" />{t('Cargando información…','Loading information…')}</div>; if(error) return <Notice error>{error} <Button variant="secondary" onPress={retry}>{t('Reintentar','Retry')}</Button></Notice>; return children; }
export function useMutation() { const [error,setError]=useState(''); const [busy,setBusy]=useState(false); const [success,setSuccess]=useState(''); async function run(action:()=>Promise<void>,message='') { setError('');setSuccess('');setBusy(true);try{await action();setSuccess(message);}catch(err){setError(err instanceof Error?err.message:'Error');}finally{setBusy(false);} } return {error,busy,success,run}; }
export function Feedback({ error, success }: { error: string; success: string }) { return <>{error && <Notice error>{error}</Notice>}{success && <Notice>{success}</Notice>}</>; }
