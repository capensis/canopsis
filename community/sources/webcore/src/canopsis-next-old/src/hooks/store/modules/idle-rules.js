import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the idle rules Vuex store module.
 * Provides access to getters and actions for managing idle rules operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The idle rules module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const useIdleRulesStoreModule = () => useStoreModuleHooks('idleRules');

/**
 * Custom hook for idle rules operations.
 * Provides convenient access to idle rules getters and actions.
 *
 * @returns {Object} An object containing idle rules getters and actions:
 * @property {Object} idleRulesMeta - Getter for idle rules metadata
 * @property {boolean} idleRulesPending - Getter for idle rules loading state
 * @property {Array} idleRules - Getter for idle rules items
 * @property {Function} fetchIdleRulesList - Action to fetch idle rules list
 * @property {Function} createIdleRule - Action to create an idle rule
 * @property {Function} updateIdleRule - Action to update an idle rule
 * @property {Function} removeIdleRule - Action to remove an idle rule
 * @property {Function} bulkEnableIdleRules - Action to bulk enable idle rules
 * @property {Function} bulkDisableIdleRules - Action to bulk disable idle rules
 * @property {Function} bulkRemoveIdleRules - Action to bulk remove idle rules
 * @property {Function} fetchIdleRulesListWithPreviousParams - Action to fetch idle rules list with previous params
 * @property {Function} fetchIdleRulesListWithoutStore - Action to fetch idle rules list without store
 *
 * @example
 * // Usage in a component
 * const { idleRules, idleRulesPending, fetchIdleRulesList } = useIdleRules();
 * await fetchIdleRulesList({ params: { page: 1, limit: 10 } });
 */
export const useIdleRules = () => {
  const { useGetters, useActions } = useIdleRulesStoreModule();

  const getters = useGetters({
    idleRulesMeta: 'meta',
    idleRulesPending: 'pending',
    idleRules: 'items',
  });

  const actions = useActions({
    fetchIdleRulesList: 'fetchList',
    fetchIdleRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchIdleRulesListWithoutStore: 'fetchListWithoutStore',
    createIdleRule: 'create',
    updateIdleRule: 'update',
    removeIdleRule: 'remove',
    bulkEnableIdleRules: 'bulkEnable',
    bulkDisableIdleRules: 'bulkDisable',
    bulkRemoveIdleRules: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
