import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the view tab entity
 */
export const entitiesViewTabMixin = {
  methods: {
    ...mapActions({
      fetchViewTab: 'fetchViewTab',
      createViewTab: 'createViewTab',
      updateViewTab: 'updateViewTab',
      copyViewTab: 'copyViewTab',
      removeViewTab: 'removeViewTab',
      updateViewTabPositions: 'updateViewTabPositions',
    }),

    async updateViewTabAndFetch({ id, data }) {
      await this.updateViewTab({ id, data });

      return this.fetchViewTab({ id });
    },
  },
};
