import { storeToRefs } from 'pinia';

import { useEventFilterStore } from '@/stores/modules/event-filter';

/**
 * Event filters store facade.
 *
 * @returns {Object}
 */
export const useEventFilter = () => {
  const eventFilterStore = useEventFilterStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(eventFilterStore);

  return {
    eventFilters: items,
    eventFiltersMeta: meta,
    eventFiltersPending: pending,
    eventFiltersFetchingParams: fetchingParams,

    fetchEventFiltersList: eventFilterStore.fetchList,
    fetchEventFiltersListWithPreviousParams: eventFilterStore.fetchListWithPreviousParams,
    fetchEventFiltersListWithoutStore: eventFilterStore.fetchListWithoutStore,
    fetchEventFilterFailuresListWithoutStore:
      eventFilterStore.fetchEventFilterFailuresListWithoutStore,
    markNewEventFilterFailuresAsRead: eventFilterStore.markNewEventFilterFailuresAsRead,
    createEventFilter: eventFilterStore.create,
    updateEventFilter: eventFilterStore.update,
    removeEventFilter: eventFilterStore.remove,
    bulkEnableEventFilters: eventFilterStore.bulkEnable,
    bulkDisableEventFilters: eventFilterStore.bulkDisable,
    bulkRemoveEventFilters: eventFilterStore.bulkRemove,
  };
};
