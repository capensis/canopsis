import { storeToRefs } from 'pinia';

import { useWidgetTemplateStore } from '@/stores/modules/widget-template';

/**
 * Widget template store facade.
 *
 * @returns {Object}
 */
export const useWidgetTemplate = () => {
  const widgetTemplateStore = useWidgetTemplateStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(widgetTemplateStore);

  return {
    widgetTemplates: items,
    widgetTemplatesMeta: meta,
    widgetTemplatesPending: pending,
    widgetTemplatesFetchingParams: fetchingParams,

    fetchWidgetTemplatesList: widgetTemplateStore.fetchList,
    fetchWidgetTemplatesListWithPreviousParams: widgetTemplateStore.fetchListWithPreviousParams,
    fetchWidgetTemplatesListWithoutStore: widgetTemplateStore.fetchListWithoutStore,
    createWidgetTemplate: widgetTemplateStore.create,
    updateWidgetTemplate: widgetTemplateStore.update,
    removeWidgetTemplate: widgetTemplateStore.remove,
  };
};
