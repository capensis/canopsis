import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Notifications Store Module.
 *
 * @returns {Object} An object containing getters and actions for notifications.
 */
const useNotificationsStoreModule = () => useStoreModuleHooks('notification');

/**
 * Hook to access notifications store.
 *
 * @returns {Object} An object containing:
 * - Actions to fetch different types of notifications and manage them.
 */
export const useNotifications = () => {
  const { useActions } = useNotificationsStoreModule();

  const actions = useActions({
    fetchNotificationsListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...actions,
  };
};
