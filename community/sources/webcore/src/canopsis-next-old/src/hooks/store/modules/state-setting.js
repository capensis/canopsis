import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the state setting Vuex store module.
 *
 * @returns {Object} An object containing store module utilities
 */
const useStateSettingStoreModule = () => useStoreModuleHooks('stateSetting');

/**
 * Hook for state setting operations and state access.
 *
 * @returns {Object} State setting getters and actions
 */
export const useStateSetting = () => {
  const { useGetters, useActions } = useStateSettingStoreModule();

  const getters = useGetters({
    stateSettings: 'items',
    stateSettingsPending: 'pending',
    stateSettingsMeta: 'meta',
  });

  const actions = useActions({
    fetchStateSettingsList: 'fetchList',
    fetchStateSettingsListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchStateSettingsListWithoutStore: 'fetchListWithoutStore',
    createStateSetting: 'create',
    updateStateSetting: 'update',
    removeStateSetting: 'remove',
  });

  return {
    ...getters,
    ...actions,
  };
};
