import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the view entity
 */
export const entitiesViewMixin = {
  methods: {
    ...mapActions({
      createView: 'createView',
      updateView: 'updateView',
      updateViewsPositions: 'updateViewPositions',
      updateViewWithoutStore: 'updateWithoutStoreView',
      removeView: 'removeView',
      copyView: 'copyView',
      exportViewsWithoutStore: 'exportViewWithoutStore',
      importViewsWithoutStore: 'importViewWithoutStore',
    }),

    async createViewWithPopup({ data }) {
      try {
        await this.createView({ data });

        this.$popups.success({ text: this.$t('modals.view.success.create') });
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.fail.create') });

        throw err;
      }
    },

    async updateViewWithPopup({ id, data }) {
      try {
        await this.updateView({ id, data });

        this.$popups.success({ text: this.$t('modals.view.success.edit') });
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.fail.edit') });

        throw err;
      }
    },

    async copyViewWithPopup({ id, data }) {
      try {
        await this.copyView({ id, data });

        this.$popups.success({ text: this.$t('modals.view.success.duplicate') });
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.fail.duplicate') });

        throw err;
      }
    },

    async removeViewWithPopup({ id }) {
      try {
        await this.removeView({ id });

        this.$popups.success({ text: this.$t('modals.view.success.delete') });
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.fail.delete') });

        throw err;
      }
    },
  },
};
