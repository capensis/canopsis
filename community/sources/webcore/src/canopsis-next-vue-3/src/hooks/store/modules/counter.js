import { storeToRefs } from 'pinia';

import { useCounterStore } from '@/stores/modules/counter';

/**
 * Counter widget store facade.
 *
 * @returns {Object}
 */
export const useCounter = () => {
  const counterStore = useCounterStore();
  const { widgets, getListByWidgetId, getPendingByWidgetId } = storeToRefs(counterStore);

  return {
    countersByWidget: widgets,
    getCountersListByWidgetId: getListByWidgetId,
    getCountersPendingByWidgetId: getPendingByWidgetId,
    fetchCountersList: counterStore.fetchList,
  };
};
