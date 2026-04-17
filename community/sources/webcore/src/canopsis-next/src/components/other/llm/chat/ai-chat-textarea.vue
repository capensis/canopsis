<template>
  <div
    :class="{ 'ai-chat-textarea--empty-chat': emptyChat }"
    class="ai-chat-textarea"
  >
    <c-alert :value="!!errorMessage" type="error">
      <span v-html="sanitizedErrorMessage" class="font-weight-regular" />
    </c-alert>
    <v-textarea
      v-field="prompt"
      v-validate="promptRules"
      ref="textareaElement"
      :placeholder="$t('llm.chat.promptPlaceholder')"
      :aria-label="$t('llm.chat.promptPlaceholder')"
      :auto-grow="emptyChat"
      :disabled="thinking || disabled"
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
        :value="llm"
        :items="llms"
        :pending="llmsPending"
        :disabled="llmsDisabled || disabled"
        @input="updateLlm"
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
        :disabled="askDisabled || disabled"
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

import AiChatLlmField from './partials/ai-chat-llm-field.vue';

export default {
  inject: ['$validator'],
  components: {
    AiChatLlmField,
  },
  model: {
    prop: 'prompt',
    event: 'input',
  },
  props: {
    prompt: {
      type: String,
      default: '',
    },
    llm: {
      type: Object,
      default: null,
    },
    llms: {
      type: Array,
      default: () => [],
    },
    llmsPending: {
      type: Boolean,
      default: false,
    },
    thinking: {
      type: Boolean,
      default: false,
    },
    emptyChat: {
      type: Boolean,
      default: false,
    },
    disabled: {
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

    const llmsDisabled = computed(() => props.llms.length <= 1 || !props.emptyChat);
    const askDisabled = computed(() => (
      props.llmsPending
      || !props.llms?.length
      || !props.prompt.trim()
    ));

    const sanitizedErrorMessage = computed(() => sanitizeHtml(props.errorMessage ?? ''));

    const promptRules = computed(() => ({
      required: true,
      max: LLM_AI_CHAT_PROMPT_MAX_LENGTH,
    }));

    /**
     * Forwards the selected LLM from `ai-chat-llm-field` to the parent via `update:llm`.

     * @param {Object|null} newLlm - Selected LLM record (same shape as the `llm` prop), or null when cleared.
     */
    const updateLlm = newLlm => emit('update:llm', newLlm);

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
        llm: props.llm?._id,
        prompt: props.prompt.trim(),
      });
    };

    /**
     * Emits `stop` so the parent can abort the current LLM request while `thinking` is true.
     */
    const stop = () => emit('stop');

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

    const focus = () => textareaElement.value?.focus?.();

    return {
      textareaElement,

      llmsDisabled,
      askDisabled,
      sanitizedErrorMessage,
      promptRules,

      updateLlm,
      ask,
      enterKeydown,
      stop,

      /**
       * Passed from the parent to focus the textarea element.
       */
      focus,
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
