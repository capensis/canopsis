<template>
  <v-layout column>
    <v-tabs
      v-model="subTab"
      slider-color="primary"
      centered
    >
      <v-tab>{{ $t('llm.promptsHistory.tabs.allPrompts') }}</v-tab>
      <v-tab-item>
        <llm-prompts-history-all-prompts-table
          :query="query"
          :local-search="localSearch"
          :data="data"
          :pending="pending"
          :total-items="totalItems"
          :options="options"
          :headers="headers"
          @update:localSearch="onUpdateLocalSearch"
          @submit-search="submitSearch"
          @clear-search="clearSearch"
          @not-related-change="onNotRelatedChange"
          @group-by-chat-change="onGroupByChatChange"
          @update:options="updateOptions"
          @see-chat="onSeeChat"
        />
      </v-tab-item>

      <v-tab>{{ $t('llm.promptsHistory.tabs.byUser') }}</v-tab>
      <v-tab-item>
        Nothing
      </v-tab-item>
    </v-tabs>
  </v-layout>
</template>

<script>
import { ref, computed, toRef, watch } from 'vue';

import { LLM_PROMPTS_HISTORY_VIEWS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import { useLlmPromptsHistory } from '../hooks/llm-prompts-history';

import LlmPromptsHistoryAllPromptsTable from './llm-prompts-history-all-prompts-table.vue';

export default {
  components: {
    LlmPromptsHistoryAllPromptsTable,
  },
  props: {
    llmId: {
      type: String,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const subTab = ref(0);

    const llmIdRef = toRef(props, 'llmId');
    const viewRef = computed(() => (subTab.value === 0
      ? LLM_PROMPTS_HISTORY_VIEWS.all
      : LLM_PROMPTS_HISTORY_VIEWS.byUser));

    const {
      data,
      meta,
      pending,
      query,
      updateQuery,
      options,
      updateOptions,
    } = useLlmPromptsHistory(llmIdRef, viewRef);

    const localSearch = ref(query.value.search ?? '');

    watch(() => query.value.search, (search) => {
      localSearch.value = search ?? '';
    });

    const totalItems = computed(() => meta.value?.total_count ?? 0);

    const headers = computed(() => [
      { text: t('llm.promptsHistory.columns.userName'), value: 'user_name' },
      { text: t('llm.promptsHistory.columns.datetime'), value: 'datetime' },
      { text: t('llm.promptsHistory.columns.tokensUsed'), value: 'tokens_used' },
      { text: t('llm.promptsHistory.columns.modal'), value: 'modal', sortable: false },
      { text: t('llm.promptsHistory.columns.name'), value: 'prompt_name', sortable: false },
      { text: t('llm.promptsHistory.columns.usage'), value: 'usage', sortable: false },
      { text: t('llm.promptsHistory.columns.canopsisRelated'), value: 'canopsis_related', sortable: false },
      { text: t('llm.promptsHistory.columns.prompt'), value: 'prompt', sortable: false },
      { text: t('llm.promptsHistory.columns.seeChat'), value: 'see_chat', sortable: false, width: '88px' },
    ]);

    const onUpdateLocalSearch = (value) => {
      localSearch.value = value;
    };

    const submitSearch = () => {
      updateQuery({ ...query.value, search: localSearch.value ?? '', page: 1 });
    };

    const clearSearch = () => {
      localSearch.value = '';
      updateQuery({ ...query.value, search: '', page: 1 });
    };

    const onNotRelatedChange = (value) => {
      updateQuery({ ...query.value, not_related_to_canopsis: value, page: 1 });
    };

    const onGroupByChatChange = (value) => {
      updateQuery({ ...query.value, group_by_chat: value, page: 1 });
    };

    const onSeeChat = (item) => {
      emit('see-chat', item);
    };

    return {
      subTab,
      data,
      pending,
      query,
      options,
      updateOptions,
      totalItems,
      headers,
      localSearch,

      onUpdateLocalSearch,
      submitSearch,
      clearSearch,
      onNotRelatedChange,
      onGroupByChatChange,
      onSeeChat,
    };
  },
};
</script>
