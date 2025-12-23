import { useStoreModuleHooks } from '@/hooks/store';

const useDeclareTicketRuleStoreModule = () => useStoreModuleHooks('declareTicketRule');

/**
 * Hook for accessing declare ticket rule store module actions
 *
 * @returns {Object} An object containing declare ticket rule actions
 */
export const useDeclareTicketRule = () => {
  const { useActions } = useDeclareTicketRuleStoreModule();

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
    ...actions,
  };
};
