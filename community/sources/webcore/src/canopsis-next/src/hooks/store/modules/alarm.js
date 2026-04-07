import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Custom hook to access the `alarm` store module.
 *
 * This hook utilizes the `useStoreModuleHooks` function to provide access to
 * the `alarm` module within the Vuex store. It allows for easy retrieval
 * of getters and dispatching of actions specific to the `alarm` namespace.
 *
 * @returns {Object} An object containing the store, module, useGetters, and useActions
 *                   functions for the `alarm` namespace.
 */
export const useAlarmStoreModule = () => useStoreModuleHooks('alarm');

/**
 * Custom hook to interact with the `alarm` store module.
 *
 * This hook provides access to the `alarm` module's getters and actions,
 * allowing for easy retrieval and manipulation of alarm data.
 *
 * @returns {Object} An object containing:
 * - `getAlarmItem`: A getter for retrieving a specific alarm item.
 * - `getAlarmsList`: A getter for retrieving a list of alarms.
 * - `getAlarmsListByWidgetId`: A getter for retrieving alarms by widget ID.
 * - `getAlarmsMetaByWidgetId`: A getter for retrieving alarm metadata by widget ID.
 * - `getAlarmsPendingByWidgetId`: A getter for checking if alarms are being fetched by widget ID.
 * - `getAlarmsFetchingParamsByWidgetId`: A getter for retrieving fetching parameters by widget ID.
 * - `fetchAlarmItem`: An action to fetch a specific alarm item.
 * - `fetchAlarmsList`: An action to fetch a list of alarms.
 * - `fetchAlarmsListWithoutStore`: An action to fetch a list of alarms without storing in the store.
 * - `createAlarmsListExport`: An action to create an alarm list export.
 * - `fetchAlarmsListExport`: An action to fetch an alarm list export.
 */
export const useAlarm = () => {
  const { useGetters, useActions } = useAlarmStoreModule();

  const getters = useGetters({
    getAlarmItem: 'getItem',
    getAlarmsList: 'getList',
    getAlarmsListByWidgetId: 'getListByWidgetId',
    getAlarmsMetaByWidgetId: 'getMetaByWidgetId',
    getAlarmsPendingByWidgetId: 'getPendingByWidgetId',
    getAlarmsFetchingParamsByWidgetId: 'getFetchingParamsByWidgetId',
  });

  const actions = useActions({
    fetchAlarmItem: 'fetchItem',
    fetchAlarmsList: 'fetchList',
    fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    createAlarmsListExport: 'createAlarmsListExport',
    fetchAlarmsListExport: 'fetchAlarmsListExport',
  });

  return {
    ...getters,
    ...actions,
  };
};
