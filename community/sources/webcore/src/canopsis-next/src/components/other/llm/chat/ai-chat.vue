<template>
  <v-layout class="ai-chat" column>
    <ai-chat-info-alert v-if="hasPatternItem && !emptyChat" :pattern-item="patternItem" />
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
      <v-spacer />
      <ai-chat-pattern-item-field
        v-if="hasPatternItem && emptyChat"
        v-model="patternItem"
        :items="patternsItems"
        :label="patternsItemsLabel"
      />
      <ai-chat-textarea
        v-model="prompt"
        ref="textareaElement"
        :llm.sync="llm"
        :llms="llms"
        :llms-pending="llmsPending"
        :error-message="errorMessage"
        :empty-chat="emptyChat"
        :thinking="!!thinkingMessage"
        :disabled="needRestart"
        @ask="ask"
        @stop="stop"
      />

      <ai-chat-suggestions v-if="emptyChat" @select="applySuggestion" />
    </v-layout>
  </v-layout>
</template>

<script>
import { toRefs } from 'vue';

import { useAiChat } from '@/hooks/ai/ai-chat';

import AiChatMessagesList from '@/components/other/llm/chat/ai-chat-messages-list.vue';
import AiChatGreeting from '@/components/other/llm/chat/ai-chat-greeting.vue';
import AiChatSuggestions from '@/components/other/llm/chat/ai-chat-suggestions.vue';
import AiChatTextarea from '@/components/other/llm/chat/ai-chat-textarea.vue';
import AiChatInfoAlert from '@/components/other/llm/chat/ai-chat-info-alert.vue';
import AiChatPatternItemField from '@/components/other/llm/chat/ai-chat-pattern-item-field.vue';
import AiChatWarningAlert from '@/components/other/llm/chat/ai-chat-warning-alert.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AiChatMessagesList,
    AiChatGreeting,
    AiChatSuggestions,
    AiChatTextarea,
    AiChatInfoAlert,
    AiChatPatternItemField,
    AiChatWarningAlert,
  },
  props: {
    patterns: {
      type: [Object, Array],
      default: () => ({}),
    },
    patternsItems: {
      type: Array,
      default: () => [],
    },
    jsonString: {
      type: String,
      default: '',
    },
    llms: {
      type: Array,
      default: () => [],
    },
    context: {
      type: String,
      default: '',
    },
    ruleId: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    return useAiChat(toRefs(props), emit);
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
