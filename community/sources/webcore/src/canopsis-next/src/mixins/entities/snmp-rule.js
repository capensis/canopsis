import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('snmpRule');

export default {
  computed: {
    ...mapGetters({
      snmpRules: 'items',
      snmpRulesMeta: 'meta',
      snmpRulesPending: 'pending',
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
