<template>
  <v-layout
    class="pa-2 secondary lighten-2"
    column
  >
    <c-advanced-data-table
      :headers="headers"
      :items="data"
      :loading="promptsPending"
      :total-items="meta.total_count ?? 0"
      :options="options"
      :no-data-text="$t('llm.promptsHistory.expandUserHistory.empty')"
      :search-label="$t('llm.promptsHistory.searchByModalOrName')"
      item-key="_id"
      hide-actions
      search
      advanced-pagination
      @update:options="updateOptions"
    >
      <template #toolbar>
        <v-layout
          class="gap-3 px-4"
          justify-end
          align-center
        >
          <v-flex shrink>
            <c-enabled-field
              v-model="options.only_off_topic"
              :label="$t('llm.promptsHistory.groupByChat')"
              class="mt-0"
              color="primary"
              hide-details
            />
          </v-flex>
        </v-layout>
      </template>
      <template #timestamp="{ item }">
        {{ item.timestamp | date('long', '-') }}
      </template>
      <template #actions="{ item }">
        <c-action-btn
          :tooltip="$t('common.seeChat')"
          icon="forum"
          @click="$emit('see-chat', item)"
        />
      </template>
    </c-advanced-data-table>
  </v-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { convertQueryToRequest } from '@/helpers/query';

import { useI18n } from '@/hooks/i18n';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useLlm } from '@/hooks/store/modules/llm';

export default {
  props: {
    llmId: {
      type: String,
      required: true,
    },
    userId: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { fetchLlmHistoryWithoutStore } = useLlm();

    const {
      data: items,
      meta,
      pending,
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
        only_off_topic: false,
        user: '',
      },
      fetchListHandler: async ({ params }) => fetchLlmHistoryWithoutStore({ id: props.llmId, params }),
      convertQueryToRequest: (query) => {
        const newQuery = convertQueryToRequest(query);

        newQuery.only_off_topic = options.only_off_topic;
        newQuery.user = props.userId;

        return newQuery;
      },
    });

    const headers = computed(() => [
      { text: t('llm.promptsHistory.expandUserHistory.lastPromptDate'), value: 'timestamp' },
      { text: t('llm.promptsHistory.columns.tokensUsed'), value: 'tokens' },
      { text: t('llm.promptsHistory.columns.context'), value: 'context', sortable: false },
      { text: t('llm.promptsHistory.columns.name'), value: 'rule.name', sortable: false },
      { text: t('common.seeChat'), value: 'actions', sortable: false, width: '88px' },
    ]);

    onMounted(fetchList);

    return {
      items,
      meta,
      options,
      updateOptions,
      pending,
      headers,
    };
  },
};
</script>
