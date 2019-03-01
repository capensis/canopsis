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
      createWebhook: 'create',
      updateWebhook: 'update',
      removeWebhook: 'remove',
    }),

    async createWebhookWithPopup(data) {
      await this.createWebhook({ data });

      this.refreshWebhooksListWithSuccessPopup({ text: this.$t('modals.webhook.create.success') });
    },

    async updateWebhookWithPopup(id, data) {
      await this.updateWebhook({ id, data });

      this.refreshWebhooksListWithSuccessPopup({ text: this.$t('modals.webhook.edit.success') });
    },

    async removeWebhookWithPopup(id) {
      await this.removeWebhook({ id });

      this.refreshWebhooksListWithSuccessPopup({ text: this.$t('modals.webhook.remove.success') });
    },

    refreshWebhooksListWithSuccessPopup(popup) {
      this.refreshWebhooksList();
      this.addSuccessPopup(popup);
    },
  },
};
