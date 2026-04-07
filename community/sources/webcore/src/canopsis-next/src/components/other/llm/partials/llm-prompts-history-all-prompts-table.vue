<template>
  <c-advanced-data-table
    :headers="headers"
    :items="data"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :search-label="$t('llm.promptsHistory.search')"
    class="llm-prompts-history-all-prompts-table"
    item-key="_id"
    hide-actions
    search
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar>
      <v-layout
        class="gap-3 px-4"
        align-center
        wrap
      >
        <v-flex xs12 md3>
          <c-enabled-field
            :value="options.only_off_topic"
            :label="$t('llm.promptsHistory.onlyOffTopic')"
            hide-details
            @change="onlyOffTopicChange"
          />
        </v-flex>
        <v-flex>
          <v-layout justify-end>
            <c-enabled-field
              :value="options.only_off_topic"
              :label="$t('llm.promptsHistory.onlyOffTopic')"
              hide-details
              @change="onlyOffTopicChange"
            />
          </v-layout>
        </v-flex>
      </v-layout>
    </template>
    <template #timestamp="{ item }">
      {{ item.timestamp | date('long', '-') }}
    </template>
    <template #prompt="{ item }">
      <span class="llm-prompts-history-all-prompts-table__prompt">{{ item.prompt }}</span>
    </template>
    <template #actions="{ item }">
      <c-action-btn
        :tooltip="$t('common.seeChat')"
        :aria-label="$t('common.seeChat')"
        icon="forum"
        @click="seeChat(item)"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed, onMounted } from 'vue';

import { convertQueryToRequest } from '@/helpers/query';

import { useI18n } from '@/hooks/i18n';
import { useLlm } from '@/hooks/store/modules/llm';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

export default {
  props: {
    llmId: {
      type: String,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { fetchLlmHistoryWithoutStore } = useLlm();

    const {
      data,
      meta,
      pending,
      query,
      updateQuery,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: 10,
        search: '',
        sortBy: [],
        sortDesc: [],
        only_off_topic: false,
      },
      convertQueryToRequest: (fetchQuery) => {
        const params = {
          ...convertQueryToRequest(fetchQuery),
        };

        if (fetchQuery.only_off_topic) {
          params.only_off_topic = true;
        }

        return params;
      },
      fetchListHandler: async ({ params }) => fetchLlmHistoryWithoutStore({ id: props.llmId, params }),
    });

    const totalItems = computed(() => meta.value?.total_count ?? 0);

    const headers = computed(() => [
      { text: t('llm.promptsHistory.columns.userName'), value: 'user.display_name' },
      { text: t('llm.promptsHistory.columns.datetime'), value: 'timestamp' },
      { text: t('llm.promptsHistory.columns.tokensUsed'), value: 'tokens' },
      { text: t('llm.promptsHistory.columns.context'), value: 'context', sortable: false },
      { text: t('llm.promptsHistory.columns.name'), value: 'rule.name', sortable: false },
      { text: t('llm.promptsHistory.columns.offTopic'), value: 'off_topic', sortable: false },
      { text: t('llm.promptsHistory.columns.prompt'), value: 'prompt', sortable: false },
      { text: t('common.seeChat'), value: 'actions', sortable: false, width: '88px' },
    ]);

    const onlyOffTopicChange = value => updateQuery({ ...query.value, only_off_topic: value, page: 1 });

    const seeChat = item => emit('see-chat', item);

    onMounted(fetchList);

    return {
      data,
      pending,
      totalItems,
      options,
      updateOptions,
      headers,
      onlyOffTopicChange,
      seeChat,
    };
  },
};
</script>

<style lang="scss">
.llm-prompts-history-all-prompts-table__prompt {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
