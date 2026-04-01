<template>
  <div
    :class="{ 'ai-chat-textarea--empty-chat': emptyChat }"
    class="ai-chat-textarea"
  >
    <c-alert :value="!!errorMessage" type="error">
      <span v-html="sanitizedErrorMessage" class="font-weight-regular" />
    </c-alert>
    <v-textarea
      v-model="prompt"
      v-validate="promptRules"
      ref="textareaElement"
      :placeholder="$t('llm.chat.promptPlaceholder')"
      :aria-label="$t('llm.chat.promptPlaceholder')"
      :auto-grow="emptyChat"
      :disabled="thinking"
      :error-messages="errors.collect('prompt')"
      name="prompt"
      rows="5"
      solo
      flat
      hide-details
      @keydown.enter="enterKeydown"
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
        v-if="thinking"
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

import { LLM_AI_CHAT_PROMPT_MAX_LENGTH } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';

import { useValidator } from '@/hooks/validator/validator';

import { useAiChatLlms } from './hooks/use-ai-chat-llms';
import AiChatLlmField from './ai-chat-llm-field.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AiChatLlmField,
  },
  props: {
    thinking: {
      type: Boolean,
      default: false,
    },
    emptyChat: {
      type: Boolean,
      default: false,
    },
    errorMessage: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const validator = useValidator();

    const textareaElement = ref(null);
    const prompt = ref('');
    const selectedLlm = ref(null);

    const { llms, llmsPending } = useAiChatLlms();

    const llmsDisabled = computed(() => llms.value.length <= 1 || !props.emptyChat);
    const askDisabled = computed(() => (
      llmsPending.value
      || !prompt.value.trim()
    ));

    const sanitizedErrorMessage = computed(() => sanitizeHtml(props.errorMessage ?? ''));

    const promptRules = computed(() => ({
      required: true,
      max: LLM_AI_CHAT_PROMPT_MAX_LENGTH,
    }));

    /**
     * Emits the user message and chosen model so the parent can run the AI request.
     *
     * Payload: `{ llm: string | null, prompt: string }` (`prompt` is trimmed).
     */
    const ask = async () => {
      const isValid = await validator.validateAll();

      if (!isValid) {
        return;
      }

      emit('ask', {
        llm: selectedLlm.value?._id,
        prompt: prompt.value,
      });
    };

    /**
     * Submits on Enter; Shift+Enter keeps the default newline. Matches `Ask` enablement.
     *
     * @param {KeyboardEvent} event
     */
    const enterKeydown = (event) => {
      if (event.shiftKey || askDisabled.value) {
        return;
      }

      event.preventDefault();
      ask();
    };

    /**
     * Emits `stop` so the parent can abort the current LLM request while `thinking` is true.
     */
    const stop = () => emit('stop');

    return {
      textareaElement,
      prompt,
      selectedLlm,
      llms,
      llmsPending,
      llmsDisabled,
      askDisabled,
      sanitizedErrorMessage,
      promptRules,
      ask,
      enterKeydown,
      stop,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat-textarea {
  --llm-ai-chat-textarea-max-height: 50vh;

  border-top: 1px solid var(--v-divider-border-color, rgba(0, 0, 0, 0.12));

  ::v-deep textarea {
    max-height: var(--llm-ai-chat-textarea-max-height);
    overflow-y: auto;
    resize: none;
  }

  &--empty-chat {
    padding-top: 8px;
    border-radius: 8px;
    border: 1px solid var(--v-divider-border-color, rgba(0, 0, 0, 0.12));
  }
}
</style>
