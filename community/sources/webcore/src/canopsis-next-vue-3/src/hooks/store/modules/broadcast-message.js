import { storeToRefs } from 'pinia';

import { useBroadcastMessageStore } from '@/stores/modules/broadcast-message';

/**
 * Broadcast message store facade.
 *
 * Keeps direct Pinia usage out of components and exposes migrated Vuex-like
 * action names with entity-specific aliases.
 *
 * @returns {Object}
 */
export const useBroadcastMessages = () => {
  const broadcastMessageStore = useBroadcastMessageStore();
  const {
    items,
    meta,
    pending,
    fetchingParams,
  } = storeToRefs(broadcastMessageStore);

  return {
    broadcastMessages: items,
    broadcastMessagesMeta: meta,
    broadcastMessagesPending: pending,
    broadcastMessagesFetchingParams: fetchingParams,

    fetchBroadcastMessagesList: broadcastMessageStore.fetchList,
    fetchBroadcastMessagesListWithPreviousParams: broadcastMessageStore.fetchListWithPreviousParams,
    fetchBroadcastMessagesListWithoutStore: broadcastMessageStore.fetchListWithoutStore,
    fetchActiveBroadcastMessagesListWithoutStore: broadcastMessageStore.fetchActiveListWithoutStore,
    createBroadcastMessage: broadcastMessageStore.create,
    updateBroadcastMessage: broadcastMessageStore.update,
    removeBroadcastMessage: broadcastMessageStore.remove,
    markBroadcastMessageAsRead: broadcastMessageStore.markAsRead,
  };
};
