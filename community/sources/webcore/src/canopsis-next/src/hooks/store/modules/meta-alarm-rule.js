import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the meta alarm rule Vuex store module.
 * Provides access to getters and actions for managing meta alarm rule operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The meta alarm rule module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const useMetaAlarmRuleStoreModule = () => useStoreModuleHooks('metaAlarmRule');

/**
 * Custom hook for meta alarm rule operations.
 * Provides convenient access to meta alarm rule getters and actions.
 *
 * @returns {Object} An object containing meta alarm rule getters and actions:
 * @property {Object} metaAlarmRulesMeta - Getter for meta alarm rule metadata
 * @property {boolean} metaAlarmRulesPending - Getter for meta alarm rule loading state
 * @property {Array} metaAlarmRules - Getter for meta alarm rule items
 * @property {Function} fetchMetaAlarmRulesList - Action to fetch meta alarm rules list
 * @property {Function} createMetaAlarmRule - Action to create a meta alarm rule
 * @property {Function} updateMetaAlarmRule - Action to update a meta alarm rule
 * @property {Function} removeMetaAlarmRule - Action to remove a meta alarm rule
 * @property {Function} fetchMetaAlarmRulesListWithoutStore - Action to fetch meta alarm rules list without store
 * @property {Function} bulkEnableMetaAlarmRules - Action to bulk enable meta alarm rules
 * @property {Function} bulkDisableMetaAlarmRules - Action to bulk disable meta alarm rules
 * @property {Function} bulkRemoveMetaAlarmRules - Action to bulk remove meta alarm rules
 *
 * @example
 * // Usage in a component
 * const { metaAlarmRules, metaAlarmRulesPending, fetchMetaAlarmRulesList } = useMetaAlarmRule();
 * await fetchMetaAlarmRulesList({ params: { page: 1, limit: 10 } });
 */
export const useMetaAlarmRule = () => {
  const { useGetters, useActions } = useMetaAlarmRuleStoreModule();

  const getters = useGetters({
    metaAlarmRulesMeta: 'meta',
    metaAlarmRulesPending: 'pending',
    metaAlarmRules: 'items',
  });

  const actions = useActions({
    fetchMetaAlarmRulesList: 'fetchList',
    createMetaAlarmRule: 'create',
    updateMetaAlarmRule: 'update',
    removeMetaAlarmRule: 'remove',
    fetchMetaAlarmRulesListWithoutStore: 'fetchListWithoutStore',
    bulkEnableMetaAlarmRules: 'bulkEnable',
    bulkDisableMetaAlarmRules: 'bulkDisable',
    bulkRemoveMetaAlarmRules: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
