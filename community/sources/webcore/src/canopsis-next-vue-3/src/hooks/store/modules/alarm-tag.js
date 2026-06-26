import { storeToRefs } from 'pinia';

import { useAlarmTagStore } from '@/stores/modules/alarm-tag';

/**
 * Alarm tag store hook.
 *
 * @returns {Object}
 */
export const useAlarmTag = () => {
  const alarmTagStore = useAlarmTagStore();
  const {
    items,
    meta,
    pending,
    fetchingParams,
  } = storeToRefs(alarmTagStore);

  const getTagColor = tag => (
    items.value
      .find(({ value }) => tag === value)
      ?.color
  );

  return {
    alarmTags: items,
    alarmTagsMeta: meta,
    alarmTagsPending: pending,
    alarmTagsFetchingParams: fetchingParams,

    fetchAlarmTagsList: alarmTagStore.fetchList,
    fetchAlarmTagsListWithPreviousParams: alarmTagStore.fetchListWithPreviousParams,
    fetchAlarmTagsListWithoutStore: alarmTagStore.fetchListWithoutStore,
    createAlarmTag: alarmTagStore.create,
    updateAlarmTag: alarmTagStore.update,
    removeAlarmTag: alarmTagStore.remove,
    bulkRemoveAlarmTag: alarmTagStore.bulkRemove,

    getTagColor,
  };
};
