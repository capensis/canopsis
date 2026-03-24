<template>
  <v-layout class="ai-chat" column>
    <v-layout
      class="ai-chat__body gap-6 pa-4"
      column
    >
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
  </v-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { SIDE_BARS } from '@/constants';

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
  setup() {
    const messages = ref([]);

    const emptyChat = computed(() => messages.value.length === 0);

    const ask = () => {};

    const stop = () => {};

    const applySuggestion = () => {};

    return {
      messages,
      emptyChat,

      ask,
      stop,
      applySuggestion,
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
