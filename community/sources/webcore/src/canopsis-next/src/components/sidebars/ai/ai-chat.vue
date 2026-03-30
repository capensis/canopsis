<template>
  <v-layout class="ai-chat" column>
    <ai-chat-warning-alert />
    <div
      ref="bodyElement"
      class="ai-chat__body layout column gap-6 pa-4"
    >
      <ai-chat-greeting v-if="emptyChat" />
      <v-layout
        class="ai-chat__messages gap-3"
        column
      >
        <template v-for="message in messages">
          <ai-chat-pattern
            v-if="message.patterns"
            :key="message._id"
            :version="message.version"
            :active="message.version === activeVersion"
            :patterns="message.patterns"
            :original-version="message.originalVersion"
            @restore:version="restoreVersion"
          />
          <ai-chat-message
            v-else
            :key="`${message._id}-message`"
            :message="message"
          />
        </template>
        <ai-chat-message
          v-if="thinking"
          :key="thinkingMessage._id"
          :message="thinkingMessage"
        />
      </v-layout>
    </div>
    <v-layout
      :justify-end="!emptyChat"
      :class="{ 'pa-4': emptyChat }"
      class="gap-6 pt-0"
      column
    >
      <ai-chat-textarea
        ref="textareaElement"
        :error-message="errorMessage"
        :empty-chat="emptyChat"
        :thinking="thinking"
        @ask="ask"
        @stop="stop"
      />

      <ai-chat-suggestions v-if="emptyChat" @select="applySuggestion" />
    </v-layout>
    <portal v-if="socketRoom" :to="$constants.PORTALS_NAMES.additionalSidebarTitle">
      <v-btn
        class="white--text"
        color="white"
        outlined
        @click="restart"
      >
        <v-icon class="mr-2" color="white" small>
          refresh
        </v-icon>
        {{ $t('common.restart') }}
      </v-btn>
    </portal>
  </v-layout>
</template>

<script>
import { isEmpty, isUndefined } from 'lodash';
import {
  computed,
  ref,
  set,
  watch,
  nextTick,
  onBeforeUnmount,
} from 'vue';

import { SOCKET_ROOMS } from '@/config';
import { LLM_AI_CHAT_MESSAGE_ROLES, SIDE_BARS, PATTERNS_FIELDS } from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import { uid } from '@/helpers/uid';
/**
 * TODO: REMOVE IT
 */
import { promisedTimeout } from '@/helpers/async';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';
import { getChangedPatternsFields } from '@/helpers/entities/pattern/form';

import { useI18n } from '@/hooks/i18n';
import { useSocket } from '@/hooks/socket';

import AiChatGreeting from './ai-chat-greeting.vue';
import AiChatMessage from './ai-chat-message.vue';
import AiChatPattern from './ai-chat-pattern.vue';
import AiChatSuggestions from './ai-chat-suggestions.vue';
import AiChatTextarea from './ai-chat-textarea.vue';
import AiChatWarningAlert from './ai-chat-warning-alert.vue';

export default {
  name: SIDE_BARS.aiChat,
  components: {
    AiChatGreeting,
    AiChatMessage,
    AiChatPattern,
    AiChatSuggestions,
    AiChatTextarea,
    AiChatWarningAlert,
  },
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t, tc } = useI18n();

    const bodyElement = ref(null);
    const textareaElement = ref(null);
    const messages = ref([]);
    const thinking = ref(false);
    const versions = ref([]);
    const lastLlmVersion = ref(null);
    const socketRoom = ref(null);
    const errorMessage = ref(null);

    let lastPrompt = '';
    let lastLlmMessageIndex = 0;
    let lastChangedPatternsFieldsMessageIndex = 0;

    const activeVersion = computed(() => Math.max(versions.value.length - 1, 0));
    const currentFormPatterns = computed(() => (
      formFilterToPatterns(props.sidebar.config?.patterns ?? {}, Object.values(PATTERNS_FIELDS))
    ));
    const changedPatternsFields = computed(() => (
      getChangedPatternsFields(currentFormPatterns.value, versions.value[lastLlmVersion.value])
    ));
    const emptyCurrentFormPatterns = computed(() => isEmpty(currentFormPatterns.value));
    const emptyChat = computed(() => messages.value.length === 0);

    const socket = useSocket();

    const thinkingMessage = computed(() => ({
      _id: uid(),
      role: LLM_AI_CHAT_MESSAGE_ROLES.model,
      thinking: true,
    }));

    const addMessage = (text, role = LLM_AI_CHAT_MESSAGE_ROLES.model) => {
      messages.value.push({
        text,
        role,

        _id: uid(),
        component: 'ai-chat-message',
        time: Date.now(),
      });
    };

    const addPattern = (patterns, role, originalVersion = null) => {
      versions.value.push(patterns);

      let messageText = '';

      if (!role && changedPatternsFields.value.length > 0) {
        messageText = tc('llm.chat.patternsEditedMessage', changedPatternsFields.value.length, {
          patterns: changedPatternsFields.value.map(field => t(`pattern.patternsFields.${field}`)).join(', ').toLowerCase(),
        });
      } else if (role === LLM_AI_CHAT_MESSAGE_ROLES.model || !isUndefined(originalVersion)) {
        lastLlmVersion.value = activeVersion.value;
        messageText = emptyCurrentFormPatterns.value ? t('llm.chat.patternCreatedMessage') : t('llm.chat.patternUpdatedMessage');
      }

      if (messageText) {
        addMessage(
          messageText,
          role,
        );
      }

      messages.value.push({
        patterns,
        originalVersion,

        _id: uid(),
        component: 'ai-chat-pattern',
        version: activeVersion.value,
      });
    };

    const resetTextarea = () => {
      textareaElement.value.prompt = lastPrompt;
      lastPrompt = '';
    };

    /**
     * Invoked when the `llmchat` room delivers a server message; extend to update `messages` or streaming state.
     * Payload shape follows `SocketRoom.call` (API-defined).
     */
    const socketListener = async ({ error, code, ...patterns }) => {
      await promisedTimeout(() => {}, 5000);

      if (error) {
        thinking.value = false;
        errorMessage.value = error;

        resetTextarea();

        return;
      }

      addPattern(patterns, LLM_AI_CHAT_MESSAGE_ROLES.model);
      props.sidebar.config.setPatterns(patterns);

      thinking.value = false;
    };

    const scrollChatBodyToBottom = () => nextTick(() => requestAnimationFrame(() => {
      const el = bodyElement.value;

      if (!el) {
        return;
      }

      el.scrollTop = el.scrollHeight;
    }));

    watch(messages, scrollChatBodyToBottom, { deep: true, flush: 'post' });
    watch(thinking, scrollChatBodyToBottom, { flush: 'post' });

    const errorHandler = ({ message, error }) => {
      const { room, code } = error ?? {};

      if (room !== SOCKET_ROOMS.llmChat) {
        return;
      }

      thinking.value = false;
      errorMessage.value = code ? t(`llm.chat.errors.${code}`) : message;

      resetTextarea();
    };

    /**
     * Joins `SOCKET_ROOMS.llmChat` merging `socketRoomData` from the sidebar with `config: llm`.
     * Registers `socketListener` on the resulting `SocketRoom`.
     *
     * @param {Object} llm - LLM config merged into the join payload as `config` (e.g. selected model row).
     */
    const joinSocketRoom = (llm) => {
      const joinData = {
        ...props.sidebar.config?.socketRoomData,

        config: llm,
      };

      socketRoom.value = socket
        .on(Socket.EVENTS_TYPES.error, errorHandler)
        .join(SOCKET_ROOMS.llmChat, joinData, true)
        .addListener(socketListener);
    };

    /**
     * Leaves the LLM chat socket room and clears the cached `SocketRoom` reference (runs on unmount).
     */
    const leaveSocketRoom = () => {
      socket
        .off(Socket.EVENTS_TYPES.error, errorHandler)
        .leave(SOCKET_ROOMS.llmChat)
        .removeListener(socketListener);

      socketRoom.value = null;
    };

    /**
     * Ensures the chat room is joined, then sends the user prompt over the socket via `SocketRoom.send`.
     *
     * @param {Object} params
     * @param {Object} params.llm - LLM config passed into `joinSocketRoom` when the room is not open yet.
     * @param {string} params.prompt - User message text sent in the outbound `{ prompt }` payload.
     */
    const ask = ({ llm, prompt }) => {
      if (!socketRoom.value) {
        joinSocketRoom(llm);
      }

      errorMessage.value = null;
      lastPrompt = prompt;
      textareaElement.value.prompt = '';
      thinking.value = true;

      const data = {
        prompt,
        type: 'send',
      };

      if (!emptyCurrentFormPatterns.value) {
        Object.entries(currentFormPatterns.value).forEach(([field, pattern]) => {
          data[field] = pattern;
        });
      }

      socketRoom.value.send(data);

      addMessage(prompt, LLM_AI_CHAT_MESSAGE_ROLES.user);
    };

    /**
     * Sends `{ type: 'cancel' }` on the room channel to stop an in-flight request when the room is joined.
     */
    const stop = () => {
      socketRoom.value?.send({ type: 'cancel' });

      textareaElement.value.prompt = lastPrompt;
      thinking.value = false;
      lastPrompt = '';
    };

    /**
     * Handles `select` from `ai-chat-suggestions` (payload: `t('llm.chat.suggestionPrompts.*')` for current locale).
     */
    const applySuggestion = (newPrompt) => {
      textareaElement.value.prompt = newPrompt;
      textareaElement.value.$refs.textareaElement?.focus?.();
    };

    const restoreVersion = (version) => {
      const newPatterns = versions.value[version];

      addPattern(newPatterns, LLM_AI_CHAT_MESSAGE_ROLES.user, version);

      props.sidebar.config.setPatterns(newPatterns);
    };

    watch(currentFormPatterns, (newPatterns) => {
      if (emptyChat.value) {
        return;
      }

      if (changedPatternsFields.value.length === 0) {
        if (lastLlmVersion.value !== activeVersion.value) {
          messages.value.splice(lastLlmMessageIndex + 1);
          versions.value.pop();
        }

        return;
      }

      if (lastLlmVersion.value !== activeVersion.value) {
        set(messages.value, lastChangedPatternsFieldsMessageIndex, {
          ...messages.value[lastChangedPatternsFieldsMessageIndex],
          text: tc('llm.chat.patternsEditedMessage', changedPatternsFields.value.length, {
            patterns: changedPatternsFields.value.map(field => t(`pattern.patternsFields.${field}`)).join(', ').toLowerCase(),
          }),
        });

        set(messages.value, lastChangedPatternsFieldsMessageIndex + 1, {
          ...messages.value[lastChangedPatternsFieldsMessageIndex + 1],
          patterns: newPatterns,
        });

        return;
      }

      lastLlmMessageIndex = messages.value.length - 1;
      lastChangedPatternsFieldsMessageIndex = lastLlmMessageIndex + 1;

      addPattern(newPatterns);
    });

    watch(
      thinking,
      newThinking => props.sidebar.config.setPending?.(newThinking, emptyCurrentFormPatterns.value, stop),
    );

    onBeforeUnmount(leaveSocketRoom);

    return {
      bodyElement,
      textareaElement,
      versions,
      activeVersion,
      messages,
      thinking,
      thinkingMessage,
      socketRoom,
      emptyChat,
      errorMessage,
      ask,
      stop,
      applySuggestion,
      restoreVersion,
      restart: leaveSocketRoom,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat {
  --header-height: 64px;
  height: calc(100% - var(--header-height));
  min-height: 0;

  &__body {
    flex: 1 1 auto;
    min-height: 0;
    overflow-y: auto;
  }

  /* Override Vuetify `.layout { flex: 1 1 auto }` so `padding-bottom` stays under the last message. */
  &__messages {
    flex: 0 0 auto !important;
    min-height: 0;

    > * {
      flex: 0 0 auto !important;
    }
  }
}
</style>
