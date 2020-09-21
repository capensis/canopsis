import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('metaAlarmRule');

export default {
  computed: {
    ...mapGetters({
      metaAlarmRulePending: 'pending',
      metaAlarmRule: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchMetaAlarmRuleList: 'fetchList',
      refreshMetaAlarmRuleList: 'fetchListWithPreviousParams',
      createMetaAlarmRule: 'create',
      updateMetaAlarmRule: 'update',
      removeMetaAlarmRule: 'remove',
    }),
  },
};
