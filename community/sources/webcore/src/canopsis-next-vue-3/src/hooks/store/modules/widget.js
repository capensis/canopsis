import { useWidgetStore } from '@/stores/modules/widget'

/**
 * Widget store facade.
 *
 * @returns {Object}
 */
export const useWidget = () => {
  const widgetStore = useWidgetStore()

  return {
    fetchWidgetWithoutStore: widgetStore.fetchItemWithoutStore,
    createWidget: widgetStore.create,
    cloneWidget: widgetStore.clone,
    updateWidget: widgetStore.update,
    removeWidget: widgetStore.remove,
    updateWidgetGridPositions: widgetStore.updateGridPositions,
    fetchWidgetFilters: widgetStore.fetchWidgetFilters,
    createWidgetFilter: widgetStore.createWidgetFilter,
    updateWidgetFilter: widgetStore.updateWidgetFilter,
    removeWidgetFilter: widgetStore.removeWidgetFilter,
    fetchWidgetFilter: widgetStore.fetchWidgetFilter,
    updateWidgetFiltersPositions: widgetStore.updateWidgetFiltersPositions,
  }
}
