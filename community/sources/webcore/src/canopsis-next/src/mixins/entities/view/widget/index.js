import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('view/widget');

export const entitiesWidgetMixin = {
  methods: {
    ...mapActions({
      fetchWidgetWithoutStore: 'fetchItemWithoutStore',
      createWidget: 'create',
      updateWidget: 'update',
      removeWidget: 'remove',
      updateWidgetGridPositions: 'updateGridPositions',
      fetchWidgetFilters: 'fetchWidgetFilters',
      fetchWidgetFilter: 'fetchWidgetFilter',
      createWidgetFilter: 'createWidgetFilter',
      updateWidgetFilter: 'updateWidgetFilter',
      removeWidgetFilter: 'removeWidgetFilter',
      updateWidgetFiltersPositions: 'updateWidgetFiltersPositions',
    }),
  },
};
