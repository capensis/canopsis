import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('eventFilterRule');

export default {
  computed: {
    ...mapGetters({
      eventFilterRulesPending: 'pending',
      eventFilterRules: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchEventFilterRulesList: 'fetchList',
      refreshEventFilterRulesList: 'fetchListWithPreviousParams',
      createEventFilterRule: 'create',
      updateEventFilterRule: 'update',
      removeEventFilterRule: 'remove',
    }),
  },
};
