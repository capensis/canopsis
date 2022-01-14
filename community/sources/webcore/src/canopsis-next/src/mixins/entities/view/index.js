import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

/**
 * TODO: move activeView mechanism to another mixin
 */

/**
 * @mixin Helpers for the view entity
 */
export const entitiesViewMixin = {
  computed: {
    ...mapGetters({
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
