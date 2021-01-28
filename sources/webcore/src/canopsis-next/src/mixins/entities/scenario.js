import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('scenario');

export default {
  computed: {
    ...mapGetters({
      scenariosPending: 'pending',
      scenarios: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchScenariosList: 'fetchList',
      refreshScenariosList: 'fetchListWithPreviousParams',
      createScenario: 'create',
      updateScenario: 'update',
      removeScenario: 'remove',
    }),
  },
};
