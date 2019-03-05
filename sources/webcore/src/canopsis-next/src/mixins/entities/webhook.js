import { createNamespacedHelpers } from 'vuex';

import popupMixin from '@/mixins/popup';
import { createMapActionsWithPopup, createMapActionsWith } from '@/helpers/store';

const { mapActions, mapGetters } = createNamespacedHelpers('webhook');

const mapActionsWithRefresh = createMapActionsWith('refreshWebhooksList');
const mapActionsWithPopup = createMapActionsWithPopup({
  createWebhookWithPopupWithRefresh() {
    return 'Created!';
  },
  updateWebhookWithPopupWithRefresh() {
    return this.$t('common.author');
  },
  removeWebhookWithPopupWithRefresh() {
    return 'Removed!';
  },
});

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

    ...mapActionsWithRefresh(mapActionsWithPopup({
      createWebhookWithPopupWithRefresh: 'createWebhook',
      updateWebhookWithPopupWithRefresh: 'updateWebhook',
      removeWebhookWithPopupWithRefresh: 'removeWebhook',
    })),
  },
};
