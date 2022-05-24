import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view/widget');

export const entitiesWidgetMixin = {
  computed: {
    ...mapGetters({
      getWidgetById: 'getItemById',
    }),
  },
  methods: {
    ...mapActions({
      fetchWidgetWithoutStore: 'fetchItemWithoutStore',
      createWidget: 'create',
      updateWidget: 'update',
      copyWidget: 'copy',
      removeWidget: 'remove',
      updateWidgetGridPositions: 'updateGridPositions',
      createWidgetFilter: 'createWidgetFilter',
      updateWidgetFilter: 'updateWidgetFilter',
      removeWidgetFilter: 'removeWidgetFilter',
      fetchWidgetFilter: 'fetchWidgetFilter',
    }),
  },
};
