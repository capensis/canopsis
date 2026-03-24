import { useStoreModuleHooks } from '@/hooks/store';

/**
 * @returns {Object} Vuex module hooks for the `llm` namespace
 */
const useLlmStoreModule = () => useStoreModuleHooks('llm');

/**
 * Maps Vuex `llm` actions to app-facing names.
 *
 * @returns {Object}
 */
export const useLlm = () => {
  const { useActions } = useLlmStoreModule();

  const actions = useActions({
    createLlm: 'create',
    updateLlm: 'update',
    removeLlm: 'remove',
    bulkRemoveLlms: 'bulkRemove',
    fetchLlmsListWithoutStore: 'fetchListWithoutStore',
    fetchModelsListWithoutStore: 'fetchModelsListWithoutStore',
    fetchPromptsHistoryWithoutStore: 'fetchPromptsHistoryWithoutStore',
  });

  return {
    ...actions,
  };
};
