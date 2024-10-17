import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to access the Vuex store module for template variables.
 * This function utilizes the `useStoreModuleHooks` to provide a modular way to interact with
 * the `templateVars` Vuex module.
 *
 * @returns {Object} An object containing the store, module, and functions to access getters and actions for
 * the `templateVars` module.
 */
const useTemplateVarsStoreModule = () => useStoreModuleHooks('templateVars');

/**
 * Custom hook for accessing and interacting with the `templateVars` Vuex store module.
 * This hook provides a convenient way to access the module's getters and actions.
 *
 * @returns {Object} An object containing the mapped getters and actions for the `templateVars` module.
 *
 * @property {Function} templateVars - Getter for accessing the `templateVars` items.
 * @property {Function} templateVarsPending - Getter for checking if the `templateVars` data is pending.
 * @property {Function} fetchTemplateVarsList - Action for fetching the list of `templateVars`.
 */
export const useTemplateVars = () => {
  const { useGetters, useActions } = useTemplateVarsStoreModule();

  const getters = useGetters({
    templateVars: 'items',
    templateVarsPending: 'pending',
  });

  const actions = useActions({
    fetchTemplateVarsList: 'fetchList',
  });

  return {
    ...getters,
    ...actions,
  };
};
