import { createNamespacedHelpers } from 'vuex';

import popupMixin from '@/mixins/popup';

const { mapActions, mapGetters } = createNamespacedHelpers('eventFilterRule');

export default {
  mixins: [popupMixin],
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
