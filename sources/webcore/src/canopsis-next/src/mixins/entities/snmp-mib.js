import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('snmpMib');

export default {
  methods: {
    ...mapActions({
      fetchSnmpMibList: 'fetchList',
      fetchSnmpMibDistinctList: 'fetchDistinctList',
    }),
  },
};
