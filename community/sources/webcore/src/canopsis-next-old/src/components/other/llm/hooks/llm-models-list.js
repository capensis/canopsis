import { keyBy } from 'lodash';
import { computed, ref, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { usePendingHandler } from '@/hooks/query/pending';
import { useLlm } from '@/hooks/store/modules/llm';

/**
 * Fetches LLM model options for a select field and loads them on mount.
 *
 * @returns {Object}
 */
export const useLlmModelsListForSelect = () => {
  const items = ref([]);

  const { fetchModelsListWithoutStore } = useLlm();

  const { pending, handler: fetchList } = usePendingHandler(async () => {
    items.value = await fetchModelsListWithoutStore({ params: { limit: MAX_LIMIT } });
  });

  const itemsByName = computed(() => keyBy(items.value, 'name'));

  onMounted(fetchList);

  return {
    items,
    itemsByName,
    pending,
    fetchList,
  };
};
