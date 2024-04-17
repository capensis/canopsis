import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('scenario');

export const entitiesScenarioMixin = {
  computed: {
    ...mapGetters({
      scenariosMeta: 'meta',
      scenariosPending: 'pending',
      scenarios: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchScenariosList: 'fetchList',
      createScenario: 'create',
      updateScenario: 'update',
      removeScenario: 'remove',
      createTestScenarioExecution: 'createTestScenarioExecution',
      fetchTestScenarioExecutionWithoutStore: 'fetchTestScenarioExecutionWithoutStore',
      fetchTestScenarioExecutionWebhooksResponse: 'fetchTestScenarioExecutionWebhooksResponse',
    }),
  },
};
