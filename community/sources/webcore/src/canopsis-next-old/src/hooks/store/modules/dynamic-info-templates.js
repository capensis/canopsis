import { get } from 'lodash';

import { ASSOCIATIVE_TABLES_NAMES } from '@/constants';

import { useAssociativeTableStoreModule } from '@/hooks/store/modules/associative-table';

/**
 * Provides CRUD helpers for dynamic info templates stored via the associative table API.
 *
 * @returns {Object} Bound async functions for fetching and mutating templates
 */
export const useDynamicInfoTemplates = () => {
  const { useActions } = useAssociativeTableStoreModule();

  const { fetchAssociativeTable, updateAssociativeTable } = useActions({
    fetchAssociativeTable: 'fetch',
    updateAssociativeTable: 'update',
  });

  /**
   * Fetches dynamic info templates from the associative table.
   *
   * @returns {Promise<Array>} Templates taken from associative table content (`templates` key).
   */
  const fetchDynamicInfoTemplatesList = async () => {
    const content = await fetchAssociativeTable({
      name: ASSOCIATIVE_TABLES_NAMES.dynamicInfoTemplates,
    });

    return get(content, 'templates', []);
  };

  /**
   * Persists the given payloads for the dynamic info templates associative table entry.
   *
   * @param {Object} data - Payload sent to associative table update (e.g. `{ templates: [...] }`).
   *
   * @returns {Promise<Array>} Templates array returned in the associative table API response content.
   */
  const updateDynamicInfoTemplateMethod = async (data) => {
    const content = await updateAssociativeTable({
      name: ASSOCIATIVE_TABLES_NAMES.dynamicInfoTemplates,
      data,
    });

    return get(content, 'templates', []);
  };

  /**
   * Appends a new template record and persists the full list server-side.
   *
   * @param {Object} params - Creation parameters.
   * @param {Object} params.data - New template payload to append.
   *
   * @returns {Promise<Array>} Full templates array after persistence.
   */
  const createDynamicInfoTemplate = async ({ data }) => {
    const templates = await fetchDynamicInfoTemplatesList();

    return updateDynamicInfoTemplateMethod({
      templates: [...templates, data],
    });
  };

  /**
   * Replaces one template entry by `_id` and persists the full list server-side.
   *
   * @param {Object} params - Update parameters.
   * @param {string} params.id - Template `_id` to replace.
   * @param {Object} params.data - Replacement template payload.
   *
   * @returns {Promise<Array>} Full templates array after persistence.
   */
  const updateDynamicInfoTemplate = async ({ id, data }) => {
    const templates = await fetchDynamicInfoTemplatesList();

    return updateDynamicInfoTemplateMethod({
      templates: templates.map(item => (item._id === id ? data : item)),
    });
  };

  /**
   * Removes a template entry by `_id` and persists the remaining list server-side.
   *
   * @param {Object} params - Removal parameters.
   * @param {string} params.id - Template `_id` to delete.
   *
   * @returns {Promise<Array>} Full templates array after persistence.
   */
  const removeDynamicInfoTemplate = async ({ id }) => {
    const templates = await fetchDynamicInfoTemplatesList();

    return updateDynamicInfoTemplateMethod({
      templates: templates.filter(item => item._id !== id),
    });
  };

  /**
   * Removes multiple templates by `_id` and persists the remaining list server-side.
   *
   * @param {Object} params - Removal parameters.
   * @param {Array<{ _id: string }>} [params.data] - Id objects (e.g. from `pickIds`) for templates to delete.
   *
   * @returns {Promise<Array>} Empty array; mass-actions `showErrorPopups` expects API status rows, not templates.
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
    createDynamicInfoTemplate,
    updateDynamicInfoTemplate,
    removeDynamicInfoTemplate,
    bulkRemoveDynamicInfoTemplates,
  };
};
