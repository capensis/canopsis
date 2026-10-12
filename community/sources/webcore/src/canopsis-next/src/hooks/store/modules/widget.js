import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Custom hook to access the `view/widget` store module.
 *
 * This hook utilizes the `useStoreModuleHooks` function to provide access to
 * the `view/widget` module within the Vuex store. It allows for easy retrieval
 * of getters and dispatching of actions specific to the `view/widget` namespace.
 *
 * @returns {Object} An object containing the store, module, useGetters, and useActions
 *                   functions for the `view/widget` namespace.
 */
export const useWidgetStoreModule = () => useStoreModuleHooks('view/widget');

/**
 * Custom hook to interact with the `view/widget` store module.
 *
 * This hook provides access to the `view/widget` module's getters and actions,
 * allowing for easy retrieval and manipulation of widget data.
 *
 * @returns {Object} An object containing:
 * - `fetchWidgetWithoutStore`: An action to fetch a widget without storing it.
 * - `createWidget`: An action to create a new widget.
 * - `updateWidget`: An action to update an existing widget.
 * - `removeWidget`: An action to remove a widget.
 * - `updateWidgetGridPositions`: An action to update widget grid positions.
 * - `fetchWidgetFilters`: An action to fetch widget filters.
 * - `fetchWidgetFilter`: An action to fetch a specific widget filter.
 * - `createWidgetFilter`: An action to create a new widget filter.
 * - `updateWidgetFilter`: An action to update an existing widget filter.
 * - `removeWidgetFilter`: An action to remove a widget filter.
 * - `updateWidgetFiltersPositions`: An action to update widget filter positions.
 */
export const useWidget = () => {
  const { useActions } = useWidgetStoreModule();

  const actions = useActions({
    fetchWidgetWithoutStore: 'fetchItemWithoutStore',
    createWidget: 'create',
    updateWidget: 'update',
    removeWidget: 'remove',
    updateWidgetGridPositions: 'updateGridPositions',
    fetchWidgetFilters: 'fetchWidgetFilters',
    fetchWidgetFilter: 'fetchWidgetFilter',
    createWidgetFilter: 'createWidgetFilter',
    updateWidgetFilter: 'updateWidgetFilter',
    removeWidgetFilter: 'removeWidgetFilter',
    updateWidgetFiltersPositions: 'updateWidgetFiltersPositions',
  });

  return {
    ...actions,
  };
};
