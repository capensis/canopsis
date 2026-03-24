<template>
  <div
    :class="{ 'ai-chat-textarea--empty-chat': emptyChat }"
    class="ai-chat-textarea"
  >
    <v-textarea
      v-model="prompt"
      :placeholder="$t('llm.chat.promptPlaceholder')"
      :aria-label="$t('llm.chat.promptPlaceholder')"
      :auto-grow="emptyChat"
      :disabled="asking"
      rows="5"
      solo
      flat
      hide-details
    />
    <v-layout
      class="px-4 pb-4 gap-2"
      wrap
      align-end
      justify-space-between
    >
      <ai-chat-llm-field
        v-model="selectedLlm"
        :items="llms"
        :pending="llmsPending"
        :disabled="llmsDisabled"
      />
      <v-btn
        v-if="asking"
        color="secondary"
        depressed
        @click="stop"
      >
        <v-icon class="mr-2">
          stop
        </v-icon>
        <span>{{ $t('common.stop') }}</span>
      </v-btn>
      <v-btn
        v-else
        :disabled="askDisabled"
        color="primary"
        depressed
        @click="ask"
      >
        {{ $t('llm.chat.ask') }}
      </v-btn>
    </v-layout>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import { useAiChatLlms } from './hooks/use-ai-chat-llms';
import AiChatLlmField from './ai-chat-llm-field.vue';

export default {
  components: {
    AiChatLlmField,
  },
  props: {
    asking: {
      type: Boolean,
      default: false,
    },
    emptyChat: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const prompt = ref('');
    const selectedLlm = ref(null);

    const { llms, llmsPending } = useAiChatLlms();

    const llmsDisabled = computed(() => llms.value.length <= 1);
    const askDisabled = computed(() => (
      llmsPending.value
      || !prompt.value.trim()
    ));

    /**
     * Emits the user message and chosen model so the parent can run the AI request.
     *
     * Payload: `{ model: string | null, prompt: string }` (`prompt` is trimmed).
     */
    const ask = () => emit('ask', {
      model: selectedLlm.value,
      prompt: prompt.value.trim(),
    });

    /**
     * Emits `stop` so the parent can abort the current LLM request while `asking` is true.
     */
    const stop = () => emit('stop');

    return {
      prompt,
      selectedLlm,
      llms,
      llmsPending,
      llmsDisabled,
      askDisabled,
      ask,
      stop,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat-textarea {
  border-top: 1px solid var(--v-divider-border-color, rgba(0, 0, 0, 0.12));

  ::v-deep textarea {
    resize: none;
  }

  &--empty-chat {
    padding-top: 8px;
    border-radius: 8px;
    border: 1px solid var(--v-divider-border-color, rgba(0, 0, 0, 0.12));
  }
}
</style>
