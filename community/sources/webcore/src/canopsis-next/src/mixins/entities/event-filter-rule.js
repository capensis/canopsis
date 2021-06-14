import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('eventFilterRule');

export const entitiesEventFilterRuleMixin = {
  computed: {
    ...mapGetters({
      eventFilterRulesPending: 'pending',
      eventFilterRules: 'items',
      eventFilterRulesMeta: 'meta',
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
