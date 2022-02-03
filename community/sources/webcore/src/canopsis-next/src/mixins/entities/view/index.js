import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

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
      copyView: 'copy',
      bulkCreateViewsWithoutStore: 'bulkCreateWithoutStore',
      exportViewsWithoutStore: 'exportWithoutStore',
      importViewsWithoutStore: 'importWithoutStore',
    }),

    async createViewWithPopup({ data }) {
      try {
        await this.createView({ data });

        this.$popups.success({ text: this.$t('modals.view.success.create') });
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.fail.create') });
      }
    },

    async updateViewWithPopup({ id, data }) {
      try {
        await this.updateView({ id, data });

        this.$popups.success({ text: this.$t('modals.view.success.edit') });
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.fail.edit') });
      }
    },

    async copyViewWithPopup({ id, data }) {
      try {
        await this.createView({ id, data });

        this.$popups.success({ text: this.$t('modals.view.success.duplicate') });
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.fail.duplicate') });
      }
    },

    async removeViewWithPopup({ id }) {
      try {
        await this.removeView({ id });

        this.$popups.success({ text: this.$t('modals.view.success.delete') });
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.fail.delete') });
      }
    },
  },
};
