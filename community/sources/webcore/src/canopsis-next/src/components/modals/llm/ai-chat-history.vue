<template>
  <modal-wrapper class="ai-chat-history-modal" close>
    <template #title="">
      {{ $t('modals.aiChatHistory.title') }}
    </template>
    <template #text="">
      <div class="position-relative">
        <v-layout class="ai-chat-history-modal__info" column>
          <c-alert type="info" class="ma-0">
            <v-layout column>
              <span class="font-weight-regular">{{ chatUserDisplayName }}</span>
              <span>{{ chatContext }}</span>
            </v-layout>
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
        </v-layout>
        <v-layout class="ai-chat-history-modal__messages_wrapper" column>
          <div
            v-if="messages.length"
            ref="messagesElement"
            class="ai-chat-history-modal__messages pa-4"
            @scroll.passive="messagesScroll"
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
        </v-layout>
      </div>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="close"
      >
        {{ $t('common.close') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import { llmChatHistoryServerMessagesToMessages } from '@/helpers/entities/llm/chat/message/entity';
import { convertQueryToRequest } from '@/helpers/query';

import { useInnerModal } from '@/hooks/modals';
import { useLlm } from '@/hooks/store/modules/llm';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';

import AiChatMessagesList from '@/components/other/llm/chat/ai-chat-messages-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

const SCROLL_BOTTOM_THRESHOLD_PX = 30;

export default {
  name: MODALS.aiChatHistory,
  components: {
    AiChatMessagesList,
    ModalWrapper,
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

    const messagesElement = ref(null);

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
    const messagesScroll = ({ target }) => {
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
      messagesElement,

      messages,
      pending,
      loadMore,
      chatContext,
      chatUserDisplayName,

      close,
      messagesScroll,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat-history-modal {
  ::v-deep .v-card__text {
    padding: 0 !important;
  }

  &__info {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 2;
  }

  &__messages {
    max-height: 60vh;
    overflow-y: auto;

    &_wrapper {
      padding-top: 76px;
    }
  }
}
</style>
