import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the link rule Vuex store module.
 * Provides access to getters and actions for managing link rule operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The link rule module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const useLinkRuleStoreModule = () => useStoreModuleHooks('linkRule');

/**
 * Custom hook for link rule operations.
 * Provides convenient access to link rule getters and actions.
 *
 * @returns {Object} An object containing link rule getters and actions:
 * @property {Object} linkRulesMeta - Getter for link rule metadata
 * @property {boolean} linkRulesPending - Getter for link rule loading state
 * @property {Array} linkRules - Getter for link rule items
 * @property {Function} fetchLinkRulesList - Action to fetch link rules list
 * @property {Function} createLinkRule - Action to create a link rule
 * @property {Function} updateLinkRule - Action to update a link rule
 * @property {Function} removeLinkRule - Action to remove a link rule
 * @property {Function} bulkRemoveLinkRules - Action to bulk remove link rules
 * @property {Function} fetchLinkCategoriesWithoutStore - Action to fetch link categories without store
 * @property {Function} bulkEnableLinkRules - Action to bulk enable link rules
 * @property {Function} bulkDisableLinkRules - Action to bulk disable link rules
 * @property {Function} fetchLinkRulesListWithPreviousParams - Action to fetch link rules list with previous params
 * @property {Function} fetchLinkRulesListWithoutStore - Action to fetch link rules list without store
 *
 * @example
 * // Usage in a component
 * const { linkRules, linkRulesPending, fetchLinkRulesList } = useLinkRule();
 * await fetchLinkRulesList({ params: { page: 1, limit: 10 } });
 */
export const useLinkRule = () => {
  const { useGetters, useActions } = useLinkRuleStoreModule();

  const getters = useGetters({
    linkRulesMeta: 'meta',
    linkRulesPending: 'pending',
    linkRules: 'items',
  });

  const actions = useActions({
    fetchLinkRulesList: 'fetchList',
    fetchLinkRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchLinkRulesListWithoutStore: 'fetchListWithoutStore',
    createLinkRule: 'create',
    updateLinkRule: 'update',
    removeLinkRule: 'remove',
    bulkRemoveLinkRules: 'bulkRemove',
    fetchLinkCategoriesWithoutStore: 'fetchLinkCategoriesWithoutStore',
    bulkEnableLinkRules: 'bulkEnable',
    bulkDisableLinkRules: 'bulkDisable',
  });

  return {
    ...getters,
    ...actions,
  };
};
