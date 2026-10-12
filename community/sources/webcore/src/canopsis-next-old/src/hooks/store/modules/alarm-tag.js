import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Custom hook to access the `alarmTag` store module.
 *
 * This hook utilizes the `useStoreModuleHooks` function to provide access to
 * the `alarmTag` module within the Vuex store. It allows for easy retrieval
 * of getters and dispatching of actions specific to the `alarmTag` namespace.
 *
 * @returns {Object} An object containing the store, module, useGetters, and useActions
 *                   functions for the `alarmTag` namespace.
 */
export const useAlarmTagStoreModule = () => useStoreModuleHooks('alarmTag');

/**
 * Custom hook to interact with the `alarmTag` store module.
 *
 * This hook provides access to the `alarmTag` module's getters and actions,
 * allowing for easy retrieval and manipulation of alarm tag data.
 *
 * @returns {Object} An object containing:
 * - `alarmTags`: A getter for retrieving the list of alarm tags.
 * - `alarmTagsPending`: A getter indicating if the alarm tags are currently being fetched.
 * - `alarmTagsMeta`: A getter for retrieving metadata associated with alarm tags.
 * - `fetchAlarmTagsList`: An action to fetch the list of alarm tags.
 * - `fetchAlarmTagsListWithoutStore`: An action to fetch the list of alarm tags without store.
 * - `createAlarmTag`: An action to create a new alarm tag.
 * - `updateAlarmTag`: An action to update an existing alarm tag.
 * - `removeAlarmTag`: An action to remove an alarm tag.
 * - `bulkRemoveAlarmTag`: An action to remove multiple alarm tags.
 * - `getTagColor`: A function to get the color of a specific tag.
 */
export const useAlarmTag = () => {
  const { useGetters, useActions } = useAlarmTagStoreModule();

  const getters = useGetters({
    alarmTags: 'items',
    alarmTagsPending: 'pending',
    alarmTagsMeta: 'meta',
  });

  const actions = useActions({
    fetchAlarmTagsList: 'fetchList',
    fetchAlarmTagsListWithoutStore: 'fetchListWithoutStore',
    createAlarmTag: 'create',
    updateAlarmTag: 'update',
    removeAlarmTag: 'remove',
    bulkRemoveAlarmTag: 'bulkRemove',
  });

  const getTagColor = tag => (
    getters.alarmTags.value
      .find(({ value }) => tag === value)
      ?.color
  );

  return {
    ...getters,
    ...actions,

    getTagColor,
  };
};
