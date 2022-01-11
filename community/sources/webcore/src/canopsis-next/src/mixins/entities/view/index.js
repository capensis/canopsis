import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the view entity
 */
export const entitiesViewMixin = {
  computed: {
    ...mapGetters({
      viewId: 'itemId',
      viewPending: 'pending',
      view: 'item',
      getViewById: 'getItemById',
    }),
  },
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
      createView: 'create',
      updateView: 'update',
      updateViewsPositions: 'updatePositions',
      updateViewWithoutStore: 'updateWithoutStore',
      removeView: 'remove',
      bulkCreateViewsWithoutStore: 'bulkCreateWithoutStore',
      exportViewsWithoutStore: 'exportWithoutStore',
      importViewsWithoutStore: 'importWithoutStore',
    }),

    async createViewWithPopup({ data }) {
      await this.createView({ data });

      this.$popups.success({ text: this.$t('modals.view.success.create') });
    },

    async updateViewWithPopup({ id, data }) {
      await this.updateView({ id, data });

      this.$popups.success({ text: this.$t('modals.view.success.edit') });
    },

    async removeViewWithPopup({ id }) {
      await this.removeView({ id });

      this.$popups.success({ text: this.$t('modals.view.success.delete') });
    },
  },
};
