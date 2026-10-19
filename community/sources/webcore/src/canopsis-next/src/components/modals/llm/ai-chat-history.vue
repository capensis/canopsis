<template>
  <c-modal-wrapper
    class="ai-chat-history-modal"
    text-class="ai-chat-history-modal__text"
    close
  >
    <template #title>
      {{ $t('modals.aiChatHistory.title') }}
    </template>
    <template #text>
      <div class="ai-chat-history-modal__content position-relative">
        <c-layout class="ai-chat-history-modal__info" column>
          <c-alert type="info" class="ma-0">
            <c-layout column>
              <span class="font-weight-regular">
                {{ chatUserDisplayName }}
              </span>
              <span>{{ chatContext }}</span>
            </c-layout>
          </c-alert>
          <v-fade-transition>
            <v-progress-linear
              v-if="pending"
              class="ma-0"
              color="primary"
              height="3"
              indeterminate
            />
          </v-fade-transition>
        </c-layout>
        <c-layout class="ai-chat-history-modal__messages_wrapper" column>
          <div
            v-if="messages.length"
            class="ai-chat-history-modal__messages pa-4"
            @scroll.passive="messagesInfinityScroll"
          >
            <ai-chat-messages-list
              :messages="messages"
              :active-version="null"
              history
            />
          </div>
          <span
            v-else-if="!pending"
            class="grey--text"
          >
            {{ $t('common.noData') }}
          </span>
        </c-layout>
      </div>
    </template>
    <template #actions>
      <v-btn
        variant="text"
        @click="close"
      >
        {{ $t('common.close') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import { llmChatHistoryServerMessagesToMessages } from '@/helpers/entities/llm/chat/message/entity';
import { convertQueryToRequest } from '@/helpers/query';

import { useInnerModal } from '@/hooks/modals';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useLlm } from '@/hooks/store/modules/llm';

import AiChatMessagesList from '@/components/other/llm/chat/ai-chat-messages-list.vue';

const SCROLL_BOTTOM_THRESHOLD_PX = 30;

export default {
  name: MODALS.aiChatHistory,
  components: {
    AiChatMessagesList,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);
    const { fetchLlmMessagesWithoutStore } = useLlm();

    const messages = ref([]);
    const meta = ref({});

    const chat = computed(() => config.value.chat);
    const chatId = computed(() => chat.value?.chat ?? chat.value?._id);
    const chatContext = computed(() => chat.value?.context);
    const chatUserDisplayName = computed(() => chat.value?.user?.display_name ?? '');
    const hasMore = computed(() => meta.value?.total_count > messages.value.length);

    const {
      pending,
      query,
      updateQueryPage,
      fetchHandlerWithQuery: fetchList,
    } = usePendingWithLocalQuery({
      initialQuery: {
        page: 1,
        limit: PAGINATION_LIMIT,
      },
      fetchHandler: async (fetchQuery) => {
        const params = convertQueryToRequest(fetchQuery);

        params.chat = chatId.value;

        const response = await fetchLlmMessagesWithoutStore({ id: config.value.llmId, params });

        messages.value.push(...llmChatHistoryServerMessagesToMessages(response.data));
        meta.value = response.meta;
      },
    });

    /**
     * Requests the next page of chat history for infinite scroll.
     */
    const loadMore = () => updateQueryPage(query.value.page + 1);

    /**
     * Loads more history when the messages container is scrolled within the bottom threshold.
     *
     * @param {{ target: HTMLElement }} event
     */
    const messagesInfinityScroll = ({ target }) => {
      if (pending.value || !hasMore.value) {
        return;
      }

      const { scrollTop, clientHeight, scrollHeight } = target;
      const distanceToBottom = scrollHeight - (scrollTop + clientHeight);

      if (distanceToBottom <= SCROLL_BOTTOM_THRESHOLD_PX) {
        loadMore();
      }
    };

    onMounted(fetchList);

    return {
      messages,
      pending,
      loadMore,
      chatContext,
      chatUserDisplayName,

      close,
      messagesInfinityScroll,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat-history-modal {
  :deep(.v-card__text) {
    padding: 0 !important;
  }

  :deep(.ai-chat-history-modal__text) {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
    min-height: 0;
    padding: 0 !important;
    overflow-x: hidden;
  }

  :deep(.ai-chat-history-modal__text > *) {
    flex: 1 1 auto;
    width: 100%;
    min-width: 0;
    min-height: 0;
  }

  &__content,
  &__messages_wrapper,
  &__messages {
    width: 100%;
    min-width: 0;
    min-height: 0;
  }

  &__info {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 2;
  }

  &__messages {
    flex: 1 1 auto;
    height: 100%;
    overflow-y: auto;
    overflow-x: hidden;

    &_wrapper {
      display: flex;
      flex: 1 1 auto;
      padding-top: 76px;
    }
  }
}
</style>
