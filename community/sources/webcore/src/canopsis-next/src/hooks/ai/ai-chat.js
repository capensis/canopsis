import { isEmpty, isUndefined, isNil } from 'lodash';
import {
  computed,
  ref,
  unref,
  set,
  watch,
  nextTick,
  onBeforeUnmount,
} from 'vue';

import { SOCKET_ROOMS } from '@/config';
import {
  MAX_LIMIT,
  LLM_AI_CHAT_MESSAGE_ROLES,
  LLM_AI_CHAT_MESSAGE_TYPES,
  LLM_AI_CHAT_ERROR_CODES,
  PATTERNS_FIELDS,
  DATETIME_FORMATS,
} from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import { uid } from '@/helpers/uid';
/**
 * TODO: REMOVE IT
 */
import { promisedTimeout } from '@/helpers/async';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';
import { getChangedPatternsFields } from '@/helpers/entities/pattern/form';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { useSocket } from '@/hooks/socket';
import { usePendingHandler } from '@/hooks/query/pending';
import { useLlm } from '@/hooks/store/modules/llm';

const CHAT_COMPONENTS = {
  message: 'ai-chat-message',
  pattern: 'ai-chat-pattern',
};

/**
 * Hook for watching sidebar JSON string changes and triggering `ask` with a localized fix-pattern prompt.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref|Object} [options.sidebar] - Sidebar instance; `config.jsonString` is observed.
 * @param {function(Object): void} [options.ask] - Invoked with `{ prompt, withoutPatterns }` when JSON changes.
 * @returns {{ jsonString: import('vue').ComputedRef<string> }}
 */
export const useAiChatJsonString = ({ sidebar, ask } = {}) => {
  const { t } = useI18n();

  const jsonString = computed(() => unref(sidebar)?.config?.jsonString ?? '');

  watch(jsonString, newJsonString => ask({
    prompt: t('llm.chat.fixPatternPrompt', { jsonString: newJsonString }),
    withoutPatterns: true,
  }));

  return {
    jsonString,
  };
};

/**
 * Hook for the stack of pattern versions from the chat and the index of the last LLM-applied version.
 *
 * @returns {Object} Version refs and mutators (`addVersion`, `removeLastVersion`, `resetVersions`).
 */
export const useAiChatVersions = () => {
  const versions = ref([]);
  const lastLlmVersion = ref(null);

  const activeVersion = computed(() => versions.value.length - 1);
  const lastLlmPatterns = computed(() => (isNil(lastLlmVersion.value) ? null : versions.value[lastLlmVersion.value]));

  /**
   * Appends a pattern snapshot to `versions` and optionally marks it as the latest LLM-produced index.
   *
   * @param {Object} patterns - Pattern payload stored for this version.
   * @param {boolean} [isLlm=false] - When true, sets `lastLlmVersion` to the new `activeVersion`.
   */
  const addVersion = (patterns, isLlm = false) => {
    versions.value.push(patterns);

    if (isLlm) {
      lastLlmVersion.value = activeVersion.value;
    }
  };

  /**
   * Drops the most recent entry from `versions`.
   */
  const removeLastVersion = () => versions.value.pop();

  /**
   * Clears all versions and does not reset `lastLlmVersion` (callers may set it to `null` if needed).
   */
  const resetVersions = () => {
    versions.value = [];
    lastLlmVersion.value = null;
  };

  return {
    versions,
    lastLlmVersion,
    activeVersion,
    lastLlmPatterns,

    addVersion,
    removeLastVersion,
    resetVersions,
  };
};

/**
 * Hook for the chat message list, pattern cards, and version bookkeeping tied to sidebar patterns.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref|Object} [options.sidebar] - Used for `setPatterns` on restore.
 * @param {import('vue').ComputedRef|import('vue').Ref} [options.currentFormPatterns]
 *   Live form patterns for diff vs last LLM output.
 * @returns {Object} Messages API, versions, diffs, `restoreVersion`.
 */
export const useAiChatMessages = ({ sidebar, currentFormPatterns } = {}) => {
  const { t, tc } = useI18n();
  const {
    versions,
    lastLlmVersion,
    activeVersion,
    addVersion,
    removeLastVersion,
    resetVersions,
  } = useAiChatVersions();

  const lastLlmPatterns = computed(() => (isNil(lastLlmVersion.value) ? null : versions.value[lastLlmVersion.value]));

  const changedPatternsFields = computed(() => (
    getChangedPatternsFields(unref(currentFormPatterns), unref(lastLlmPatterns))
  ));

  const emptyCurrentFormPatterns = computed(() => isEmpty(unref(currentFormPatterns)));

  const messages = ref([]);

  /**
   * Appends a text or pattern chat row with `_id`, `time`, and `CHAT_COMPONENTS.message` tag; extra keys are merged.
   *
   * @param {Object} payload
   * @param {string} [payload.text] - Visible body for a text bubble.
   * @param {string} [payload.role] - Defaults to `LLM_AI_CHAT_MESSAGE_ROLES.model`.
   * @returns {number} New `messages` length from `Array.prototype.push`.
   */
  const addMessage = ({ text, role, ...rest }) => messages.value.push({
    text,
    role,

    _id: uid(),
    component: CHAT_COMPONENTS.message,
    time: Date.now(),

    ...rest,
  });

  /**
   * Pushes a pattern version, optionally a status text line, then a pattern card message.
   *
   * @param {Object} options
   * @param {Object} options.patterns - Pattern payload for this version.
   * @param {string} [options.role] - `LLM_AI_CHAT_MESSAGE_ROLES.model` for LLM turns; omit when syncing user edits.
   * @param {number} [options.originalVersion] - Set when restoring a prior version (skips extra status text).
   */
  const addPattern = ({ patterns, role, originalVersion }) => {
    let messageText = '';

    if (!role && changedPatternsFields.value.length > 0) {
      messageText = tc('llm.chat.patternsEditedMessage', changedPatternsFields.value.length, {
        patterns: changedPatternsFields.value.map(field => t(`pattern.patternsFields.${field}`)).join(', ').toLowerCase(),
      });

      addVersion(patterns);
    } else if (role === LLM_AI_CHAT_MESSAGE_ROLES.model) {
      addVersion(patterns, true);

      messageText = emptyCurrentFormPatterns.value ? t('llm.chat.patternCreatedMessage') : t('llm.chat.patternUpdatedMessage');
    }

    if (messageText && isUndefined(originalVersion)) {
      addMessage({
        role,
        text: messageText,
      });
    }

    addMessage({
      patterns,
      role: LLM_AI_CHAT_MESSAGE_ROLES.model,
      originalVersion,
      version: activeVersion.value,

      _id: uid(),
      component: CHAT_COMPONENTS.pattern,
    });
  };

  /**
   * Replaces one message at `index` (Vue 2 reactivity-safe).
   *
   * @param {number} index - Index in `messages`.
   * @param {Object} newMessage - Full replacement message object.
   */
  const updateMessage = (index, newMessage) => set(messages.value, index, { ...messages.value[index], ...newMessage });

  /**
   * Drops the last `count` messages from the tail of the list.
   *
   * @param {number} [count=1] - How many entries to remove from the end.
   * @returns {Array} Elements removed by `splice`.
   */
  const removeLastMessages = (count = 1) => messages.value.splice(messages.value.length - count);

  /**
   * Replays patterns from `versions[version]` into chat and applies them to the sidebar form.
   *
   * @param {number} version - Version index inside `versions`.
   */
  const restoreVersion = (version) => {
    const newPatterns = versions.value[version];

    addPattern({ patterns: newPatterns, role: LLM_AI_CHAT_MESSAGE_ROLES.model, originalVersion: version });

    unref(sidebar)?.config?.setPatterns?.(newPatterns);
  };

  /**
   * Clears the chat message list (does not reset versions; use `resetVersions` separately if needed).
   */
  const resetMessages = () => messages.value = [];

  return {
    messages,
    addMessage,
    addPattern,
    updateMessage,
    removeLastMessages,
    resetMessages,

    versions,
    activeVersion,
    lastLlmVersion,
    addVersion,
    removeLastVersion,
    restoreVersion,
    resetVersions,

    changedPatternsFields,
    emptyCurrentFormPatterns,
  };
};

/**
 * Hook for binding sidebar filter patterns to chat: syncs manual edits with messages and version stack.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref|Object} [options.sidebar] - Source for `config.patterns` and `setPatterns`.
 * @returns {Object} Subset of message/version API for the AI chat UI (`resetMessages`, `resetVersions`, etc.).
 */
export const useAiChatPattern = ({ sidebar } = {}) => {
  const { t, tc } = useI18n();

  let lastChangedPatternsFieldsMessageIndex = -1;
  let lastLlmMessageIndex = -1;

  const currentFormPatterns = computed(() => (
    formFilterToPatterns(unref(sidebar)?.config?.patterns ?? {}, Object.values(PATTERNS_FIELDS))
  ));

  const {
    messages,
    addMessage,
    addPattern,
    updateMessage,
    removeLastMessages,
    resetMessages,

    versions,
    lastLlmVersion,
    activeVersion,
    removeLastVersion,
    resetVersions,
    restoreVersion,

    changedPatternsFields,
    emptyCurrentFormPatterns,
  } = useAiChatMessages({ sidebar, currentFormPatterns });

  watch(currentFormPatterns, async () => {
    if (!messages.value.length) {
      return;
    }

    if (changedPatternsFields.value.length === 0) {
      if (activeVersion.value >= 0 && lastLlmVersion.value >= 0 && lastLlmVersion.value !== activeVersion.value) {
        removeLastMessages(messages.value.length - lastLlmMessageIndex - 1);
        removeLastVersion();
      }

      return;
    }

    if (lastLlmVersion.value >= 0 && activeVersion.value >= 0 && lastLlmVersion.value !== activeVersion.value) {
      updateMessage(lastChangedPatternsFieldsMessageIndex, {
        text: tc('llm.chat.patternsEditedMessage', changedPatternsFields.value.length, {
          patterns: changedPatternsFields.value.map(field => t(`pattern.patternsFields.${field}`)).join(', ').toLowerCase(),
        }),
      });

      updateMessage(lastChangedPatternsFieldsMessageIndex + 1, {
        patterns: currentFormPatterns.value,
      });

      return;
    }

    lastLlmMessageIndex = messages.value.length - 1;
    lastChangedPatternsFieldsMessageIndex = lastLlmMessageIndex + 1;

    addPattern({ patterns: currentFormPatterns.value });
  });

  return {
    messages,
    addMessage,
    addPattern,
    resetMessages,

    versions,
    activeVersion,
    resetVersions,
    restoreVersion,

    currentFormPatterns,
    emptyCurrentFormPatterns,
  };
};

/**
 * Hook for the error line above the chat input (socket/validation failures).
 *
 * @returns {{ errorMessage: import('vue').Ref<string|null>, resetErrorMessage: function(): void }}
 */
export const useAiChatErrors = () => {
  const errorMessage = ref(null);

  /**
   * Clears the chat error so the alert above the input is hidden.
   */
  const resetErrorMessage = () => errorMessage.value = null;

  return {
    errorMessage,
    resetErrorMessage,
  };
};

/**
 * Hook for loading state while waiting for the LLM socket response; drives the thinking bubble message.
 *
 * @returns {Object} `thinking` ref, `thinkingMessage` computed, `enableThinking` / `disableThinking`.
 */
export const useAiChatThinking = () => {
  const thinking = ref(false);

  const thinkingMessage = computed(() => (thinking.value ? {
    _id: uid(),
    role: LLM_AI_CHAT_MESSAGE_ROLES.model,
    thinking: true,
  } : null));

  /**
   * Sets `thinking` to true so the UI shows the in-flight LLM state (thinking bubble).
   */
  const enableThinking = () => thinking.value = true;

  /**
   * Sets `thinking` to false after a response, error, or cancel.
   */
  const disableThinking = () => thinking.value = false;

  return {
    thinking,

    thinkingMessage,

    enableThinking,
    disableThinking,
  };
};

/**
 * Hook for joining `SOCKET_ROOMS.llmChat`, handling responses via `addPattern` / `setPatterns`, and `sendMessage`.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref|Object} [options.sidebar]
 *   Join payload from `config.socketRoomData`; applies patterns on success.
 * @param {function(Object): void} [options.addPattern] - LLM pattern payload on success.
 * @param {function(Object): void} [options.addMessage] - Add message to the chat.
 * @param {function(): void} [options.restorePrompt] - Restore prompt on socket error.
 * @param {function(): Promise<Array>} [options.fetchLlms] - Refetch LLM list (e.g. after Gone) to detect empty models.
 * @returns {Object} Thinking, errors, `sendMessage`, join/leave (unmount leaves room).
 */
export const useAiChatSocket = ({ sidebar, addPattern, addMessage, restorePrompt, fetchLlms } = {}) => {
  const { t } = useI18n();
  const socket = useSocket();

  let socketRoom = null;

  const { thinking, thinkingMessage, enableThinking, disableThinking } = useAiChatThinking();
  const { errorMessage, resetErrorMessage } = useAiChatErrors();

  /**
   * Socket error callback: only handles `SOCKET_ROOMS.llmChat`, maps codes to i18n, restores prompt.
   *
   * @param {Object} event
   * @param {string} [event.message] - Fallback text when `error.code` is missing.
   * @param {Object} [event.error] - Meta from the plugin (`room`, `code`, `payload` e.g. `retry_at`).
   */
  const errorHandler = async ({ message, error }) => {
    const { room, code, payload = {} } = error ?? {};

    if (room !== SOCKET_ROOMS.llmChat) {
      return;
    }

    const data = payload?.retry_at
      ? { retryAt: convertDateToStringWithFormatForToday(payload.retry_at, DATETIME_FORMATS.long, '', DATETIME_FORMATS.timeWithoutSeconds) }
      : {};

    let newErrorMessage = code ? t(`llm.chat.errors.${code}`, data) : message;

    if (code === LLM_AI_CHAT_ERROR_CODES.gone) {
      const llms = await fetchLlms?.() ?? [];

      if (!llms.length) {
        newErrorMessage = t('llm.chat.errors.noModels');
      }
    }

    errorMessage.value = newErrorMessage;

    disableThinking();
    restorePrompt();
  };

  /**
   * Room message handler: validation HTML, plain error, or success (`addPattern` + sidebar `setPatterns`).
   *
   * @param {Object} payload - Server payload; keys other than those below are spread as `patterns` for the form.
   * @param {*} [payload.error] - When set, shown as HTML error and prompt restored.
   * @param {*} [payload.code] - May accompany failures from the wire format.
   * @param {Array<string>} [payload.val_errors] - Renders a `<ul>` of validation messages.
   */
  const socketListener = async ({ error, code, val_errors: validationErrors = [], ...patterns }) => {
    await promisedTimeout(() => {}, 1000); // TODO: remove this

    resetErrorMessage();

    disableThinking();

    if (validationErrors?.length) {
      const validationMessage = `<ul>${validationErrors.map(validationError => `<li>${validationError}</li>`).join('')}</ul>`;

      addMessage({
        text: `${t('llm.chat.patternCannotBeCreatedReasons')}${validationMessage}`,
        role: LLM_AI_CHAT_MESSAGE_ROLES.model,
      });

      restorePrompt();

      return;
    }

    if (error) {
      addMessage({
        text: error,
        role: LLM_AI_CHAT_MESSAGE_ROLES.model,
      });

      restorePrompt();

      return;
    }

    addPattern({ patterns, role: LLM_AI_CHAT_MESSAGE_ROLES.model });

    unref(sidebar)?.config?.setPatterns?.(patterns);
  };

  /**
   * Joins `SOCKET_ROOMS.llmChat` merging `socketRoomData` from the sidebar with `config: llm`.
   * Registers `socketListener` on the resulting `SocketRoom`.
   *
   * @param {Object} llm - LLM config merged into the join payload as `config` (e.g. selected model row).
   */
  const joinSocketRoom = (llm) => {
    const joinData = {
      ...unref(sidebar)?.config?.socketRoomData,

      config: llm,
    };

    socketRoom = socket
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

    socketRoom = null;
  };

  /**
   * Turns on thinking state, ensures the LLM room is joined, then sends `data` on the socket.
   *
   * @param {Object} data - Outbound payload (e.g. `prompt` + `LLM_AI_CHAT_MESSAGE_TYPES.send`, or cancel).
   * @param {Object} llm - Merged as `config` on join when `socketRoom` is not open yet.
   * @returns {*} Whatever `SocketRoom.send` returns, or `undefined` if the room is missing.
   */
  const sendMessage = (data, llm) => {
    enableThinking();

    if (!socketRoom) {
      joinSocketRoom(llm);
    }

    return socketRoom?.send(data);
  };

  onBeforeUnmount(leaveSocketRoom);

  return {
    thinking,
    thinkingMessage,
    enableThinking,
    disableThinking,
    errorMessage,
    resetErrorMessage,
    sendMessage,
    joinSocketRoom,
    leaveSocketRoom,
  };
};

/**
 * Hook for the current textarea value and stash used when sending (restore on error/cancel).
 *
 * @returns {Object} `prompt`, `lastPrompt`, `updatePrompt`, `restorePrompt`, `resetPrompt`.
 */
export const useAiChatPrompt = () => {
  const prompt = ref('');
  const lastPrompt = ref('');

  /**
   * Sets the visible textarea ref value (e.g. typing, suggestion pick, programmatic clear).
   *
   * @param {string} [newPrompt=''] - Next prompt string.
   */
  const updatePrompt = (newPrompt = '') => prompt.value = newPrompt;

  /**
   * Puts `lastPrompt` back into `prompt` after send/error/cancel, then clears the stash.
   */
  const restorePrompt = () => {
    prompt.value = lastPrompt.value;
    lastPrompt.value = '';
  };

  /**
   * Clears both the field and the stashed copy (e.g. new question or full chat reset).
   */
  const resetPrompt = () => {
    prompt.value = '';
    lastPrompt.value = '';
  };

  return {
    prompt,
    lastPrompt,
    updatePrompt,
    restorePrompt,
    resetPrompt,
  };
};

/**
 * Hook for fetching LLM rows, keeping `llm` in sync with the default model, and exposing `fetchLlms`.
 *
 * @param {Object} [options]
 * @returns {Object} `llm`, `llms`, `llmsPending`, `fetchLlms`, `resetLlm`.
 */
export const useAiChatLlms = ({ initialLlms = [] } = {}) => {
  const llms = ref([...unref(initialLlms)]);

  const { fetchLlmsListWithoutStore } = useLlm();

  const { pending: llmsPending, handler: fetchLlms } = usePendingHandler(async () => {
    const { data } = await fetchLlmsListWithoutStore({ params: { enabled: true, limit: MAX_LIMIT } });

    llms.value = data;

    return data;
  });

  return {
    llms,
    llmsPending,
    fetchLlms,
  };
};

/**
 * Hook for fetching LLM rows, keeping `llm` in sync with the default model, and exposing `fetchLlms`.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref|Object} [options.sidebar]
 * @returns {Object} `llm`, `llms`, `llmsPending`, `fetchLlms`, `resetLlm`.
 */
export const useAiChatLlmModel = ({ sidebar } = {}) => {
  const llm = ref(null);

  const initialLlms = computed(() => unref(sidebar)?.config?.llms ?? []);

  const { llms, llmsPending, fetchLlms } = useAiChatLlms({ initialLlms });

  const defaultLlm = computed(() => llms.value.find(llmItem => llmItem.default) ?? null);

  /**
   * Resets `llm` to the row marked `default` in the current `llms` list (or `null` if none).
   */
  const resetLlm = () => llm.value = defaultLlm.value;

  watch(llms, () => {
    if (defaultLlm.value) {
      llm.value = defaultLlm.value;
    }
  }, { immediate: true });

  return {
    llm,
    llms,
    llmsPending,

    fetchLlms,
    resetLlm,
  };
};

/**
 * Hook for keeping the chat body scrolled to the bottom when messages or thinking state change.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref} [options.messages] - Chat messages ref (deep watch).
 * @param {import('vue').Ref} [options.thinking] - Thinking flag ref.
 * @returns {{ bodyElement: import('vue').Ref, scrollChatBodyToBottom: function(): void }}
 */
export const useAiChatScroll = ({ messages, thinking } = {}) => {
  const bodyElement = ref(null);

  /**
   * Schedules a scroll of `bodyElement` to the bottom after DOM updates (`nextTick` + `requestAnimationFrame`).
   *
   * No-op when `bodyElement` is not mounted.
   */
  const scrollChatBodyToBottom = () => nextTick(() => requestAnimationFrame(() => {
    const el = bodyElement.value;

    if (!el) {
      return;
    }

    el.scrollTop = el.scrollHeight;
  }));

  watch(messages, scrollChatBodyToBottom, { deep: true, flush: 'post' });
  watch(thinking, scrollChatBodyToBottom, { flush: 'post' });

  return {
    bodyElement,
    scrollChatBodyToBottom,
  };
};

/**
 * Hook for the prompt textarea template ref and applying a suggestion string with focus.
 *
 * @param {Object} options
 * @param {function(string): void} options.updatePrompt - Sets the prompt text (e.g. from `useAiChatPrompt`).
 * @returns {{ textareaElement: import('vue').Ref, applySuggestion: function(string): void }}
 */
export const useAiChatSuggestions = ({ updatePrompt }) => {
  const textareaElement = ref(null);

  /**
   * Fills the prompt from a starter suggestion and focuses the bound textarea ref.
   *
   * @param {string} newPrompt - Full suggestion text (e.g. from `llm.chat.suggestionPrompts`).
   */
  const applySuggestion = (newPrompt) => {
    updatePrompt(newPrompt);
    textareaElement.value?.focus?.();
  };

  return {
    textareaElement,

    applySuggestion,
  };
};

/**
 * Hook for composing AI chat sidebar state: prompt, model, messages, socket, scroll, JSON auto-ask, restart.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref|Object} [options.sidebar]
 *   Sidebar `config`: patterns, jsonString, socket, setPatterns, setPending.
 * @returns {Object} `ai-chat.vue` API: `ask`, `stop`, `restart`, refs.
 */
export const useAiChat = ({ sidebar } = {}) => {
  const { prompt, updatePrompt, resetPrompt, restorePrompt } = useAiChatPrompt();
  const { llm, llms, llmsPending, resetLlm, fetchLlms } = useAiChatLlmModel({ sidebar });
  const { textareaElement, applySuggestion } = useAiChatSuggestions({ updatePrompt });

  const {
    currentFormPatterns,
    emptyCurrentFormPatterns,

    messages,
    addMessage,
    addPattern,
    resetMessages,
    resetVersions,

    versions,
    activeVersion,
    restoreVersion,
  } = useAiChatPattern({ sidebar });

  const {
    thinking,
    thinkingMessage,
    disableThinking,

    errorMessage,
    resetErrorMessage,

    sendMessage,

    leaveSocketRoom,
  } = useAiChatSocket({
    sidebar,
    addPattern,
    addMessage,
    restorePrompt,
    fetchLlms,
  });

  const { bodyElement } = useAiChatScroll({ messages, thinking });

  const emptyChat = computed(() => messages.value.length === 0);

  /**
   * Sends the user prompt over the LLM socket and appends a user-role message to the transcript.
   *
   * @param {Object} params
   * @param {string} params.prompt - Text sent as `data.prompt` with `LLM_AI_CHAT_MESSAGE_TYPES.send`.
   * @param {boolean} [params.withoutPatterns=false] - Reserved: when backend merges form patterns into `data`
   *   (currently gated by a TODO condition).
   */
  const ask = ({ prompt: newPrompt, withoutPatterns = false }) => {
    resetErrorMessage();
    resetPrompt();

    const data = {
      prompt: newPrompt,
      type: LLM_AI_CHAT_MESSAGE_TYPES.send,
    };

    if (false && !withoutPatterns && !emptyCurrentFormPatterns.value) { // TODO: remove false
      Object.entries(currentFormPatterns.value).forEach(([field, pattern]) => {
        data[field] = pattern;
      });
    }

    sendMessage(data, llm.value?._id);
    addMessage({ text: newPrompt, role: LLM_AI_CHAT_MESSAGE_ROLES.user });
  };

  /**
   * Emits cancel on the room and restores the stashed prompt after aborting an in-flight request.
   */
  const stop = () => {
    sendMessage({ type: LLM_AI_CHAT_MESSAGE_TYPES.cancel }, llm.value?._id);

    restorePrompt();
  };

  /**
   * Local-only reset: prompt, selected model, errors, messages, versions, thinking (does not leave socket).
   */
  const resetChat = () => {
    resetPrompt();
    resetLlm();
    resetErrorMessage();
    resetMessages();
    resetVersions();
    disableThinking();
  };

  /**
   * Leaves the LLM socket room then runs `resetChat` (full session reset from the UI).
   */
  const restart = () => {
    leaveSocketRoom();
    resetChat();
  };

  useAiChatJsonString({
    sidebar,
    ask,
  });

  watch(thinking, newThinking => (
    unref(sidebar)?.config?.setPending?.(newThinking, emptyCurrentFormPatterns.value, stop)
  ));

  return {
    bodyElement,
    textareaElement,

    prompt,

    llm,
    llms,
    llmsPending,

    messages,
    emptyChat,

    versions,
    activeVersion,
    restoreVersion,

    thinkingMessage,
    errorMessage,

    ask,
    stop,
    applySuggestion,
    restart,
  };
};
