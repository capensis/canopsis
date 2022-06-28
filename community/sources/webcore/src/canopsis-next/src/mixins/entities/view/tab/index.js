import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('view/tab');

/**
 * @mixin Helpers for the view tab entity
 */
export const entitiesViewTabMixin = {
  methods: {
    ...mapActions({
      fetchViewTab: 'fetchItem',
      createViewTab: 'create',
      updateViewTab: 'update',
      copyViewTab: 'copy',
      removeViewTab: 'remove',
      updateViewTabPositions: 'updatePositions',
    }),

    async updateViewTabAndFetch({ id, data }) {
      await this.updateViewTab({ id, data });

      return this.fetchViewTab({ id });
    },
  },
};
