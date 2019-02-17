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

    async createWebhookWithPopup(data) {
      await this.createWebhook({ data });
      this.refreshWebhooksList();
      this.addSuccessPopup({ text: this.$t('modals.webhook.create.success') });
    },

    async duplicateWebhookWithPopup(data) {
      await this.createWebhook({ data });
      this.refreshWebhooksList();
      this.addSuccessPopup({ text: this.$t('modals.webhook.duplicate.success') });
    },

    async editEventFilterRuleWithPopup(id, data) {
      await this.editWebhook({ id, data });
      this.refreshWebhooksList();
      this.addSuccessPopup({ text: this.$t('modals.webhook.edit.success') });
    },
  },
};
