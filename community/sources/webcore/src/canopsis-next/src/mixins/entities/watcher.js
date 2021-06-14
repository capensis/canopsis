import { createNamespacedHelpers } from 'vuex';

import { CANOPSIS_STACK } from '@/constants';

import entitiesInfoMixin from '@/mixins/entities/info';

const { mapGetters, mapActions } = createNamespacedHelpers('watcher');

/**
 * @mixin
 */
export default {
  mixins: [entitiesInfoMixin],
  computed: {
    ...mapGetters({
      getWatchersListByWidgetId: 'getListByWidgetId',
      getWatchersPendingByWidgetId: 'getPendingByWidgetId',
      getWatchersErrorByWidgetId: 'getErrorByWidgetId',
      getWatcher: 'getItem',
    }),
    watchers() {
      return this.getWatchersListByWidgetId(this.widget._id);
    },
    watchersPending() {
      return this.getWatchersPendingByWidgetId(this.widget._id);
    },
    watchersError() {
      return this.getWatchersErrorByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchWatcherItem: 'fetchItem',
      fetchWatchersList: 'fetchList',
      createWatcher: 'createWatcher',
      createWatcherNg: 'createWatcherNg',
      editWatcher: 'editWatcher',
      editWatcherNg: 'editWatcherNg',
    }),

    async createWatcherWithPopup(data) {
      if (this.stack === CANOPSIS_STACK.go) {
        await this.createWatcherNg({ data });
      } else {
        await this.createWatcher({ data });
      }

      this.$popups.success({ text: this.$t('modals.createWatcher.success.create') });
    },

    async duplicateWatcherWithPopup(data) {
      if (this.stack === CANOPSIS_STACK.go) {
        await this.createWatcherNg({ data });
      } else {
        await this.createWatcher({ data });
      }

      this.$popups.success({ text: this.$t('modals.createWatcher.success.duplicate') });
    },

    async editWatcherWithPopup(data) {
      if (this.stack === CANOPSIS_STACK.go) {
        await this.editWatcherNg({ data });
      } else {
        await this.editWatcher({ data });
      }

      this.$popups.success({ text: this.$t('modals.createWatcher.success.edit') });
    },
  },
};
