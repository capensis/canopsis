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

    async editWebhookWithPopup(id, data) {
      await this.editWebhook({ id, data });
      this.refreshWebhooksList();
      this.addSuccessPopup({ text: this.$t('modals.webhook.edit.success') });
    },

    async removeWebhookWithPopup(id) {
      await this.removeWebhook({ id });
      this.refreshWebhooksList();
      this.addSuccessPopup({ text: this.$t('modals.webhook.remove.success') });
    },
  },
};
