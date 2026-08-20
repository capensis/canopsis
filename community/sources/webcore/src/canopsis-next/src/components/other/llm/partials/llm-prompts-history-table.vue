<template>
  <c-advanced-data-table
    :headers="headers"
    :items="preparedItems"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :search-label="searchLabel || $t('llm.promptsHistory.search')"
    class="llm-prompts-history-table"
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
        <v-flex v-if="!groupByChat" xs12 md3>
          <c-enabled-field
            :value="options.only_off_topic"
            :label="$t('llm.promptsHistory.notRelatedToCanopsis')"
            hide-details
            @input="onlyOffTopicChange"
          />
        </v-flex>
        <v-flex>
          <v-layout justify-end>
            <c-enabled-field
              v-model="groupByChat"
              :label="$t('llm.promptsHistory.groupByChat')"
              hide-details
            />
          </v-layout>
        </v-flex>
      </v-layout>
    </template>
    <template #timestamp="{ item }">
      {{ item.timestamp | date('long', '-') }}
    </template>
    <template #usage>
      {{ $t('pattern.patterns') }}
    </template>
    <template #rule="{ item }">
      <span v-if="item.rule?.name">{{ item.rule.name }}</span>
      <span v-else class="font-italic">{{ $t('llm.promptsHistory.ruleNotSaved') }}</span>
    </template>
    <template #off_topic="{ item }">
      <c-enabled :value="!item.off_topic" />
    </template>
    <template #prompt="{ item }">
      <span class="llm-prompts-history-table__prompt">{{ item.prompt }}</span>
    </template>
    <template #actions="{ item }">
      <c-action-btn
        :tooltip="$t('common.seeChat')"
        :aria-label="$t('common.seeChat')"
        icon="chat"
        @click="seeChat(item)"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed, ref, watch, onMounted } from 'vue';

import { MODALS, LLM_SOCKET_CONTEXTS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useLlm } from '@/hooks/store/modules/llm';
import { useChildObserver } from '@/hooks/observer';

export default {
  props: {
    llmId: {
      type: String,
      required: true,
    },
    userId: {
      type: String,
      default: '',
    },
    searchLabel: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const { t, te } = useI18n();
    const modals = useModals();
    const popups = usePopups();

    const { fetchLlmHistoryWithoutStore, fetchLlmChatsWithoutStore } = useLlm();

    const groupByChat = ref(false);

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

        if (props.userId) {
          params.user = props.userId;
        }

        return params;
      },
      fetchListHandler: async ({ params }) => (
        (groupByChat.value ? fetchLlmChatsWithoutStore : fetchLlmHistoryWithoutStore)({ id: props.llmId, params })
      ),
    });

    const totalItems = computed(() => meta.value?.total_count ?? 0);

    const headers = computed(() => [
      !props.userId && { text: t('llm.promptsHistory.columns.userName'), value: 'user.display_name' },
      {
        text: t(groupByChat.value ? 'llm.promptsHistory.expandUserHistory.lastPromptDate' : 'llm.promptsHistory.columns.datetime'),
        value: 'timestamp',
      },
      { text: t('llm.promptsHistory.columns.tokensUsed'), value: 'tokens' },
      { text: t('llm.promptsHistory.columns.modal'), value: 'context', sortable: false },
      { text: t('common.name'), value: 'rule', sortable: false },
      { text: t('llm.promptsHistory.columns.usage'), value: 'usage', sortable: false },
      !groupByChat.value && { text: t('llm.promptsHistory.columns.canopsisRelated'), value: 'off_topic', sortable: false },
      !groupByChat.value && { text: t('llm.promptsHistory.columns.prompt'), value: 'prompt', sortable: false },
      { text: t('common.seeChat'), value: 'actions', sortable: false, width: '88px' },
    ].filter(Boolean));

    const preparedItems = computed(() => items.value.map((item) => {
      const context = item.context.replace(`${LLM_SOCKET_CONTEXTS.widgetFilter}_`, '');
      const widgetMessageKey = `modals.createWidget.types.${context}.title`;

      return {
        ...item,
        context: te(widgetMessageKey)
          ? `${t(widgetMessageKey)} / ${t('common.filter')}`
          : t(`llm.promptsHistory.contextTitles.${context}`),
      };
    }));

    /**
     * Applies the "not related to Canopsis" filter and reloads the first page.

     * @param {boolean} value - Whether to restrict history to prompts not related to Canopsis.
     */
    const onlyOffTopicChange = value => updateOptions({ ...options.value, only_off_topic: value, page: 1 });

    /**
     * Opens the AI chat history modal for the row’s chat.

     * @param {Object} item - LLM history or chat aggregate row.
     */
    const seeChat = (item) => {
      const chatId = item.chat ?? item._id;

      if (!chatId) {
        popups.error({ text: t('errors.default') });

        return;
      }

      modals.show({
        name: MODALS.aiChatHistory,
        config: {
          llmId: props.llmId,
          chat: item,
        },
      });
    };

    watch(groupByChat, () => {
      if (options.value.page > 1) {
        updateOptions({ ...options.value, page: 1 });

        return;
      }

      fetchList();
    });

    useChildObserver({ handler: fetchList, key: '$refresh' });

    onMounted(fetchList);

    return {
      groupByChat,
      preparedItems,
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
.llm-prompts-history-table__prompt {
  display: -webkit-box;
  min-width: 0;
  overflow: hidden;
  -webkit-box-orient: vertical;
  line-clamp: 3;
  -webkit-line-clamp: 3;
  word-break: break-word;
}
</style>
