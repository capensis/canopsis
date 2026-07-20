import { storeToRefs } from 'pinia'

import { useDeclareTicketRuleStore } from '@/stores/modules/declare-ticket-rule'

/**
 * Declare ticket rule store facade.
 *
 * @returns {Object}
 */
export const useDeclareTicketRule = () => {
  const declareTicketRuleStore = useDeclareTicketRuleStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(declareTicketRuleStore)

  return {
    declareTicketRules: items,
    declareTicketRulesMeta: meta,
    declareTicketRulesPending: pending,
    declareTicketRulesFetchingParams: fetchingParams,

    fetchDeclareTicketRulesList: declareTicketRuleStore.fetchList,
    fetchDeclareTicketRulesListWithPreviousParams:
      declareTicketRuleStore.fetchListWithPreviousParams,
    fetchDeclareTicketRulesListWithoutStore: declareTicketRuleStore.fetchListWithoutStore,
    createDeclareTicketRule: declareTicketRuleStore.create,
    updateDeclareTicketRule: declareTicketRuleStore.update,
    removeDeclareTicketRule: declareTicketRuleStore.remove,
    bulkRemoveDeclareTicketRules: declareTicketRuleStore.bulkRemove,
    bulkEnableDeclareTicketRules: declareTicketRuleStore.bulkEnable,
    bulkDisableDeclareTicketRules: declareTicketRuleStore.bulkDisable,
    createTestDeclareTicketExecution: declareTicketRuleStore.createTestDeclareTicketExecution,
    fetchDeclareTicketExecutionWithoutStore:
      declareTicketRuleStore.fetchDeclareTicketExecutionWithoutStore,
    fetchTestDeclareTicketExecutionWebhooksResponse:
      declareTicketRuleStore.fetchTestDeclareTicketExecutionWebhooksResponse,
    bulkCreateDeclareTicketExecution: declareTicketRuleStore.bulkCreateDeclareTicketExecution,
    fetchAssignedTicketsWithoutStore: declareTicketRuleStore.fetchAssignedTicketsWithoutStore,
  }
}
