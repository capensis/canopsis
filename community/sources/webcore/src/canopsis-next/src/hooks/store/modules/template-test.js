import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Template Test Store Module.
 *
 * @returns {Object} An object containing getters and actions for the template test.
 */
const useTemplateTestStoreModule = () => useStoreModuleHooks('template/test');

/**
 * Hook to access template test store.
 *
 * @returns {Object} An object containing:
 * - Actions to fetch lists without store, create, update, and remove template tests.
 */
export const useTemplateTest = () => {
  const { useActions } = useTemplateTestStoreModule();

  const actions = useActions({
    fetchTemplateTestListWithoutStore: 'fetchListWithoutStore',
    createTemplateTest: 'create',
    updateTemplateTest: 'update',
    removeTemplateTest: 'remove',
  });

  return {
    ...actions,
  };
};
