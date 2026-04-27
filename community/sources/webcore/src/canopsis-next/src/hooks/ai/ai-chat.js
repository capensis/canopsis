import { pick, isEmpty, isNil } from 'lodash';
import {
  computed,
  ref,
  unref,
  set,
  watch,
  inject,
  nextTick,
  onBeforeUnmount,
} from 'vue';

import { SOCKET_ROOMS } from '@/config';
import {
  MAX_LIMIT,
  MODALS,
  LLM_AI_CHAT_MESSAGE_ROLES,
  LLM_AI_CHAT_MESSAGE_TYPES,
  LLM_AI_CHAT_ERROR_CODES,
  PATTERNS_FIELDS,
  DATETIME_FORMATS,
  LLM_SOCKET_CONTEXTS,
} from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import { uid } from '@/helpers/uid';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';
import { getChangedPatternsFields } from '@/helpers/entities/pattern/form';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSocket } from '@/hooks/socket';
import { usePendingHandler } from '@/hooks/query/pending';
import { useLlm } from '@/hooks/store/modules/llm';

/**
 * Base fields merged into every new chat row (`_id`, `timestamp`).
 *
 * @returns {{ _id: string, timestamp: number }}
 */
const getEmptyMessage = () => ({
  _id: uid(),
  timestamp: Date.now(),
});

/**
 * Watches a pattern JSON string ref and calls `ask` with the localized fix-pattern prompt on each non-empty
 * value.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<string>|import('vue').ComputedRef<string>|string} [options.jsonString] -
 *   Observed JSON text (e.g. from the pattern editor).
 * @param {function({ prompt: string, withoutPatterns: boolean })} [options.ask] - Invoked when `jsonString`
 *   updates to a truthy string.
 * @returns {{ jsonString: import('vue').Ref<string>|import('vue').ComputedRef<string>|string }} Echo of
 *   `options.jsonString` (for chaining); the composable only registers the watcher.
 */
export const useAiChatJsonString = ({ jsonString, ask } = {}) => {
  const { t } = useI18n();

  watch(jsonString, newJsonString => newJsonString && ask({
    prompt: t('llm.chat.fixPatternPrompt', { jsonString: newJsonString }),
    withoutPatterns: true,
  }));
};

/**
 * Stack of pattern snapshots produced in the chat and index of the version last attributed to the LLM.
 *
 * @returns {{
 *   versions: import('vue').Ref<Array<Object>>,
 *   lastLlmVersion: import('vue').Ref<number|null>,
 *   activeVersion: import('vue').ComputedRef<number>,
 *   lastLlmPatterns: import('vue').ComputedRef<Object|null>,
 *   addVersion: function(Object, boolean),
 *   removeLastVersion: function(),
 *   resetVersions: function(),
 * }}
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
   * Clears `versions` and sets `lastLlmVersion` back to `null`.
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
 * Chat transcript, pattern-card rows, and version history; diffs the live form patterns against the last LLM
 * snapshot for status text and version bumps.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<Object>|import('vue').ComputedRef<Object>} [options.currentFormPatterns] -
 *   Normalized patterns from the host, compared to the last LLM version.
 * @param {function(Object)} [options.updateFormPatterns] - Writes patterns into the host when restoring a
 *   version from `restoreVersion`.
 * @returns {{
 *   messages: import('vue').Ref<Array<Object>>,
 *   addMessage: function(Object): number,
 *   addPattern: function(Object),
 *   updateMessage: function(number, Object),
 *   removeLastMessages: function(number=): Array,
 *   resetMessages: function,
 *   versions: import('vue').Ref<Array<Object>>,
 *   activeVersion: import('vue').ComputedRef<number>,
 *   lastLlmVersion: import('vue').Ref<number|null>,
 *   addVersion: function(Object, boolean=),
 *   removeLastVersion: function,
 *   restoreVersion: function(number),
 *   resetVersions: function,
 *   changedPatternsFields: import('vue').ComputedRef<Array<string>>,
 *   emptyCurrentFormPatterns: import('vue').ComputedRef<boolean>,
 * }}
 */
export const useAiChatMessages = ({ currentFormPatterns, updateFormPatterns } = {}) => {
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
   * Appends a chat row; merges `getEmptyMessage()` (`_id`, `timestamp`) with `payload`.
   *
   * @param {Object} payload - Row fields; additional keys (e.g. `patterns`, `error`, `val_errors`, `local`,
   *   `thinking`) are kept on the stored object.
   * @param {string} [payload.prompt] - Visible body for a prompt bubble.
   * @param {string} [payload.role] - One of `LLM_AI_CHAT_MESSAGE_ROLES` when set.
   * @returns {number} New `messages.length` after `Array.prototype.push`.
   */
  const addMessage = ({ prompt, role, ...rest }) => messages.value.push({
    prompt,
    role,

    ...getEmptyMessage(),
    ...rest,
  });

  /**
   * Pushes a pattern version, optionally a status text line, then a pattern card message.
   *
   * @param {Object} options
   * @param {Object} options.patterns - Pattern payload for this version.
   * @param {string} [options.role] - `LLM_AI_CHAT_MESSAGE_ROLES.model` for LLM turns; omit when syncing user edits.
   * @param {number} [options.fromVersion] - Set when restoring a prior version (skips extra status text).
   * @param {boolean} [options.local] - Set when adding a local message.
   */
  const addPattern = ({ patterns, role, fromVersion, local }) => {
    let messagePrompt = '';

    if (!role && changedPatternsFields.value.length > 0) {
      messagePrompt = tc('llm.chat.patternsEditedMessage', changedPatternsFields.value.length, {
        patterns: changedPatternsFields.value.map(field => t(`pattern.patternsFields.${field}`)).join(', ').toLowerCase(),
      });

      addVersion(patterns);
    } else if (role === LLM_AI_CHAT_MESSAGE_ROLES.model) {
      addVersion(patterns, true);
    }

    addMessage({
      patterns,
      local,
      role,

      _id: uid(),
      prompt: messagePrompt,
      version: activeVersion.value,
      from_version: fromVersion,
      emptyPatterns: versions.value.length === 1 && emptyCurrentFormPatterns.value,
    });
  };

  /**
   * Replaces one message at `index` (Vue 2 reactivity-safe).
   *
   * @param {number} index - Index in `messages`.
   * @param {Object} newMessage - Full replacement message object.
   */
  const updateMessage = (index, newMessage) => (
    set(messages.value, index, { ...(messages.value[index] ?? getEmptyMessage()), ...newMessage })
  );

  /**
   * Drops the last `count` messages from the tail of the list.
   *
   * @param {number} [count=1] - How many entries to remove from the end.
   * @returns {Array} Elements removed by `splice`.
   */
  const removeLastMessages = (count = 1) => messages.value.splice(messages.value.length - count);

  /**
   * Replays patterns from `versions[version]` into the transcript and applies them via `updateFormPatterns`.
   *
   * @param {number} version - Index in `versions`.
   */
  const restoreVersion = (version) => {
    const newPatterns = versions.value[version];

    addPattern({ patterns: newPatterns, role: LLM_AI_CHAT_MESSAGE_ROLES.model, fromVersion: version });

    updateFormPatterns(newPatterns);
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
 * Keeps filter patterns, chat messages, and the version stack aligned when the user edits the form after LLM
 * replies (watches `currentFormPatterns`).
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<Object>|import('vue').ComputedRef<Object>} [options.patterns] - Raw pattern map
 *   from the host (e.g. modal form).
 * @param {import('vue').Ref<Object|null>} [options.patternItem] - Optional selected item (`value` is the key
 *   into `patterns`); narrows `currentFormPatterns` when set.
 * @param {function(Object)} [options.updateFormPatterns] - Callback used when restoring versions (same shape
 *   as in `useAiChatMessages`).
 * @returns {{
 *   messages: import('vue').Ref<Array<Object>>,
 *   addMessage: function(Object): number,
 *   addPattern: function(Object),
 *   resetMessages: function,
 *   versions: import('vue').Ref<Array<Object>>,
 *   activeVersion: import('vue').ComputedRef<number>,
 *   resetVersions: function,
 *   restoreVersion: function(number),
 *   currentFormPatterns: import('vue').ComputedRef<Object>,
 *   emptyCurrentFormPatterns: import('vue').ComputedRef<boolean>,
 * }}
 */
export const useAiChatPattern = ({ patterns, patternItem, updateFormPatterns } = {}) => {
  const { t, tc } = useI18n();

  const currentFormPatterns = computed(() => {
    const unwrappedPatternItem = unref(patternItem);

    return formFilterToPatterns(
      unwrappedPatternItem
        ? { [unwrappedPatternItem.value]: unref(patterns)?.[unwrappedPatternItem.value] }
        : unref(patterns) ?? {},
      Object.values(PATTERNS_FIELDS),
    );
  });

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
  } = useAiChatMessages({ currentFormPatterns, updateFormPatterns });

  const lastLlmMessageIndex = computed(() => messages.value.findLastIndex(message => (
    message.role === LLM_AI_CHAT_MESSAGE_ROLES.model && !message.local
  )));

  const lastChangedPatternsFieldsMessageIndex = computed(() => messages.value.findLastIndex(message => (
    message.prompt && message.local
  )));

  watch(currentFormPatterns, async () => {
    if (!versions.value.length) {
      return;
    }

    if (changedPatternsFields.value.length === 0) {
      if (lastLlmVersion.value < activeVersion.value) {
        removeLastMessages(1);
        removeLastVersion();
      }

      return;
    }

    if (lastLlmVersion.value < activeVersion.value) {
      updateMessage(lastChangedPatternsFieldsMessageIndex.value, {
        prompt: tc('llm.chat.patternsEditedMessage', changedPatternsFields.value.length, {
          patterns: changedPatternsFields.value.map(field => t(`pattern.patternsFields.${field}`)).join(', ').toLowerCase(),
        }),
      });

      const newMessage = emptyCurrentFormPatterns.value ? {
        prompt: t('llm.chat.emptyPatternsMessage'),
        patterns: undefined,
      } : {
        patterns: currentFormPatterns.value,
      };

      if (emptyCurrentFormPatterns.value) {
        removeLastVersion();
      }

      updateMessage(lastChangedPatternsFieldsMessageIndex.value, newMessage);

      return;
    }

    if (lastChangedPatternsFieldsMessageIndex.value > lastLlmMessageIndex.value) {
      removeLastMessages(1);
    }

    if (emptyCurrentFormPatterns.value) {
      addMessage({
        prompt: t('llm.chat.emptyPatternsMessage'),
        local: true,
      });
    } else {
      addPattern({ patterns: currentFormPatterns.value, local: true });
    }
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
 * Inline error string above the chat input (socket, validation, or pattern-item errors).
 *
 * @returns {{
 *   errorMessage: import('vue').Ref<string|null>,
 *   resetErrorMessage: function,
 * }}
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
 * In-flight flag and a synthetic “thinking” placeholder message for the transcript while the socket request
 * runs.
 *
 * @returns {{
 *   thinking: import('vue').Ref<boolean>,
 *   thinkingMessage: import('vue').ComputedRef<Object|null>,
 *   enableThinking: function,
 *   disableThinking: function,
 * }}
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
 * Joins `SOCKET_ROOMS.llmChat`, wires listeners, and exposes `sendMessage` with shared thinking/error state.
 * Unregisters and leaves the room in `onBeforeUnmount`.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<string>|string} [options.context] - LLM socket context (e.g. `LLM_SOCKET_CONTEXTS.*`).
 * @param {import('vue').Ref<string>|string} [options.ruleId] - Optional rule id sent as `rule` on join.
 * @param {function(Object)} [options.addPattern] - Appends a pattern-card message when the payload is a
 *   successful apply.
 * @param {function(Object)} [options.updateFormPatterns] - Applies parsed pattern fields to the host on success.
 * @param {function(Object)} [options.addMessage] - Appends chat rows (wire errors, `val_errors`, etc.).
 * @param {function} [options.restorePrompt] - Restores the prompt after recoverable socket or validation errors.
 * @param {function(): Promise<Array>|Array} [options.fetchLlms] - Optional refetch of enabled LLMs (e.g. after
 *   `LLM_AI_CHAT_ERROR_CODES.gone`) to refine the error message.
 * @returns {{
 *   thinking: import('vue').Ref<boolean>,
 *   thinkingMessage: import('vue').ComputedRef<Object|null>,
 *   enableThinking: function,
 *   disableThinking: function,
 *   errorMessage: import('vue').Ref<string|null>,
 *   resetErrorMessage: function,
 *   sendMessage: function(Object, Object|string): *,
 *   joinSocketRoom: function(Object|string),
 *   leaveSocketRoom: function,
 * }}
 */
export const useAiChatSocket = ({
  context,
  ruleId,
  addPattern,
  updateFormPatterns,
  addMessage,
  restorePrompt,
  fetchLlms,
} = {}) => {
  const { t } = useI18n();
  const socket = useSocket();

  const aiChat = inject('$aiChat', {});

  const needRestart = ref(false);

  let socketRoom = null;

  const { thinking, thinkingMessage, enableThinking, disableThinking } = useAiChatThinking();
  const { errorMessage, resetErrorMessage } = useAiChatErrors();

  /**
   * Writes the `needRestart` ref so callers (or nested composables) can surface a "restart chat" state
   * without reaching into the ref directly.
   *
   * @param {boolean} newNeedRestart - Next value for `needRestart` (e.g. `true` when pattern items go stale
   *   mid-session).
   */
  const updateNeedRestart = newNeedRestart => needRestart.value = newNeedRestart;

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
      ? { retryAt: convertDateToStringWithFormatForToday(payload.retry_at, DATETIME_FORMATS.long, '', DATETIME_FORMATS.timePicker) }
      : {};

    let newErrorMessage = code ? t(`llm.chat.errors.${code}`, data) : message;

    if (code === LLM_AI_CHAT_ERROR_CODES.gone) {
      const llms = await fetchLlms?.() ?? [];

      if (!llms.length) {
        newErrorMessage = t('llm.chat.errors.noModels');
      }

      updateNeedRestart(true);
    }

    errorMessage.value = newErrorMessage;

    disableThinking();
    restorePrompt();
  };

  /**
   * Room message handler: validation list, wire error, or success (`addPattern` + `updateFormPatterns`).
   *
   * @param {Object} payload - Server payload; keys other than those below are spread as `patterns` for the form.
   * @param {*} [payload.error] - When set, shown as HTML error and prompt restored.
   * @param {*} [payload.code] - May accompany failures from the wire format.
   * @param {Array<string>} [payload.val_errors] - Renders a `<ul>` of validation messages.
   */
  const socketListener = async ({ error, code, val_errors: validationErrors = [], chat, ...rest }) => {
    const patterns = pick(rest, Object.values(PATTERNS_FIELDS));

    if (chat) {
      aiChat.registerChatId?.(chat);
    }

    resetErrorMessage();

    disableThinking();

    if (error) {
      addMessage({
        error,
        role: LLM_AI_CHAT_MESSAGE_ROLES.model,
      });

      restorePrompt();

      return;
    }

    if (validationErrors?.length) {
      addMessage({
        val_errors: validationErrors,
        role: LLM_AI_CHAT_MESSAGE_ROLES.model,
      });

      restorePrompt();

      return;
    }

    addPattern({ patterns, role: LLM_AI_CHAT_MESSAGE_ROLES.model });
    updateFormPatterns(patterns);
  };

  /**
   * Generic HTTP/server failure on the socket: clears thinking, restores the prompt, sets a generic error.
   */
  const serverErrorHandler = () => {
    resetErrorMessage();
    disableThinking();
    restorePrompt();

    errorMessage.value = t('llm.chat.errors.serverError');
  };

  /**
   * Joins `SOCKET_ROOMS.llmChat` with `{ context, rule, config: llm }` and registers `socketListener`.
   *
   * @param {Object|string} llm - Model config for the join payload (`config`); may be a full document or id
   *   depending on the caller.
   */
  const joinSocketRoom = (llm) => {
    const joinData = {
      context: unref(context),
      rule: unref(ruleId),
      config: llm,
    };

    socketRoom = socket
      .on(Socket.EVENTS_TYPES.error, errorHandler)
      .on(Socket.EVENTS_TYPES.serverError, serverErrorHandler)
      .join(SOCKET_ROOMS.llmChat, joinData, true)
      .addListener(socketListener);
  };

  /**
   * Leaves the LLM chat socket room and clears the cached `SocketRoom` reference (runs on unmount).
   */
  const leaveSocketRoom = () => {
    socket
      .off(Socket.EVENTS_TYPES.error, errorHandler)
      .off(Socket.EVENTS_TYPES.serverError, serverErrorHandler)
      .leave(SOCKET_ROOMS.llmChat)
      .removeListener(socketListener);

    socketRoom = null;
  };

  /**
   * Enables thinking, joins the LLM room on first use, then sends `data` on the active room.
   *
   * @param {Object} data - Outbound payload (e.g. `prompt` + `LLM_AI_CHAT_MESSAGE_TYPES.send`, or cancel).
   * @param {Object|string} llm - Model config for `joinSocketRoom` when the room is not open yet (callers may
   *   pass a full document or an id, depending on the socket layer).
   * @returns {*} Return value of `SocketRoom.send`, or `undefined` if the room could not be created.
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
    needRestart,
    updateNeedRestart,
  };
};

/**
 * Prompt field, `lastPrompt` stash for send/cancel/error recovery, and small setters.
 *
 * @returns {{
 *   prompt: import('vue').Ref<string>,
 *   lastPrompt: import('vue').Ref<string>,
 *   updatePrompt: function(string=),
 *   restorePrompt: function,
 *   resetPrompt: function,
 * }}
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
 * Local `llms` list and a `usePendingHandler`-wrapped fetch of enabled models from the API.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<Array>|Array} [options.initialLlms=[]] - Seed list before the first fetch.
 * @returns {{
 *   llms: import('vue').Ref<Array<Object>>,
 *   llmsPending: import('vue').Ref<boolean>,
 *   fetchLlms: function(): Promise<Array<Object>>,
 * }}
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
 * Selected model (`llm`), list + fetch from `useAiChatLlms`, and sync to the row with `default: true`.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<Array>|Array} [options.initialLlms=[]] - Initial rows passed to `useAiChatLlms`.
 * @returns {{
 *   llm: import('vue').Ref<Object|null>,
 *   llms: import('vue').Ref<Array<Object>>,
 *   llmsPending: import('vue').Ref<boolean>,
 *   fetchLlms: function(): Promise<Array<Object>>,
 *   resetLlm: function,
 * }}
 */
export const useAiChatLlmModel = ({ initialLlms = [] } = {}) => {
  const llm = ref(null);

  const { llms, llmsPending, fetchLlms } = useAiChatLlms({ initialLlms });

  const defaultLlm = computed(() => llms.value.find(llmItem => llmItem.default) ?? llms.value[0] ?? null);

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
 * Binds the chat body element and keeps it scrolled to the bottom when `messages` (deep) or `thinking` changes.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<Array>} [options.messages] - Watched with `{ deep: true, flush: 'post' }`.
 * @param {import('vue').Ref<boolean>} [options.thinking] - Watched with `{ flush: 'post' }` for the thinking row.
 * @returns {{
 *   bodyElement: import('vue').Ref<HTMLElement|undefined>,
 *   scrollChatBodyToBottom: function,
 * }}
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
 * Template ref for the prompt textarea plus `applySuggestion` for starter prompts (sets text and focuses).
 *
 * @param {Object} options
 * @param {function(string=)} options.updatePrompt - Setter from `useAiChatPrompt`.
 * @returns {{
 *   textareaElement: import('vue').Ref<HTMLTextAreaElement|undefined>,
 *   applySuggestion: function(string),
 * }}
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
 * Pattern-item selector for scenario context: selection ref, visibility flag, label, and reset.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<string>} [options.context] - Compared with `LLM_SOCKET_CONTEXTS.scenario` for
 *   `hasPatternItem`; `patternsItemsLabel` uses `unref(context)`.
 * @returns {{
 *   patternItem: import('vue').Ref<Object|null>,
 *   hasPatternItem: import('vue').ComputedRef<boolean>,
 *   patternsItemsLabel: import('vue').ComputedRef<string>,
 *   resetPatternItem: function,
 * }}
 */
export const useAiChatPatternItem = ({ context } = {}) => {
  const { t, te } = useI18n();

  const patternItem = ref(null);

  const hasPatternItem = computed(() => context.value === LLM_SOCKET_CONTEXTS.scenario);

  const patternsItemsLabel = computed(() => {
    const labelKey = `llm.chat.patternsItemsLabel.${unref(context)}`;

    if (te(labelKey)) {
      return t(labelKey);
    }

    return te(labelKey) ? t(labelKey) : t('common.patternsItems');
  });

  const resetPatternItem = () => patternItem.value = null;

  return {
    patternItem,
    hasPatternItem,
    patternsItemsLabel,

    resetPatternItem,
  };
};

/**
 * Full AI chat composable: prompt, models, pattern sync, socket, scroll, JSON-triggered ask, and restart flow.
 * Emits `update:patterns` when patterns are applied from the socket or restore, and `update:pending` while the
 * LLM request is in flight.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<Object>} [options.patterns] - Live pattern map from the host.
 * @param {import('vue').Ref<string>|string} [options.context] - Socket join `context`.
 * @param {import('vue').Ref<Array>|Array} [options.patternsItems] - Selectable pattern groups; watched to detect
 *   stale `patternItem` and set `needRestart` / error when needed.
 * @param {import('vue').Ref<string>|string} [options.ruleId] - Optional `rule` on socket join.
 * @param {import('vue').Ref<string>|string} [options.jsonString] - When it updates, triggers `ask` with the
 *   fix-pattern prompt via `useAiChatJsonString`.
 * @param {import('vue').Ref<Array>|Array} [options.llms] - Initial LLM rows (`initialLlms` for `useAiChatLlmModel`).
 * @param {function(string, ...*)} emit - Vue `setup` emit; must handle `update:patterns` and `update:pending`.
 * @returns {{
 *   bodyElement: import('vue').Ref<HTMLElement|undefined>,
 *   textareaElement: import('vue').Ref<HTMLTextAreaElement|undefined>,
 *   needRestart: import('vue').Ref<boolean>,
 *   prompt: import('vue').Ref<string>,
 *   llm: import('vue').Ref<Object|null>,
 *   llms: import('vue').Ref<Array<Object>>,
 *   llmsPending: import('vue').Ref<boolean>,
 *   messages: import('vue').Ref<Array<Object>>,
 *   emptyChat: import('vue').ComputedRef<boolean>,
 *   versions: import('vue').Ref<Array<Object>>,
 *   activeVersion: import('vue').ComputedRef<number>,
 *   restoreVersion: function(number),
 *   thinkingMessage: import('vue').ComputedRef<Object|null>,
 *   errorMessage: import('vue').Ref<string|null>,
 *   patternItem: import('vue').Ref<Object|null>,
 *   patternsItems: import('vue').Ref<Array>|Array|undefined,
 *   hasPatternItem: import('vue').ComputedRef<boolean>,
 *   patternsItemsLabel: import('vue').ComputedRef<string>,
 *   ask: function({ prompt: string, withoutPatterns?: boolean }),
 *   stop: function,
 *   applySuggestion: function(string),
 *   restart: function,
 * }}
 */
export const useAiChat = ({
  patterns,
  context,
  patternsItems,
  ruleId,
  jsonString,
  llms: initialLlms,
} = {}, emit) => {
  const { t } = useI18n();
  const modals = useModals();
  const { prompt, updatePrompt, resetPrompt, restorePrompt } = useAiChatPrompt();
  const { llm, llms, llmsPending, resetLlm, fetchLlms } = useAiChatLlmModel({ initialLlms });
  const { textareaElement, applySuggestion } = useAiChatSuggestions({ updatePrompt });

  const {
    patternItem,
    hasPatternItem,
    patternsItemsLabel,
    resetPatternItem,
  } = useAiChatPatternItem({ context });

  const updateFormPatterns = newPatterns => emit('update:patterns', newPatterns, patternItem.value?.value);

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
  } = useAiChatPattern({ patterns, patternItem, updateFormPatterns });

  const {
    thinking,
    thinkingMessage,
    disableThinking,

    errorMessage,
    resetErrorMessage,

    sendMessage,

    leaveSocketRoom,

    needRestart,
    updateNeedRestart,
  } = useAiChatSocket({
    context,
    ruleId,
    addPattern,
    updateFormPatterns,
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

    if (!withoutPatterns && !emptyCurrentFormPatterns.value) {
      Object.entries(currentFormPatterns.value).forEach(([field, pattern]) => {
        data[field] = pattern;
      });
    }

    sendMessage(data, llm.value?._id);
    addMessage({ prompt: newPrompt, role: LLM_AI_CHAT_MESSAGE_ROLES.user });
  };

  /**
   * Emits cancel on the room and restores the stashed prompt after aborting an in-flight request.
   */
  const stop = () => {
    sendMessage({ type: LLM_AI_CHAT_MESSAGE_TYPES.cancel }, llm.value?._id);
    disableThinking();
    restorePrompt();
  };

  /**
   * Local-only reset: prompt, selected model, errors, messages, versions, thinking (does not leave socket).
   */
  const resetChat = () => {
    updateNeedRestart(false);
    resetPatternItem();
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
  const restart = () => modals.show({
    name: MODALS.confirmation,
    config: {
      text: t('llm.chat.restartConfirmation'),
      action: () => {
        leaveSocketRoom();
        resetChat();
      },
    },
  });

  useAiChatJsonString({
    jsonString,
    ask,
  });

  watch(thinking, newThinking => (
    emit('update:pending', newThinking, emptyCurrentFormPatterns.value)
  ));

  watch(patternsItems, (items) => {
    if (patternItem.value?.value && !items.some(item => item.value === patternItem.value.value)) {
      if (emptyChat.value) {
        resetPatternItem();

        return;
      }

      errorMessage.value = t(`llm.chat.patternsItemsError.${unref(context)}`);
      updateNeedRestart(true);
    }
  });

  return {
    bodyElement,
    textareaElement,

    needRestart,

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

    patternItem,
    patternsItems,
    hasPatternItem,
    patternsItemsLabel,

    ask,
    stop,
    applySuggestion,
    restart,
  };
};
