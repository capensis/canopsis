<template>
  <v-layout class="ai-chat" column>
    <v-layout
      class="ai-chat__body gap-6 pa-4"
      column
    >
      <ai-chat-greeting />
      <v-layout
        class="gap-3"
        column
      >
        <ai-chat-message
          :text="demoUserMessageText"
          :timestamp="demoUserMessageTime"
          from-user
        />
        <ai-chat-message
          :text="'Pattern created'"
          :timestamp="demoUserMessageTime"
        />
        <ai-chat-message
          :text="'Pattern created'"
          :timestamp="demoUserMessageTime"
          from-system
        />
        <ai-chat-pattern :pattern="testPattern" />
        <ai-chat-pattern :pattern="testPattern" active />

        <ai-chat-message thinking />
      </v-layout>
      <ai-chat-textarea
        :selected-model="selectedModel"
        @update:selected-model="updateSelectedModel"
        @ask="forwardAsk"
      />

      <ai-chat-suggestions @select="applySuggestion" />
    </v-layout>
  </v-layout>
</template>

<script>
import { ref } from 'vue';

import { SIDE_BARS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

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
  setup(_, { emit }) {
    const { t } = useI18n();

    const demoUserMessageText = 'Create pattern when v.component has "test-attribute-correlation-1" prefix and when entity name is "test-entity-1" and last event date in range 2025-01-01T00:00:00Z to 2025-03-01T00:05:00Z.';
    const demoUserMessageTime = 12345615965;

    const prompt = ref('');
    const selectedModel = ref('');
    const updatePrompt = (value) => {
      prompt.value = value;
    };

    const updateSelectedModel = (value) => {
      selectedModel.value = value;
    };

    const forwardAsk = (payload) => {
      emit('ask', payload);
    };

    const suggestionPromptKeys = {
      createPattern: 'llm.chat.suggestionPrompts.createPattern',
      editPattern: 'llm.chat.suggestionPrompts.editPattern',
      validatePattern: 'llm.chat.suggestionPrompts.validatePattern',
    };

    const applySuggestion = (type) => {
      const key = suggestionPromptKeys[type];

      if (key) {
        prompt.value = t(key);
      }

      emit('suggestion', { type });
    };

    return {
      demoUserMessageText,
      demoUserMessageTime,
      prompt,
      selectedModel,
      updatePrompt,
      updateSelectedModel,
      forwardAsk,
      applySuggestion,
      testPattern: {
        alarm_pattern: [
          [
            {
              field: 'v.component',
              cond: {
                type: 'eq',
                value: 'test-attribute-correlation-1',
              },
            },
          ],
        ],
      },
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat {
  --header-height: 64px;
  height: calc(100% - var(--header-height));

  &__body {
    flex: 1 1 auto;
  }

}
</style>
