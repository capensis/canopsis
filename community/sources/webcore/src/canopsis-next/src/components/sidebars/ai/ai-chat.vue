<template>
  <v-layout class="ai-chat" column>
    <ai-chat-warning-alert />
    <div
      ref="bodyElement"
      class="ai-chat__body layout column gap-6 pa-4"
    >
      <ai-chat-greeting v-if="emptyChat" />
      <v-layout
        class="ai-chat__messages gap-3"
        column
      >
        <template v-for="message in messages">
          <ai-chat-pattern
            v-if="message.patterns"
            :key="message._id"
            :version="message.version"
            :active="message.version === activeVersion"
            :patterns="message.patterns"
            :original-version="message.originalVersion"
            @restore:version="restoreVersion"
          />
          <ai-chat-message
            v-else
            :key="`${message._id}-message`"
            :message="message"
          />
        </template>
        <ai-chat-message
          v-if="thinkingMessage"
          :key="thinkingMessage._id"
          :message="thinkingMessage"
        />
      </v-layout>
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

import AiChatGreeting from './ai-chat-greeting.vue';
import AiChatMessage from './ai-chat-message.vue';
import AiChatPattern from './ai-chat-pattern.vue';
import AiChatSuggestions from './ai-chat-suggestions.vue';
import AiChatTextarea from './ai-chat-textarea.vue';
import AiChatWarningAlert from './ai-chat-warning-alert.vue';

export default {
  name: SIDE_BARS.aiChat,
  components: {
    AiChatGreeting,
    AiChatMessage,
    AiChatPattern,
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

  /* Override Vuetify `.layout { flex: 1 1 auto }` so `padding-bottom` stays under the last message. */
  &__messages {
    flex: 0 0 auto !important;
    min-height: 0;

    > * {
      flex: 0 0 auto !important;
    }
  }
}
</style>
