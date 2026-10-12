import { useDynamicInfoTemplatesStore } from '@/stores/modules/dynamic-info-templates';

/**
 * Provides CRUD helpers for dynamic info templates stored via the associative table API.
 *
 * @returns {Object} Bound async functions for fetching and mutating templates
 */
export const useDynamicInfoTemplates = () => {
  const dynamicInfoTemplatesStore = useDynamicInfoTemplatesStore();

  return {
    fetchDynamicInfoTemplatesList: dynamicInfoTemplatesStore.fetchDynamicInfoTemplatesList,
    updateDynamicInfoTemplateMethod: dynamicInfoTemplatesStore.updateDynamicInfoTemplateMethod,
    createDynamicInfoTemplate: dynamicInfoTemplatesStore.createDynamicInfoTemplate,
    updateDynamicInfoTemplate: dynamicInfoTemplatesStore.updateDynamicInfoTemplate,
    removeDynamicInfoTemplate: dynamicInfoTemplatesStore.removeDynamicInfoTemplate,
    bulkRemoveDynamicInfoTemplates: dynamicInfoTemplatesStore.bulkRemoveDynamicInfoTemplates,
  };
};
