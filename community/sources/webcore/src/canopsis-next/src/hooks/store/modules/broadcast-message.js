import { useStoreModuleHooks } from '@/hooks/store';

const useBroadcastMessageStoreModule = () => useStoreModuleHooks('broadcastMessage');

/**
 * Custom hook for accessing actions related to broadcast messages from Vuex store module.
 *
 * @returns {Object} An object containing functions to create, update, remove, and fetch broadcast messages list.
 *
 * @example
 * // Usage example within a Vue component
 * import { useBroadcastMessages } from './path/to/useBroadcastMessages';
 *
 * export default {
 *   setup() {
 *     const {
 *       createBroadcastMessage,
 *       updateBroadcastMessage,
 *       removeBroadcastMessage,
 *       fetchBroadcastMessagesListWithoutStore
 *     } = useBroadcastMessages();
 *
 *     // Access and use the functions as needed
 *
 *     return {
 *       createBroadcastMessage,
 *       updateBroadcastMessage,
 *       removeBroadcastMessage,
 *       fetchBroadcastMessagesListWithoutStore
 *     };
 *   }
 * }
 */
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
