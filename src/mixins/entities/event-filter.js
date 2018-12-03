import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('eventFilterRule');

export default {
  computed: {
    ...mapGetters({
      getEventFilterRulePending: 'pending',
      items: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchEventFilterRulesList: 'fetchList',
      removeEventFilterRule: 'remove',
    }),
  },
};
