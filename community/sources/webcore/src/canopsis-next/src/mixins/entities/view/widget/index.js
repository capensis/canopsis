import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view/widget');

/**
 * @mixin Helpers for the widget entity
 */
export const entitiesWidgetMixin = {
  computed: {
    ...mapGetters({
      getWidgetById: 'getItemById',
    }),
  },
  methods: {
    ...mapActions({
      createWidget: 'create',
      updateWidget: 'update',
      copyWidget: 'copy',
      removeWidget: 'remove',
      updateWidgetGridPositions: 'updateGridPositions',
    }),
  },
};
