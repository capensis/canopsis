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
    bulkEnableLlms: 'bulkEnable',
    bulkDisableLlms: 'bulkDisable',
    fetchLlmsListWithoutStore: 'fetchListWithoutStore',
    fetchModelsListWithoutStore: 'fetchModelsListWithoutStore',
    fetchLlmHistoryWithoutStore: 'fetchLlmHistoryWithoutStore',
    fetchLlmChatsWithoutStore: 'fetchLlmChatsWithoutStore',
    fetchLlmUsersWithoutStore: 'fetchLlmUsersWithoutStore',
    fetchLlmMessagesWithoutStore: 'fetchLlmMessagesWithoutStore',
    bulkLinkLlmHistory: 'bulkLinkLlmHistory',
  });

  const fetchDefaultLlmWithoutStore = () => actions.fetchLlmsListWithoutStore({
    params: { only_default: true, enabled: true },
  });

  return {
    ...actions,

    fetchDefaultLlmWithoutStore,
  };
};
