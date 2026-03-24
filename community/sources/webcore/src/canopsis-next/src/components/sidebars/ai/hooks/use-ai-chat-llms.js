import { computed, onMounted, ref } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { usePendingHandler } from '@/hooks/query/pending';
import { useLlm } from '@/hooks/store/modules/llm';

/**
 * Fetches LLM configurations for the AI chat UI and exposes enabled entries plus loading state.
 *
 * @returns {Object} `llms` (computed ref), `llmsPending` (ref), `fetchLlms` (function).
 */
export const useAiChatLlms = () => {
  const llmsRaw = ref([]);

  const { fetchLlmsListWithoutStore } = useLlm();

  const { pending: llmsPending, handler: fetchLlms } = usePendingHandler(async () => {
    const { data } = await fetchLlmsListWithoutStore({ params: { limit: MAX_LIMIT } });

    llmsRaw.value = data;
  });

  const llms = computed(() => llmsRaw.value.filter(llm => llm.enabled));

  onMounted(fetchLlms);

  return {
    llms,
    llmsPending,
    fetchLlms,
  };
};
