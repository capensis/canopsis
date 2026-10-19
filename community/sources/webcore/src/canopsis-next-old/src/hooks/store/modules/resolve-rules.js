import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the resolve rules Vuex store module.
 * Provides access to getters and actions for managing resolve rules operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The resolve rules module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const useResolveRulesStoreModule = () => useStoreModuleHooks('resolveRules');

/**
 * Custom hook for resolve rules operations.
 * Provides convenient access to resolve rules getters and actions.
 *
 * @returns {Object} An object containing resolve rules getters and actions:
 * @property {Object} resolveRulesMeta - Getter for resolve rules metadata
 * @property {boolean} resolveRulesPending - Getter for resolve rules loading state
 * @property {Array} resolveRules - Getter for resolve rules items
 * @property {Function} fetchResolveRulesList - Action to fetch resolve rules list
 * @property {Function} createResolveRule - Action to create a resolve rule
 * @property {Function} updateResolveRule - Action to update a resolve rule
 * @property {Function} removeResolveRule - Action to remove a resolve rule
 * @property {Function} bulkEnableResolveRules - Action to bulk enable resolve rules
 * @property {Function} bulkDisableResolveRules - Action to bulk disable resolve rules
 * @property {Function} bulkRemoveResolveRules - Action to bulk remove resolve rules
 * @property {Function} fetchResolveRulesListWithPreviousParams
 *                      - Action to fetch resolve rules list with previous params
 * @property {Function} fetchResolveRulesListWithoutStore
 *                      - Action to fetch resolve rules list without store
 *
 * @example
 * // Usage in a component
 * const { resolveRules, resolveRulesPending, fetchResolveRulesList } = useResolveRules();
 * await fetchResolveRulesList({ params: { page: 1, limit: 10 } });
 */
export const useResolveRules = () => {
  const { useGetters, useActions } = useResolveRulesStoreModule();

  const getters = useGetters({
    resolveRulesMeta: 'meta',
    resolveRulesPending: 'pending',
    resolveRules: 'items',
  });

  const actions = useActions({
    fetchResolveRulesList: 'fetchList',
    fetchResolveRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchResolveRulesListWithoutStore: 'fetchListWithoutStore',
    createResolveRule: 'create',
    updateResolveRule: 'update',
    removeResolveRule: 'remove',
    bulkEnableResolveRules: 'bulkEnable',
    bulkDisableResolveRules: 'bulkDisable',
    bulkRemoveResolveRules: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
