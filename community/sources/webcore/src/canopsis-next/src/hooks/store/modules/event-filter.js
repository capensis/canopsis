import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Event Filter Store Module.
 *
 * @returns {Object} An object containing getters and actions for event filters.
 */
const useEventFilterStoreModule = () => useStoreModuleHooks('eventFilter');

/**
 * Hook to access event filter store.
 *
 * @returns {Object} An object containing:
 * - Actions to fetch event filter rules and manage them.
 */
export const useEventFilterStore = () => {
  const { useGetters, useActions } = useEventFilterStoreModule();

  const getters = useGetters({
    eventFiltersPending: 'pending',
    eventFilters: 'items',
    eventFiltersMeta: 'meta',
  });

  const actions = useActions({
    fetchEventFiltersList: 'fetchList',
    fetchEventFiltersListWithoutStore: 'fetchListWithoutStore',
    refreshEventFiltersList: 'fetchListWithPreviousParams',
    fetchEventFilterFailuresListWithoutStore: 'fetchEventFilterFailuresListWithoutStore',
    markNewEventFilterFailuresAsRead: 'markNewEventFilterFailuresAsRead',
    createEventFilter: 'create',
    updateEventFilter: 'update',
    removeEventFilter: 'remove',
  });

  return {
    ...getters,
    ...actions,
  };
};
