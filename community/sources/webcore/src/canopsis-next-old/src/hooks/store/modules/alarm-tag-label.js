import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Custom hook to access the `alarmTag/label` store module.
 *
 * This hook utilizes the `useStoreModuleHooks` function to provide access to
 * the `alarmTag` module within the Vuex store. It allows for easy retrieval
 * of getters and dispatching of actions specific to the `alarmTag/label` namespace.
 *
 * @returns {Object} An object containing the store, module, useGetters, and useActions
 *                   functions for the `alarmTag/label` namespace.
 */
export const useAlarmTagLabelStoreModule = () => useStoreModuleHooks('alarmTag/label');

/**
 * Custom hook to interact with the `alarmTag/label` store module.
 *
 * This hook provides access to the `alarmTag/label` module's getters and actions,
 * allowing for easy retrieval and manipulation of alarm tag data.
 *
 * @returns {Object} An object containing:
 * - `fetchAlarmTagsLabelsListWithoutStore`: An action to fetch the list of alarm tags without store.
 */
export const useAlarmTagLabel = () => {
  const { useActions } = useAlarmTagLabelStoreModule();

  const actions = useActions({
    fetchAlarmTagsLabelsListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...actions,
  };
};
