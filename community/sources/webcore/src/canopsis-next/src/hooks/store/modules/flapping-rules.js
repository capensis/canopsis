import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the flapping rules Vuex store module.
 * Provides access to getters and actions for managing flapping rules operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The flapping rules module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const useFlappingRulesStoreModule = () => useStoreModuleHooks('flappingRules');

/**
 * Hook for accessing flapping rules store module actions
 * Provides convenient access to flapping rules getters and actions.
 *
 * @returns {Object} An object containing flapping rules getters and actions:
 * @property {Object} flappingRulesMeta - Getter for flapping rules metadata
 * @property {boolean} flappingRulesPending - Getter for flapping rules loading state
 * @property {Array} flappingRules - Getter for flapping rules items
 * @property {Function} fetchFlappingRulesList - Action to fetch flapping rules list
 * @property {Function} createFlappingRule - Action to create a flapping rule
 * @property {Function} updateFlappingRule - Action to update a flapping rule
 * @property {Function} removeFlappingRule - Action to remove a flapping rule
 * @property {Function} bulkEnableFlappingRules - Action to bulk enable flapping rules
 * @property {Function} bulkDisableFlappingRules - Action to bulk disable flapping rules
 * @property {Function} bulkRemoveFlappingRules - Action to bulk remove flapping rules
 * @property {Function} fetchFlappingRulesListWithPreviousParams - Fetches flapping rules list with previous params
 * @property {Function} fetchFlappingRulesListWithoutStore - Fetches flapping rules list without store
 *
 * @example
 * // Usage in a component
 * const { flappingRules, flappingRulesPending, fetchFlappingRulesList } = useFlappingRules();
 * await fetchFlappingRulesList({ params: { page: 1, limit: 10 } });
 */
export const useFlappingRules = () => {
  const { useGetters, useActions } = useFlappingRulesStoreModule();

  const getters = useGetters({
    flappingRulesMeta: 'meta',
    flappingRulesPending: 'pending',
    flappingRules: 'items',
  });

  const actions = useActions({
    fetchFlappingRulesList: 'fetchList',
    fetchFlappingRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchFlappingRulesListWithoutStore: 'fetchListWithoutStore',
    createFlappingRule: 'create',
    updateFlappingRule: 'update',
    removeFlappingRule: 'remove',
    bulkEnableFlappingRules: 'bulkEnable',
    bulkDisableFlappingRules: 'bulkDisable',
    bulkRemoveFlappingRules: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
