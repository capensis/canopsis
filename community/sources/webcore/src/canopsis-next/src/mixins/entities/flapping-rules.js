import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('flappingRules');

export const entitiesFlappingRulesMixin = {
  computed: {
    ...mapGetters({
      flappingRulesMeta: 'meta',
      flappingRulesPending: 'pending',
      flappingRules: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchFlappingRulesList: 'fetchList',
      createFlappingRule: 'create',
      updateFlappingRule: 'update',
      removeFlappingRule: 'remove',
    }),
  },
};
