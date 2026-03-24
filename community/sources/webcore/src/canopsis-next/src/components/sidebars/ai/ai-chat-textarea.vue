<template>
  <div class="ai-chat-textarea pt-1">
    <v-textarea
      v-model="prompt"
      :placeholder="$t('llm.chat.promptPlaceholder')"
      :aria-label="$t('llm.chat.promptPlaceholder')"
      :auto-grow="emptyChat"
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
        :disabled="askDisabled"
        :aria-label="$t('llm.chat.ask')"
        color="primary"
        class="text-uppercase"
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

    return {
      prompt,
      selectedLlm,
      llms,
      llmsPending,
      llmsDisabled,
      askDisabled,
      ask,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat-textarea {
  border: 1px solid var(--v-divider-border-color, rgba(0, 0, 0, 0.12));
  border-radius: 8px;

  ::v-deep textarea {
    resize: none;
  }
}
</style>
