import { useEffect, useState } from 'react';

export function useOnlineStatus(): boolean {
  // Always start "online" to match the server-rendered markup — the real
  // value (which can legitimately differ, that's the whole point of a PWA)
  // is only read after hydration, inside the effect below.
  const [online, setOnline] = useState(true);

  useEffect(() => {
    setOnline(navigator.onLine);
    const setTrue = () => setOnline(true);
    const setFalse = () => setOnline(false);
    window.addEventListener('online', setTrue);
    window.addEventListener('offline', setFalse);
    return () => {
      window.removeEventListener('online', setTrue);
      window.removeEventListener('offline', setFalse);
    };
  }, []);

  return online;
}
