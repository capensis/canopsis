import { useMetaAlarmStore } from '@/stores/modules/meta-alarm';

/**
 * Meta alarm store facade.
 *
 * @returns {Object}
 */
export const useMetaAlarm = () => {
  const metaAlarmStore = useMetaAlarmStore();

  return {
    fetchMetaAlarmsListWithoutStore: metaAlarmStore.fetchListWithoutStore,
    createMetaAlarm: metaAlarmStore.create,
    addMetaAlarmAlarms: metaAlarmStore.addAlarms,
    removeMetaAlarmAlarms: metaAlarmStore.removeAlarms,
  };
};
