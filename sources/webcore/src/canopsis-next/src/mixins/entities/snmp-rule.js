import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('snmpRule');

export default {
  computed: {
    ...mapGetters({
      snmpRulesPending: 'pending',
      snmpRules: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchSnmpRulesList: 'fetchList',
      createSnmpRule: 'create',
      updateSnmpRule: 'update',
      removeSnmpRule: 'remove',
    }),
  },
};
