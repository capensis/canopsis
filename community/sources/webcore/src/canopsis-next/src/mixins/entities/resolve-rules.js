import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('resolveRules');

export const entitiesResolveRulesMixin = {
  computed: {
    ...mapGetters({
      resolveRulesMeta: 'meta',
      resolveRulesPending: 'pending',
      resolveRules: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchResolveRulesList: 'fetchList',
      createResolveRule: 'create',
      updateResolveRule: 'update',
      removeResolveRule: 'remove',
    }),
  },
};
