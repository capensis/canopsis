import { useAlarmTagLabelStore } from '@/stores/modules/alarm-tag-label';

/**
 * Alarm tag label store hook.
 *
 * @returns {Object}
 */
export const useAlarmTagLabel = () => {
  const alarmTagLabelStore = useAlarmTagLabelStore();

  return {
    fetchAlarmTagsLabelsListWithoutStore: alarmTagLabelStore.fetchListWithoutStore,
  };
};
