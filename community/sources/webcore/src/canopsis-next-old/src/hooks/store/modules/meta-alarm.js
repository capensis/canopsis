import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Retrieves the store module hooks for managing meta alarm.
 *
 * @returns {Object} An object containing store module.
 */
const useMetaAlarmStoreModule = () => useStoreModuleHooks('metaAlarm');

/**
 * Provides actions to interact with meta alarm.
 *
 * @returns {Object} An object containing actions to create, update, remove, and fetch meta alarm list without
 * storing.
 */
export const useMetaAlarm = () => {
  const { useActions } = useMetaAlarmStoreModule();

  const actions = useActions({
    fetchMetaAlarmsListWithoutStore: 'fetchListWithoutStore',
    createMetaAlarm: 'create',
    addAlarmsIntoMetaAlarm: 'addAlarms',
    removeAlarmsFromMetaAlarm: 'removeAlarms',
  });

  return actions;
};
