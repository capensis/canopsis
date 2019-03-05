import { createNamespacedHelpers } from 'vuex';

import popupMixin from '@/mixins/popup';
import { mapActionsWith, mapActionsWithPopup } from '@/helpers/store';

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
    }),

    ...mapActionsWith(mapActionsWithPopup(mapActions({
      createEventFilterRuleWithPopupWithRefresh: 'create',
      updateEventFilterRuleWithPopupWithRefresh: 'update',
      removeEventFilterRuleWithPopupWithRefresh: 'remove',
    })), 'refreshEventFilterRulesList'),

    async createEventFilterRuleWithPopup(data) {
      await this.createEventFilterRule({ data });

      this.refreshEventFilterRulesListWithSuccessPopup({
        text: this.$t('modals.eventFilterRule.create.success'),
      });
    },

    async duplicateEventFilterRuleWithPopup(data) {
      await this.createEventFilterRule({ data });

      this.refreshEventFilterRulesListWithSuccessPopup({
        text: this.$t('modals.eventFilterRule.duplicate.success'),
      });
    },

    async updateEventFilterRuleWithPopup(id, data) {
      await this.updateEventFilterRule({ id, data });

      this.refreshEventFilterRulesListWithSuccessPopup({
        text: this.$t('modals.eventFilterRule.edit.success'),
      });
    },

    async removeEventFilterRuleWithPopup(id) {
      await this.removeEventFilterRule({ id });

      this.refreshEventFilterRulesListWithSuccessPopup({
        text: this.$t('modals.eventFilterRule.remove.success'),
      });
    },

    refreshEventFilterRulesListWithSuccessPopup(popup) {
      this.refreshEventFilterRulesList();
      this.addSuccessPopup(popup);
    },
  },
};
