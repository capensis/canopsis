import { createNamespacedHelpers } from 'vuex';

import popupMixin from '@/mixins/popup';
import { createMapActionsWithPopup, createMapActionsWith } from '@/helpers/store';

const { mapActions, mapGetters } = createNamespacedHelpers('eventFilterRule');

const mapActionsWithRefresh = createMapActionsWith('refreshEventFilterRulesList');
const mapActionsWithPopup = createMapActionsWithPopup({
  createEventFilterRuleWithPopupWithRefresh() {
    return this.$t('modals.eventFilterRule.create.success');
  },
  duplicateEventFilterRuleWithPopupWithRefresh() {
    return this.$t('modals.eventFilterRule.duplicate.success');
  },
  updateEventFilterRuleWithPopupWithRefresh() {
    return this.$t('modals.eventFilterRule.edit.success');
  },
  removeEventFilterRuleWithPopupWithRefresh() {
    return this.$t('modals.eventFilterRule.remove.success');
  },
});

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
    }),

    ...mapActionsWithRefresh(mapActionsWithPopup(mapActions({
      createEventFilterRuleWithPopupWithRefresh: 'create',
      duplicateEventFilterRuleWithPopupWithRefresh: 'create',
      updateEventFilterRuleWithPopupWithRefresh: 'update',
      removeEventFilterRuleWithPopupWithRefresh: 'remove',
    }))),
  },
};
