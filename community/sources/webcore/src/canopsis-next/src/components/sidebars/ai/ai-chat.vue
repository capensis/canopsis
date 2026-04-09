<template>
  <v-layout class="ai-chat" column>
    <ai-chat-warning-alert />
    <div
      ref="bodyElement"
      class="ai-chat__body layout column gap-6 pa-4"
    >
      <ai-chat-greeting v-if="emptyChat" />
      <ai-chat-messages-list
        :messages="messages"
        :thinking-message="thinkingMessage"
        :active-version="activeVersion"
        @restore:version="restoreVersion"
      />
    </div>
    <v-layout
      :justify-end="!emptyChat"
      :class="{ 'pa-4': emptyChat }"
      class="gap-6 pt-0"
      column
    >
      <ai-chat-textarea
        v-model="prompt"
        ref="textareaElement"
        :llm.sync="llm"
        :llms="llms"
        :llms-pending="llmsPending"
        :error-message="errorMessage"
        :empty-chat="emptyChat"
        :thinking="!!thinkingMessage"
        @ask="ask"
        @stop="stop"
      />

      <ai-chat-suggestions v-if="emptyChat" @select="applySuggestion" />
    </v-layout>
    <portal v-if="!emptyChat" :to="$constants.PORTALS_NAMES.additionalSidebarTitle">
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
import { toRef } from 'vue';

import { SIDE_BARS } from '@/constants';

import { useAiChat } from '@/hooks/ai/ai-chat';

import AiChatMessagesList from '@/components/other/llm/chat/ai-chat-messages-list.vue';
import AiChatGreeting from '@/components/other/llm/chat/ai-chat-greeting.vue';
import AiChatSuggestions from '@/components/other/llm/chat/ai-chat-suggestions.vue';
import AiChatTextarea from '@/components/other/llm/chat/ai-chat-textarea.vue';
import AiChatWarningAlert from '@/components/other/llm/chat/ai-chat-warning-alert.vue';

export default {
  name: SIDE_BARS.aiChat,
  components: {
    AiChatMessagesList,
    AiChatGreeting,
    AiChatSuggestions,
    AiChatTextarea,
    AiChatWarningAlert,
  },
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    return useAiChat({
      sidebar: toRef(props, 'sidebar'),
    });
  },
};
</script>

<style lang="scss" scoped>
.ai-chat {
  --header-height: 64px;
  height: calc(100% - var(--header-height));
  overflow-y: auto;
  min-height: 0;

  &__body {
    flex: 1 1 auto;
    min-height: 0;
    overflow-y: auto;
  }
}
</style>
