<template>
  <div>
    <pattern-progress
      v-if="pending"
      :in-progress-text="pendingTexts.inProgress"
      :cancel-button-label="pendingTexts.cancel"
      @cancel="cancel"
    />
    <ai-chat-sidebar-wrapper
      :initial-minimized="initialMinimized"
      :restartable="restartable"
      @restart="restart"
      @update:minimized="updateMinimized"
    >
      <ai-chat
        ref="aiChatElement"
        :json-string="jsonString"
        :patterns="patterns"
        :patterns-items="patternsItems"
        :llms="llms"
        :context="context"
        :rule-id="ruleId"
        @update:pending="updatePending"
        @update:patterns="updatePatterns"
        @update:json-string="updateJsonString"
      />
    </ai-chat-sidebar-wrapper>
  </div>
</template>

<script>
import { ref, computed } from 'vue';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';
import AiChatSidebarWrapper from '@/components/other/llm/chat/ai-chat-sidebar-wrapper.vue';
import AiChat from '@/components/other/llm/chat/ai-chat.vue';

export default {
  components: {
    PatternProgress,
    AiChatSidebarWrapper,
    AiChat,
  },
  props: {
    initialMinimized: {
      type: Boolean,
      default: false,
    },
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
    pending: {
      type: Boolean,
      default: false,
    },
    pendingTexts: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const aiChatElement = ref(null);

    const restartable = computed(() => !aiChatElement.value?.emptyChat);

    /**
     * Opens the restart confirmation on the nested `ai-chat` (leaves socket room and resets session on confirm).
     */
    const restart = () => aiChatElement.value?.restart();

    /**
     * Sends cancel on the LLM socket and restores the prompt via the nested `ai-chat` `stop` action.
     */
    const cancel = () => aiChatElement.value?.stop();

    /**
     * Forwards `ai-chat` pending state to the parent (`update:pending`), same signature as `useAiChat` emit
     * (pending flag, empty-patterns flag, `stop` handler).
     *
     * @param {...*} args - Arguments from `@update:pending` on `ai-chat`.
     */
    const updatePending = (...args) => emit('update:pending', ...args);

    /**
     * Forwards minimized sidebar chrome to the parent (`update:minimized`).
     *
     * @param {boolean} newMinimized - Next minimized value from `ai-chat-sidebar-wrapper`.
     */
    const updateMinimized = newMinimized => emit('update:minimized', newMinimized);

    /**
     * Forwards applied patterns from `ai-chat` to the parent (`update:patterns`).
     *
     * @param {...*} args - Arguments from `@update:patterns` on `ai-chat`.
     */
    const updatePatterns = (...args) => emit('update:patterns', ...args);

    /**
     * Forwards pattern JSON pushes from `ai-chat` to the parent (`update:json-string`).
     *
     * @param {string} newJsonString - Value from `@update:json-string` on `ai-chat`.
     */
    const updateJsonString = newJsonString => emit('update:json-string', newJsonString);

    return {
      aiChatElement,
      restartable,

      restart,
      cancel,
      updatePending,
      updateMinimized,
      updatePatterns,
      updateJsonString,
    };
  },
};
</script>
