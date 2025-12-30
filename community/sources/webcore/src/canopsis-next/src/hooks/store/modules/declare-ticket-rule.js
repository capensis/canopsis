import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the declare ticket rule Vuex store module.
 * Provides access to getters and actions for managing declare ticket rule operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The declare ticket rule module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const useDeclareTicketRuleStoreModule = () => useStoreModuleHooks('declareTicketRule');

/**
 * Custom hook for declare ticket rule operations.
 * Provides convenient access to declare ticket rule getters and actions.
 *
 * @returns {Object} An object containing declare ticket rule getters and actions:
 * @property {Object} declareTicketRulesMeta - Getter for declare ticket rule metadata
 * @property {boolean} declareTicketRulesPending - Getter for declare ticket rule loading state
 * @property {Array} declareTicketRules - Getter for declare ticket rule items
 * @property {Function} fetchDeclareTicketRulesList - Fetches declare ticket rules list
 * @property {Function} createDeclareTicketRule - Creates a declare ticket rule
 * @property {Function} updateDeclareTicketRule - Updates a declare ticket rule
 * @property {Function} removeDeclareTicketRule - Removes a declare ticket rule
 * @property {Function} bulkRemoveDeclareTicketRules - Bulk removes declare ticket rules
 * @property {Function} createTestDeclareTicketExecution - Creates a test declare ticket execution
 * @property {Function} fetchDeclareTicketExecutionWithoutStore
 *                      - Fetches declare ticket execution without store
 * @property {Function} fetchTestDeclareTicketExecutionWebhooksResponse
 *                      - Fetches test declare ticket execution webhooks response
 * @property {Function} bulkCreateDeclareTicketExecution - Bulk creates declare ticket execution
 * @property {Function} fetchAssignedTicketsWithoutStore - Fetches assigned tickets without store
 * @property {Function} bulkEnableDeclareTicketRules - Bulk enables declare ticket rules
 * @property {Function} bulkDisableDeclareTicketRules - Bulk disables declare ticket rules
 * @property {Function} fetchDeclareTicketRulesListWithPreviousParams
 *                      - Fetches declare ticket rules list with previous params
 * @property {Function} fetchDeclareTicketRulesListWithoutStore
 *                      - Fetches declare ticket rules list without store
 *
 * @example
 * // Usage in a component
 * const { declareTicketRules, declareTicketRulesPending, fetchDeclareTicketRulesList } = useDeclareTicketRule();
 * await fetchDeclareTicketRulesList({ params: { page: 1, limit: 10 } });
 */
export const useDeclareTicketRule = () => {
  const { useGetters, useActions } = useDeclareTicketRuleStoreModule();

  const getters = useGetters({
    declareTicketRulesMeta: 'meta',
    declareTicketRulesPending: 'pending',
    declareTicketRules: 'items',
  });

  const actions = useActions({
    fetchDeclareTicketRulesList: 'fetchList',
    fetchDeclareTicketRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchDeclareTicketRulesListWithoutStore: 'fetchListWithoutStore',
    createDeclareTicketRule: 'create',
    updateDeclareTicketRule: 'update',
    removeDeclareTicketRule: 'remove',
    bulkRemoveDeclareTicketRules: 'bulkRemove',
    createTestDeclareTicketExecution: 'createTestDeclareTicketExecution',
    fetchDeclareTicketExecutionWithoutStore: 'fetchDeclareTicketExecutionWithoutStore',
    fetchTestDeclareTicketExecutionWebhooksResponse: 'fetchTestDeclareTicketExecutionWebhooksResponse',
    bulkCreateDeclareTicketExecution: 'bulkCreateDeclareTicketExecution',
    fetchAssignedTicketsWithoutStore: 'fetchAssignedTicketsWithoutStore',
    bulkEnableDeclareTicketRules: 'bulkEnable',
    bulkDisableDeclareTicketRules: 'bulkDisable',
  });

  return {
    ...getters,
    ...actions,
  };
};
