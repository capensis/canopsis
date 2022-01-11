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
      updateViewTabPositions: 'updatePositions',
      removeViewTab: 'remove',
    }),
  },
};
