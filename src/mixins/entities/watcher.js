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
      createWatcher: 'create',
      editWatcher: 'edit',
    }),

    async createWatcherWithPopup(data) {
      await this.createWatcher({ data });
      this.addSuccessPopup({ text: this.$t('modals.createWatcher.success.create') });
    },

    async duplicateWatcherWithPopup(data) {
      await this.createWatcher({ data });
      this.addSuccessPopup({ text: this.$t('modals.createWatcher.success.duplicate') });
    },

    async editWatcherWithPopup(data) {
      await this.editWatcher({ data });
      this.addSuccessPopup({ text: this.$t('modals.createWatcher.success.edit') });
    },
  },
};
