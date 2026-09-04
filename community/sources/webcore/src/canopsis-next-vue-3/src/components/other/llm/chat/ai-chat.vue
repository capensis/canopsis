<template>
  <c-layout class="ai-chat" column>
    <ai-chat-info-alert
      v-if="hasPatternItem && !emptyChat"
      :pattern-item="patternItem"
      :context="context"
    />
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
    <c-layout
      :justify="!emptyChat ? 'end' : undefined"
      :class="{ 'pa-4': emptyChat }"
      class="ai-chat__footer pt-0"
      gap="6"
      column
    >
      <ai-chat-pattern-item-field
        v-if="hasPatternItem && emptyChat"
        v-model="patternItem"
        :items="patternsItems"
        :label="patternsItemsLabel"
      />
      <ai-chat-textarea
        v-model="prompt"
        v-model:llm="llm"
        ref="textareaElement"
        :llms="llms"
        :llms-pending="llmsPending"
        :error-message="errorMessage"
        :empty-chat="emptyChat"
        :thinking="!!thinkingMessage"
        :disabled="needRestart || patternItemMissing"
        @ask="ask"
        @stop="stop"
      />

      <ai-chat-suggestions v-if="emptyChat" @select="applySuggestion" />
    </c-layout>
  </c-layout>
</template>

<script>
import { toRefs } from 'vue';

import { useAiChat } from '@/hooks/ai/ai-chat';
import { useForm } from '@/hooks/form';

import AiChatGreeting from '@/components/other/llm/chat/ai-chat-greeting.vue';
import AiChatInfoAlert from '@/components/other/llm/chat/ai-chat-info-alert.vue';
import AiChatMessagesList from '@/components/other/llm/chat/ai-chat-messages-list.vue';
import AiChatPatternItemField from '@/components/other/llm/chat/ai-chat-pattern-item-field.vue';
import AiChatSuggestions from '@/components/other/llm/chat/ai-chat-suggestions.vue';
import AiChatTextarea from '@/components/other/llm/chat/ai-chat-textarea.vue';
import AiChatWarningAlert from '@/components/other/llm/chat/ai-chat-warning-alert.vue';

export default {
  components: {
    AiChatMessagesList,
    AiChatGreeting,
    AiChatSuggestions,
    AiChatTextarea,
    AiChatInfoAlert,
    AiChatPatternItemField,
    AiChatWarningAlert,
  },
  emits: ['update:patterns', 'update:pending', 'update:json-string'],
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
    useForm();

    return useAiChat(toRefs(props), emit);
  },
};
</script>

<style lang="scss" scoped>
.ai-chat {
  --header-height: 64px;

  height: calc(100% - var(--header-height));
  min-height: 0;
  overflow: hidden;

  > :deep(.v-alert),
  > .ai-chat__footer {
    flex: 0 0 auto;
  }

  &__body {
    flex: 1 1 auto;
    min-height: 0;
    overflow-y: auto;
  }
}
</style>
