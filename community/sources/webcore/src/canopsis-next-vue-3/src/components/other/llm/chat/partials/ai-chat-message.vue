<template>
  <c-layout
    :class="wrapperClass"
    class="ai-chat-message"
    gap="3"
  >
    <div class="ai-chat-message__content pa-4 flex-grow-1">
      <span
        v-if="message.thinking"
        class="primary--text ai-chat-message__thinking"
      >{{ $t('llm.chat.thinking') }}</span>
      <c-layout
        v-else
        class="ai-chat-message__body"
        gap="2"
        align="start"
      >
        <c-layout
          class="ai-chat-message__main"
          gap="3"
          align="start"
          column
        >
          <div
            v-if="sanitizedPrompt"
            :class="{ 'ai-chat-message__text--collapsed': isCollapsible && !isExpanded }"
            class="ai-chat-message__text"
          >
            <span v-html="sanitizedPrompt" class="pre-wrap" />
          </div>
          <span v-if="message.timestamp" class="ai-chat-message__timestamp">
            {{ timestampString }}
          </span>
        </c-layout>
        <c-expand-btn
          v-if="isCollapsible"
          v-model="isExpanded"
          color="grey darken-1"
          class="ai-chat-message__expand-btn"
        />
      </c-layout>
    </div>
    <v-avatar v-if="!isNoRole" :color="isUserRole ? 'success' : 'primary'" size="40">
      <v-icon
        v-if="isUserRole"
        color="white"
      >
        person
      </v-icon>
      <v-icon
        v-else
        color="white"
      >
        $vuetify.icons.robot
      </v-icon>
    </v-avatar>
  </c-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { DATETIME_FORMATS, LLM_AI_CHAT_MESSAGE_ROLES } from '@/constants';

import { convertDateToString, convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { sanitizeHtml } from '@/helpers/html';

import { useI18n } from '@/hooks/i18n';

const COLLAPSE_PROMPT_MIN_LENGTH = 100;

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
    const { t } = useI18n();

    const isExpanded = ref(false);

    const isUserRole = computed(() => props.message.role === LLM_AI_CHAT_MESSAGE_ROLES.user);
    const isModelRole = computed(() => props.message.role === LLM_AI_CHAT_MESSAGE_ROLES.model);
    const isNoRole = computed(() => !props.message.role);

    const isCollapsible = computed(() => (
      isUserRole.value && (props.message.prompt || '').length > COLLAPSE_PROMPT_MIN_LENGTH
    ));
    const wrapperClass = computed(() => ({
      'ai-chat-message--from-no-role': isNoRole.value,
      'ai-chat-message--from-user': isUserRole.value,
      'ai-chat-message--from-assistant': isModelRole.value,
      'ai-chat-message--thinking': props.message.thinking,
    }));

    const sanitizedPrompt = computed(() => {
      if (props.message.error) {
        return sanitizeHtml(props.message.error);
      }

      if (props.message.val_errors) {
        const validationMessage = `<ul>${props.message.val_errors.map(validationError => `<li>${validationError}</li>`).join('')}</ul>`;

        return `${t('llm.chat.patternCannotBeCreatedReasons')}${validationMessage}`;
      }

      if (props.message.prompt) {
        return sanitizeHtml(props.message.prompt);
      }

      if (props.message.patterns) {
        return props.message.emptyPatterns ? t('llm.chat.patternCreatedMessage') : t('llm.chat.patternUpdatedMessage');
      }

      return '';
    });

    const timestampString = computed(() => (
      props.history
        ? convertDateToString(props.message.timestamp, DATETIME_FORMATS.dateTimePicker)
        : convertDateToStringWithFormatForToday(props.message.timestamp, DATETIME_FORMATS.dateTimePicker, '', DATETIME_FORMATS.timePicker)
    ));

    return {
      isExpanded,
      isUserRole,
      isNoRole,
      isCollapsible,
      wrapperClass,
      sanitizedPrompt,
      timestampString,
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
  --collapsed-line-clamp: 4;

  flex-direction: row;
  padding-left: 16px;

  :deep(svg) {
    max-height: 24px;
  }

  &__thinking {
    animation: ai-chat-message-thinking-pulse 1.2s ease-in-out infinite;
  }

  &__content {
    border-radius: 8px;
  }

  &__body {
    width: 100%;
  }

  &__main {
    min-width: 0;
    flex-grow: 1;
  }

  &__text {
    min-width: 0;
    word-break: break-word;

    &--collapsed {
      display: -webkit-box;
      overflow: hidden;
      -webkit-box-orient: vertical;
      line-clamp: var(--collapsed-line-clamp);
      -webkit-line-clamp: var(--collapsed-line-clamp);
    }
  }

  &__timestamp {
    color: var(--v-grey-darken1, #757575) !important;
  }

  &__expand-btn {
    margin-left: auto;
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
