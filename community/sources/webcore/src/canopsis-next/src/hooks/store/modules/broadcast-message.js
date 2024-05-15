import { useStoreModuleHooks } from '@/hooks/store';

const useBroadcastMessageStoreModule = () => useStoreModuleHooks('broadcastMessage');

export const useBroadcastMessages = () => {
  const { useActions } = useBroadcastMessageStoreModule();

  const actions = useActions({
    createBroadcastMessage: 'create',
    updateBroadcastMessage: 'update',
    removeBroadcastMessage: 'remove',
    fetchBroadcastMessagesListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...actions,
  };
};
