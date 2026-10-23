import { useCallback, useEffect, useState } from 'react';
import { savePendingTicket, getPendingTickets, type OfflineIdentity, type PendingTicket } from '../offline/db';
import { registerBackgroundSync, syncPendingTickets } from '../offline/sync';

export function useOfflineTickets(identity: OfflineIdentity) {
  const [pendingCount, setPendingCount] = useState(0);

  const refresh = useCallback(async () => {
    const pending = await getPendingTickets(identity);
    setPendingCount(pending.length);
  }, [identity]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const save = useCallback(
    async (data: PendingTicket['data']) => {
      const id = await savePendingTicket(data, identity);
      await registerBackgroundSync();
      await refresh();
      return id;
    },
    [refresh, identity],
  );

  const sync = useCallback(async () => {
    const synced = await syncPendingTickets(identity);
    await refresh();
    return synced;
  }, [refresh, identity]);

  return { pendingCount, save, sync, refresh };
}
