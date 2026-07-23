import { storeToRefs } from 'pinia';

import { useLlmStore } from '@/stores/modules/llm';

/**
 * LLM store facade.
 *
 * @returns {Object}
 */
export const useLlm = () => {
  const llmStore = useLlmStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(llmStore);

  return {
    llms: items,
    llmsMeta: meta,
    llmsPending: pending,
    llmsFetchingParams: fetchingParams,

    fetchLlmsList: llmStore.fetchList,
    fetchLlmsListWithPreviousParams: llmStore.fetchListWithPreviousParams,
    fetchLlmsListWithoutStore: llmStore.fetchListWithoutStore,
    createLlm: llmStore.create,
    updateLlm: llmStore.update,
    removeLlm: llmStore.remove,
    bulkRemoveLlms: llmStore.bulkRemove,
    bulkEnableLlms: llmStore.bulkEnable,
    bulkDisableLlms: llmStore.bulkDisable,
    fetchModelsListWithoutStore: llmStore.fetchModelsListWithoutStore,
    fetchLlmHistoryWithoutStore: llmStore.fetchLlmHistoryWithoutStore,
    fetchLlmChatsWithoutStore: llmStore.fetchLlmChatsWithoutStore,
    fetchLlmUsersWithoutStore: llmStore.fetchLlmUsersWithoutStore,
    fetchLlmMessagesWithoutStore: llmStore.fetchLlmMessagesWithoutStore,
    bulkLinkLlmHistory: llmStore.bulkLinkLlmHistory,
  };
};
