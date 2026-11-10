<template>
  <c-advanced-data-table
    :headers="headers"
    :items="preparedUsers"
    :loading="usersPending"
    :total-items="totalItems"
    :options="options"
    :no-data-text="t('llm.promptsHistory.userHistoryEmpty')"
    :search-label="t('llm.promptsHistory.searchByUserName')"
    item-key="_id"
    hide-actions
    search
    expand
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #input_tokens="{ item }">
      {{ n(item.inputTokensUsed) }}
    </template>

    <template #output_tokens="{ item }">
      {{ n(item.output_tokens ?? 0) }}
    </template>

    <template #last_used="{ item }">
      {{ formatDate(item.last_used) }}
    </template>

    <template #expand="{ item }">
      <div class="pa-4 secondary lighten-2">
        <v-card>
          <v-card-text>
            <llm-prompts-history-table
              :llm-id="llmId"
              :user-id="item._id"
              :search-label="t('llm.promptsHistory.searchByModalOrName')"
            />
          </v-card-text>
        </v-card>
      </div>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { useChildObserver } from '@/hooks/observer';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useLlm } from '@/hooks/store/modules/llm';

import LlmPromptsHistoryTable from './llm-prompts-history-table.vue';

export default {
  components: {
    LlmPromptsHistoryTable,
  },

  props: {
    llmId: {
      type: String,
      required: true,
    },
  },

  setup(props) {
    const { t, n } = useI18n();
    const { fetchLlmUsersWithoutStore } = useLlm();

    const {
      data: users,
      meta,
      pending: usersPending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        search: '',
        sortBy: [],
        sortDesc: [],
      },
      fetchListHandler: async ({ params }) => fetchLlmUsersWithoutStore({
          id: props.llmId,
          params,
        }),
    });

    const totalItems = computed(() => meta.value?.total_count ?? 0);
    const headers = computed(() => [
      { text: t('llm.promptsHistory.columns.userName'), value: 'display_name' },
      { text: t('llm.promptsHistory.columns.inputTokensUsed'), value: 'input_tokens' },
      { text: t('llm.promptsHistory.columns.outputTokensUsed'), value: 'output_tokens' },
      { text: t('llm.promptsHistory.columns.lastUsed'), value: 'last_used' },
    ]);
    const preparedUsers = computed(() => users.value.map(user => ({
      ...user,
      inputTokensUsed: (user.input_tokens ?? 0) + (user.cache_write_tokens ?? 0) + (user.cache_read_tokens ?? 0),
    })));
    const formatDate = value => convertDateToString(value, 'long', '-');

    useChildObserver({ handler: fetchList, key: '$refresh' });

    onMounted(fetchList);

    return {
      t,
      n,
      preparedUsers,
      usersPending,
      totalItems,
      headers,
      options,

      updateOptions,
      formatDate,
    };
  },
};
</script>
