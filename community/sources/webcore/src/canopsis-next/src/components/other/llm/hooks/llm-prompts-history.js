import { ref, unref, watch, onMounted } from 'vue';

import { LLM_PROMPTS_HISTORY_VIEWS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useQueryOptions } from '@/hooks/query/options';
import { useLlm } from '@/hooks/store/modules/llm';

/**
 * Fetches paginated prompts history for an LLM configuration.
 *
 * @param {import('vue').Ref<string>} llmIdRef
 * @param {import('vue').Ref<string>} viewRef - LLM_PROMPTS_HISTORY_VIEWS
 * @returns {Object}
 */
export const useLlmPromptsHistory = (llmIdRef, viewRef) => {
  const { fetchPromptsHistoryWithoutStore } = useLlm();
  const data = ref([]);
  const meta = ref({});

  const {
    pending,
    query,
    updateQuery,
    resetQuery,
    fetchHandlerWithQuery: fetchList,
  } = usePendingWithLocalQuery({
    initialQuery: {
      page: 1,
      itemsPerPage: 10,
      search: '',
      sortBy: [],
      sortDesc: [],
      not_related_to_canopsis: false,
      group_by_chat: false,
      view: unref(viewRef) ?? LLM_PROMPTS_HISTORY_VIEWS.all,
    },
    fetchHandler: async (fetchQuery) => {
      const id = unref(llmIdRef);
      if (!id) {
        data.value = [];
        meta.value = {};

        return;
      }

      const params = {
        ...convertQueryToRequest(fetchQuery),
        not_related_to_canopsis: fetchQuery.not_related_to_canopsis,
        group_by_chat: fetchQuery.group_by_chat,
        view: fetchQuery.view,
      };

      const response = await fetchPromptsHistoryWithoutStore({ id, params });
      const rows = response?.data ?? [];

      data.value = rows.map((row, index) => ({
        ...row,
        _id: row._id ?? row.id ?? `llm-prompt-history-${index}`,
      }));
      meta.value = response?.meta ?? {};
    },
  });

  const { options, updateOptions } = useQueryOptions(query, updateQuery);

  watch(llmIdRef, () => {
    updateQuery({ ...query.value, page: 1 });
  });

  watch(viewRef, (newView) => {
    updateQuery({ ...query.value, view: newView, page: 1 });
  });

  onMounted(() => fetchList());

  return {
    data,
    meta,
    pending,
    query,
    updateQuery,
    resetQuery,
    options,
    updateOptions,
    fetchList,
  };
};
