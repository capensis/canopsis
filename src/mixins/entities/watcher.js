import { createNamespacedHelpers } from 'vuex';

import popupMixin from '@/mixins/popup';

const { mapGetters, mapActions } = createNamespacedHelpers('watcher');

/**
 * @mixin
 */
export default {
  mixins: [popupMixin],
  computed: {
    ...mapGetters({
      getWatchersListByWidgetId: 'getListByWidgetId',
      getWatchersPendingByWidgetId: 'getPendingByWidgetId',
      getWatcher: 'getItem',
    }),
    watchers() {
      return this.getWatchersListByWidgetId(this.widget._id);
    },
    watchersPending() {
      return this.getWatchersPendingByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchWatcherItem: 'fetchItem',
      fetchWatchersList: 'fetchList',
      create: 'create',
      edit: 'edit',
    }),

    async createWatcher(data) {
      await this.create({ data });
      this.addSuccessPopup({ text: this.$t('modals.createWatcher.success.create') });
    },

    async duplicateWatcher(data) {
      await this.create({ data });
      this.addSuccessPopup({ text: this.$t('modals.createWatcher.success.duplicate') });
    },

    async editWatcher(data) {
      await this.edit({ data });
      this.addSuccessPopup({ text: this.$t('modals.createWatcher.success.edit') });
    },
  },
};
