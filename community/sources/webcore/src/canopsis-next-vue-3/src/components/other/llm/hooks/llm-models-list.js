import { keyBy } from 'lodash';
import { computed, ref, unref, watch } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { usePendingHandler } from '@/hooks/query/pending';
import { useLlm } from '@/hooks/store/modules/llm';

/**
 * Fetches LLM model options for a select field by LLM type.
 *
 * @param {import('vue').Ref<string>|string} type - LLM model type
 * @returns {Object}
 */
export const useLlmModelsListForSelect = (type) => {
  const items = ref([]);

  const { fetchModelsListWithoutStore } = useLlm();

  const { pending, handler: fetchList } = usePendingHandler(async () => {
    const typeValue = unref(type);

    if (!typeValue) {
      items.value = [];

      return;
    }

    items.value = await fetchModelsListWithoutStore({
      params: { limit: MAX_LIMIT, type: typeValue },
    });
  });

  const itemsByName = computed(() => keyBy(items.value, 'name'));

  watch(() => unref(type), fetchList, { immediate: true });

  return {
    items,
    itemsByName,
    pending,
    fetchList,
  };
};
