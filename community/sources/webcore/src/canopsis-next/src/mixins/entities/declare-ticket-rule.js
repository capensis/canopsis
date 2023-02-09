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
      createDeclareTicketRule: 'create',
      updateDeclareTicketRule: 'update',
      removeDeclareTicketRule: 'remove',
      createTestDeclareTicketExecution: 'createTestDeclareTicketExecution',
      fetchTestDeclareTicketExecutionWebhooksResponse: 'fetchTestDeclareTicketExecutionWebhooksResponse',
      bulkRemoveDeclareTicketRules: 'bulkRemove',
    }),
  },
};
