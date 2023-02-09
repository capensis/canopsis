import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('declareTicketRule');

export const entitiesDeclareTicketRuleMixin = {
  computed: {
    ...mapGetters({
      declareTicketRulesMeta: 'meta',
      declareTicketRulesPending: 'pending',
      declareTicketRules: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchDeclareTicketRulesList: 'fetchList',
      fetchAssignedDeclareTicketsWithoutStore: 'fetchAssignedTicketsWithoutStore',
      createDeclareTicketRule: 'create',
      updateDeclareTicketRule: 'update',
      removeDeclareTicketRule: 'remove',
      createTestDeclareTicketExecution: 'createTestDeclareTicketExecution',
      fetchDeclareTicketExecutionWithoutStore: 'fetchDeclareTicketExecutionWithoutStore',
      fetchTestDeclareTicketExecutionWebhooksResponse: 'fetchTestDeclareTicketExecutionWebhooksResponse',
      bulkRemoveDeclareTicketRules: 'bulkRemove',
      bulkCreateDeclareTicketExecution: 'bulkCreateDeclareTicketExecution',
    }),
  },
};
