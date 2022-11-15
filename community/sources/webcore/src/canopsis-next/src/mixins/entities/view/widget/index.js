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
      fetchWidget: 'fetchItem',
      fetchWidgetWithoutStore: 'fetchItemWithoutStore',
      createWidget: 'create',
      updateWidget: 'update',
      copyWidget: 'copy',
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
