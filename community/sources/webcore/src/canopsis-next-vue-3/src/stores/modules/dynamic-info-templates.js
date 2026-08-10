import { get } from 'lodash';
import { defineStore } from 'pinia';

import { ASSOCIATIVE_TABLES_NAMES } from '@/constants';
import { useAssociativeTableStore } from '@/stores/modules/associative-table';

export const useDynamicInfoTemplatesStore = defineStore('dynamicInfoTemplates', () => {
  const associativeTableStore = useAssociativeTableStore();

  /**
   * Fetch dynamic info templates.
   *
   * @returns {Promise<Array>}
   */
  const fetchDynamicInfoTemplatesList = async () => {
    const content = await associativeTableStore.fetch({
      name: ASSOCIATIVE_TABLES_NAMES.dynamicInfoTemplates,
    });

    return get(content, 'templates', []);
  };

  /**
   * Main update dynamic info templates method.
   *
   * @param {Object} data - Associative table payload.
   * @returns {Promise<Array>}
   */
  const updateDynamicInfoTemplateMethod = async (data) => {
    const content = await associativeTableStore.update({
      name: ASSOCIATIVE_TABLES_NAMES.dynamicInfoTemplates,
      data,
    });

    return get(content, 'templates', []);
  };

  /**
   * Add new dynamic info template into templates.
   *
   * @param {Object} params
   * @param {Object} params.data
   * @returns {Promise<Array>}
   */
  const createDynamicInfoTemplate = async ({ data }) => {
    const templates = await fetchDynamicInfoTemplatesList();

    return updateDynamicInfoTemplateMethod({
      templates: [...templates, data],
    });
  };

  /**
   * Update existing dynamic info template in templates.
   *
   * @param {Object} params
   * @param {string} params.id
   * @param {Object} params.data
   * @returns {Promise<Array>}
   */
  const updateDynamicInfoTemplate = async ({ id, data }) => {
    const templates = await fetchDynamicInfoTemplatesList();

    return updateDynamicInfoTemplateMethod({
      templates: templates.map(item => (item._id === id ? data : item)),
    });
  };

  /**
   * Remove dynamic info template from templates.
   *
   * @param {Object} params
   * @param {string} params.id
   * @returns {Promise<Array>}
   */
  const removeDynamicInfoTemplate = async ({ id }) => {
    const templates = await fetchDynamicInfoTemplatesList();

    return updateDynamicInfoTemplateMethod({
      templates: templates.filter(item => item._id !== id),
    });
  };

  /**
   * Remove dynamic info templates from templates.
   *
   * @param {Object} [params={}]
   * @param {Array<{ _id: string }>} [params.data=[]]
   * @returns {Promise<Array>}
   */
  const bulkRemoveDynamicInfoTemplates = async ({ data = [] } = {}) => {
    const idSet = new Set(data.map(item => item._id).filter(Boolean));

    if (!idSet.size) {
      return [];
    }

    const templates = await fetchDynamicInfoTemplatesList();

    await updateDynamicInfoTemplateMethod({
      templates: templates.filter(item => !idSet.has(item._id)),
    });

    return [];
  };

  return {
    fetchDynamicInfoTemplatesList,
    updateDynamicInfoTemplateMethod,
    createDynamicInfoTemplate,
    updateDynamicInfoTemplate,
    removeDynamicInfoTemplate,
    bulkRemoveDynamicInfoTemplates,

    fetchList: fetchDynamicInfoTemplatesList,
    create: createDynamicInfoTemplate,
    update: updateDynamicInfoTemplate,
    remove: removeDynamicInfoTemplate,
    bulkRemove: bulkRemoveDynamicInfoTemplates,
  };
});
