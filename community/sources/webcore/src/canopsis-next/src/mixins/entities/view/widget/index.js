import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('view/widget');

/**
 * @mixin Helpers for the widget entity
 */
export const entitiesWidgetMixin = {
  methods: {
    ...mapActions({
      createWidget: 'create',
      updateWidget: 'update',
      updateWidgetPositions: 'updatePositions',
      removeWidget: 'remove',
    }),
  },
};
