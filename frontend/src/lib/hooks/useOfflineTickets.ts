import { useCallback, useEffect, useState } from 'react';
import { savePendingTicket, getPendingTickets, type PendingTicket } from '../offline/db';
import { registerBackgroundSync, syncPendingTickets } from '../offline/sync';

export function useOfflineTickets() {
  const [pendingCount, setPendingCount] = useState(0);

  const refresh = useCallback(async () => {
    const pending = await getPendingTickets();
    setPendingCount(pending.length);
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const save = useCallback(
    async (data: PendingTicket['data']) => {
      const id = await savePendingTicket(data);
      await registerBackgroundSync();
      await refresh();
      return id;
    },
    [refresh],
  );

  const sync = useCallback(async () => {
    const synced = await syncPendingTickets();
    await refresh();
    return synced;
  }, [refresh]);

  return { pendingCount, save, sync, refresh };
}
