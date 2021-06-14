import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('metaAlarmRule');

export const entitiesMetaAlarmRuleMixin = {
  computed: {
    ...mapGetters({
      metaAlarmRulesPending: 'pending',
      metaAlarmRules: 'items',
      metaAlarmRulesMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchMetaAlarmRulesList: 'fetchList',
      createMetaAlarmRule: 'create',
      updateMetaAlarmRule: 'update',
      removeMetaAlarmRule: 'remove',
    }),
  },
};
