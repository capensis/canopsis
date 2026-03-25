<template>
  <v-layout class="ai-chat" column>
    <v-layout
      class="ai-chat__body gap-6 pa-4"
      column
    >
      <h5>{{ changedPatterns.join(', ') }}</h5>
      <ai-chat-greeting v-if="emptyChat" />
      <v-layout
        class="gap-3"
        column
      >
        <component
          v-for="message in messages"
          :is="message.type"
          :key="message.id"
          v-bind="message"
        />
      </v-layout>
    </v-layout>
    <v-layout class="gap-6 pa-4 pt-0" column>
      <ai-chat-textarea
        :empty-chat="emptyChat"
        @ask="ask"
        @stop="stop"
      />

      <ai-chat-suggestions @select="applySuggestion" />
    </v-layout>
    <portal v-if="socketRoom" :to="$constants.PORTALS_NAMES.additionalSidebarTitle">
      <v-btn
        class="white--text"
        color="white"
        outlined
        @click="restart"
      >
        <v-icon class="mr-2" color="white" small>
          refresh
        </v-icon>
        {{ $t('common.restart') }}
      </v-btn>
    </portal>
  </v-layout>
</template>

<script>
import { computed, ref, onBeforeUnmount } from 'vue';

import { SOCKET_ROOMS } from '@/config';
import { SIDE_BARS } from '@/constants';

import { useSocket } from '@/hooks/socket';

import AiChatGreeting from './ai-chat-greeting.vue';
import AiChatMessage from './ai-chat-message.vue';
import AiChatPattern from './ai-chat-pattern.vue';
import AiChatSuggestions from './ai-chat-suggestions.vue';
import AiChatTextarea from './ai-chat-textarea.vue';

export default {
  name: SIDE_BARS.aiChat,
  components: {
    AiChatGreeting,
    AiChatMessage,
    AiChatPattern,
    AiChatSuggestions,
    AiChatTextarea,
  },
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const messages = ref([]);
    const socketRoom = ref(null);

    const changedPatterns = computed(() => Object.keys(props.sidebar.config?.patterns ?? {}));
    const emptyChat = computed(() => messages.value.length === 0);

    const socket = useSocket();

    /**
     * Invoked when the `llmchat` room delivers a server message; extend to update `messages` or streaming state.
     * Payload shape follows `SocketRoom.call` (API-defined).
     */
    const socketListener = () => {};

    /**
     * Joins `SOCKET_ROOMS.llmChat` merging `socketRoomData` from the sidebar with `config: llm`.
     * Registers `socketListener` on the resulting `SocketRoom`.
     *
     * @param {Object} llm - LLM config merged into the join payload as `config` (e.g. selected model row).
     */
    const joinSocketRoom = (llm) => {
      const joinData = {
        ...props.sidebar.config?.socketRoomData,

        config: llm,
      };

      socketRoom.value = socket.join(SOCKET_ROOMS.llmChat, joinData, true).addListener(socketListener);
    };

    /**
     * Leaves the LLM chat socket room and clears the cached `SocketRoom` reference (runs on unmount).
     */
    const leaveSocketRoom = () => {
      socketRoom.value?.leave?.()?.removeListener?.(socketListener);
      socketRoom.value = null;
    };

    /**
     * Ensures the chat room is joined, then sends the user prompt over the socket via `SocketRoom.send`.
     *
     * @param {Object} params
     * @param {Object} params.llm - LLM config passed into `joinSocketRoom` when the room is not open yet.
     * @param {string} params.prompt - User message text sent in the outbound `{ prompt }` payload.
     */
    const ask = ({ llm, prompt }) => {
      if (!socketRoom.value) {
        joinSocketRoom(llm);
      }

      socketRoom.value.send({
        prompt,

        type: 'send',
      });
    };

    /**
     * Sends `{ type: 'cancel' }` on the room channel to stop an in-flight request when the room is joined.
     */
    const stop = () => socketRoom.value?.send({
      type: 'cancel',
    });

    /**
     * Handles suggestion chips from `ai-chat-suggestions`; connect to `prompt` / `ask` when implemented.
     */
    const applySuggestion = () => {};

    onBeforeUnmount(leaveSocketRoom);

    return {
      changedPatterns,
      messages,
      socketRoom,
      emptyChat,

      ask,
      stop,
      applySuggestion,
      restart: leaveSocketRoom,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat {
  --header-height: 64px;
  height: calc(100% - var(--header-height));
  min-height: 0;

  &__body {
    flex: 1 1 auto;
    min-height: 0;
    overflow-y: auto;
  }
}
</style>
