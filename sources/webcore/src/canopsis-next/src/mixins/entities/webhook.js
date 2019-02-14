import { createNamespacedHelpers } from 'vuex';

import popupMixin from '@/mixins/popup';

const { mapActions, mapGetters } = createNamespacedHelpers('webhook');

export default {
  mixins: [popupMixin],
  computed: {
    ...mapGetters({
      webhooksPending: 'pending',
      webhooks: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchWebhooksList: 'fetchList',
      refreshWebhooksList: 'fetchListWithPreviousParams',
      removeWebhook: 'remove',
      createWebhook: 'create',
      editWebhook: 'edit',
    }),

    async createEventFilterRuleWithPopup(rule) {
      await this.createEventFilterRule({ data: rule });
      this.refreshEventFilterList();
      this.addSuccessPopup({ text: this.$t('modals.eventFilterRule.create.success') });
    },

    async duplicateEventFilterRuleWithPopup(rule) {
      await this.createEventFilterRule({ data: rule });
      this.refreshEventFilterList();
      this.addSuccessPopup({ text: this.$t('modals.eventFilterRule.duplicate.success') });
    },

    async editEventFilterRuleWithPopup(ruleId, editedRule) {
      await this.editEventFilterRule({ id: ruleId, data: editedRule });
      this.refreshEventFilterList();
      this.addSuccessPopup({ text: this.$t('modals.eventFilterRule.edit.success') });
    },
  },
};
