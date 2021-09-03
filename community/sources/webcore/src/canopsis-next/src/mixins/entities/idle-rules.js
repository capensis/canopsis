import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('idleRules');

export const entitiesIdleRulesMixin = {
  computed: {
    ...mapGetters({
      idleRulesMeta: 'meta',
      idleRulesPending: 'pending',
      idleRules: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchIdleRulesList: 'fetchList',
      createIdleRule: 'create',
      updateIdleRule: 'update',
      removeIdleRule: 'remove',
    }),
  },
};
