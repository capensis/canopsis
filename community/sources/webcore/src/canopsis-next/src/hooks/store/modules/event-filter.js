import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Event Filter Store Module.
 *
 * @returns {Object} An object containing getters and actions for event filters.
 */
const useEventFilterStore = () => useStoreModuleHooks('eventFilter');

/**
 * Hook to access event filter store.
 *
 * @returns {Object} An object containing:
 * - Actions to fetch event filter rules and manage them.
 */
export const useEventFilter = () => {
  const { useGetters, useActions } = useEventFilterStore();

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
    bulkEnableEventFilters: 'bulkEnable',
    bulkDisableEventFilters: 'bulkDisable',
    bulkRemoveEventFilters: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
