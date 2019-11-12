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
  },
};
