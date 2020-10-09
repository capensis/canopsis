import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('webhook');

export default {
  computed: {
    ...mapGetters({
      webhooksPending: 'pending',
      webhooks: 'items',
      webhooksMeta: 'meta',
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
