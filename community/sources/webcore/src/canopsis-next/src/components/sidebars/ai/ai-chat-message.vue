<template>
  <v-layout
    :class="wrapperClass"
    class="ai-chat-message gap-3"
  >
    <v-flex class="ai-chat-message__content pa-4">
      <span
        v-if="message.thinking"
        class="primary--text ai-chat-message__thinking"
      >{{ $t('llm.chat.thinking') }}</span>
      <v-layout
        v-else
        class="gap-2"
        align-start
      >
        <v-layout
          class="gap-3"
          column
          align-start
        >
          <v-flex
            v-if="message.text"
            :class="{ 'ai-chat-message__text--collapsed': isCollapsible && !isExpanded }"
            class="ai-chat-message__text"
          >
            {{ message.text }}
          </v-flex>
          <span v-if="message.time" class="grey--text text--darken-1">
            {{ timeString }}
          </span>
        </v-layout>
        <c-expand-btn
          v-if="isCollapsible"
          v-model="isExpanded"
          color="grey darken-1"
          class="ai-chat-message__expand-btn"
        />
      </v-layout>
    </v-flex>
    <v-avatar v-if="!isNoRole" :color="isUserRole ? 'success' : 'primary'" size="40">
      <v-icon
        v-if="isUserRole"
        dark
      >
        person
      </v-icon>
      <v-icon
        v-else
        dark
      >
        $vuetify.icons.robot
      </v-icon>
    </v-avatar>
  </v-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { DATETIME_FORMATS, LLM_AI_CHAT_MESSAGE_ROLES } from '@/constants';

import { convertDateToString, convertDateToStringWithFormatForToday } from '@/helpers/date/date';

const COLLAPSE_TEXT_MIN_LENGTH = 100;

export default {
  props: {
    message: {
      type: Object,
      default: () => ({}),
    },
    history: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const isExpanded = ref(false);

    const isUserRole = computed(() => props.message.role === LLM_AI_CHAT_MESSAGE_ROLES.user);
    const isModelRole = computed(() => props.message.role === LLM_AI_CHAT_MESSAGE_ROLES.model);
    const isNoRole = computed(() => !props.message.role);

    const isCollapsible = computed(() => (
      isUserRole.value && (props.message.text || '').length > COLLAPSE_TEXT_MIN_LENGTH
    ));
    const wrapperClass = computed(() => ({
      'ai-chat-message--from-no-role': isNoRole.value,
      'ai-chat-message--from-user': isUserRole.value,
      'ai-chat-message--from-assistant': isModelRole.value,
      'ai-chat-message--thinking': props.message.thinking,
    }));

    const timeString = computed(() => (
      props.history
        ? convertDateToString(props.message.time, DATETIME_FORMATS.long)
        : convertDateToStringWithFormatForToday(props.message.time, DATETIME_FORMATS.long, '', DATETIME_FORMATS.timeWithoutSeconds)
    ));

    return {
      isExpanded,
      isUserRole,
      isNoRole,
      isCollapsible,
      wrapperClass,
      timeString,
    };
  },
};
</script>

<style lang="scss" scoped>
@keyframes ai-chat-message-thinking-pulse {
  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.35;
  }
}

.ai-chat-message {
  flex-direction: row;
  padding-left: 16px;

  ::v-deep svg {
    max-height: 24px;
  }

  &__thinking {
    animation: ai-chat-message-thinking-pulse 1.2s ease-in-out infinite;
  }

  &__content {
    border-radius: 8px;
  }

  &__text {
    min-width: 0;
    word-break: break-word;

    &--collapsed {
      display: -webkit-box;
      overflow: hidden;
      -webkit-box-orient: vertical;
      line-clamp: 2;
      -webkit-line-clamp: 2;
    }
  }

  &__expand-btn {
    opacity: .5;
  }

  &--from-user {
    .ai-chat-message__content {
      background-color: var(--v-success-background-base);
    }
  }

  &--from-assistant {
    flex-direction: row-reverse;
    padding-right: 16px;

    .ai-chat-message__content {
      background-color: #E6E6E6;
    }
  }

  &--from-no-role {
    .ai-chat-message__content {
      border: 1px solid #949494;
    }
  }

  &--thinking {
    flex-direction: row-reverse;
    align-items: center;

    .ai-chat-message__content {
      padding: 0 !important;
      background-color: transparent;
    }
  }
}

@media (prefers-reduced-motion: reduce) {
  .ai-chat-message__thinking {
    animation: none;
  }
}
</style>
