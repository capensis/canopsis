import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Template Data Store Module.
 *
 * @returns {Object} An object containing getters and actions for the template data.
 */
const useTemplateDataStoreModule = () => useStoreModuleHooks('template/data');

/**
 * Hook to access template data store.
 *
 * @returns {Object} An object containing:
 * - Actions to fetch lists without store, create, update, and remove template data.
 */
export const useTemplateData = () => {
  const { useActions } = useTemplateDataStoreModule();

  const actions = useActions({
    fetchTemplateDataListWithoutStore: 'fetchListWithoutStore',
    createTemplateData: 'create',
    updateTemplateData: 'update',
    removeTemplateData: 'remove',
  });

  return {
    ...actions,
  };
};
