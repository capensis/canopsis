<template>
  <div class="ai-chat-sidebar">
    <pattern-progress
      v-if="pending"
      :in-progress-text="pendingTexts.inProgress"
      :cancel-button-label="pendingTexts.cancel"
      :style="progressStyle"
      class="ai-chat-sidebar__progress"
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
import {
  ref,
  watch,
  computed,
  onBeforeUnmount,
} from 'vue';

import { getMaxZIndex } from '@/helpers/vuetify';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';
import AiChatSidebarWrapper from '@/components/other/llm/chat/ai-chat-sidebar-wrapper.vue';
import AiChat from '@/components/other/llm/chat/ai-chat.vue';

const AI_CHAT_PENDING_CLASS = 'ai-chat-sidebar--pending';

export default {
  components: {
    PatternProgress,
    AiChatSidebarWrapper,
    AiChat,
  },
  emits: ['update:pending', 'update:minimized', 'update:patterns', 'update:json-string'],
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
    const progressStyle = {
      zIndex: getMaxZIndex() + 1,
    };

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

    watch(
      () => props.pending,
      (pending) => {
        document.body.classList.toggle(AI_CHAT_PENDING_CLASS, pending);
      },
      { immediate: true },
    );

    onBeforeUnmount(() => document.body.classList.remove(AI_CHAT_PENDING_CLASS));

    return {
      aiChatElement,
      restartable,
      progressStyle,

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

<style lang="scss">
.ai-chat-sidebar {
  display: block;

  &__progress.c-progress-overlay {
    position: fixed;
    inset: 0;
  }
}

body.ai-chat-sidebar--pending {
  overflow: hidden;
}

body.ai-chat-sidebar--pending .modal-wrapper .v-card-text {
  overflow: hidden;
}
</style>
