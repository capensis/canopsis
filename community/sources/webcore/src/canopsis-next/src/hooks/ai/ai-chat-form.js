import { throttle, isArray } from 'lodash';
import {
  computed,
  watch,
  ref,
  unref,
  set,
  provide,
  inject,
  nextTick,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { LLM_AI_CHAT_WIDTH, LLM_SOCKET_CONTEXTS, LLM_AI_CHAT_TOURS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';
import { aiChatFormToPatterns } from '@/helpers/entities/llm/chat/form';
import { patternToForm } from '@/helpers/entities/pattern/form';

import { useAuth } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useValidator } from '@/hooks/validator/validator';
import { useUser } from '@/hooks/store/modules/user';
import { useLlm } from '@/hooks/store/modules/llm';

import { useAiChatLlmModel } from './ai-chat';

/** Throttle window for syncing `formPatterns` from the live form into sidebar props. */
const THROTTLE_WAIT_MS = 1000;

/**
 * Builds `patternsItems` options for the AI sidebar from the modal form (scenario: one row per action).
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<Array>|Array} [options.form=[]] - Host form rows (e.g. scenario `actions`); only
 *   `LLM_SOCKET_CONTEXTS.scenario` is mapped today.
 * @param {string} [options.context] - LLM socket context key; must match a key in the internal map (pass a
 *   resolved string; a bare `Ref`/`ComputedRef` is not unwrapped here).
 * @returns {{ patternsItems: import('vue').ComputedRef<Array<{ text: string, value: string }>> }}
 */
export const useAiChatPatternsItems = ({ form = [], context } = {}) => {
  const { t } = useI18n();

  const contextToPatternsItems = {
    [LLM_SOCKET_CONTEXTS.scenario]: () => unref(form)?.map?.((formItem, index) => ({
      text: `${index + 1} - ${t(`scenario.actions.${formItem.type}`)}`,
      value: formItem.key,
    })),
  };

  const patternsItems = computed(() => contextToPatternsItems[unref(context)]?.(form) ?? []);

  return {
    patternsItems,
  };
};

/**
 * Persists “minimized” LLM chrome in user tours and mirrors padding on the modal content wrapper.
 *
 * @param {Object} [options]
 * @param {{ id: string, dialogProps: Object }} [options.modal] - Active modal descriptor passed to
 *   `modals.updateDialogProps`.
 * @returns {{
 *   initialMinimized: import('vue').ComputedRef<boolean|undefined>,
 *   updateMinimized: function(boolean),
 * }}
 */
const useAiChatMinimized = ({ modal } = {}) => {
  const modals = useModals();
  const { currentUser } = useAuth();
  const { updateCurrentUserTours } = useUser();

  const initialMinimized = computed(() => currentUser.value?.ui_tours?.[LLM_AI_CHAT_TOURS.minimized]);

  /**
   * Saves minimized state to the user profile and updates `contentWrapperStyle` on the host modal.
   *
   * @param {boolean} newMinimized - Next minimized flag from the sidebar.
   */
  const updateMinimized = (newMinimized) => {
    updateCurrentUserTours({ data: { [LLM_AI_CHAT_TOURS.minimized]: newMinimized } });

    const unwrappedModal = unref(modal) ?? {};

    if (unwrappedModal.id) {
      modals.updateDialogProps({
        id: unwrappedModal.id,
        dialogProps: {
          ...unwrappedModal.dialogProps,

          contentWrapperStyle: {
            paddingRight: newMinimized ? 0 : `${LLM_AI_CHAT_WIDTH}px`,
          },
        },
      });
    }
  };

  return {
    initialMinimized,

    updateMinimized,
  };
};

/**
 * Overlay “pending” state for pattern apply from the LLM chat, plus cancel callback wiring.
 *
 * @returns {{
 *   pending: import('vue').Ref<boolean>,
 *   pendingTexts: import('vue').ComputedRef<{ inProgress: string, cancel: string }>,
 *   updatePending: function(boolean, boolean|null=, function=),
 *   cancel: function,
 * }}
 */
export const useAiChatPending = () => {
  const { t } = useI18n();

  const pending = ref(false);
  const creation = ref(false);

  const pendingTexts = computed(() => (creation.value ? {
    inProgress: t('pattern.patternsCreationInProgress'),
    cancel: t('pattern.cancelPatternsCreation'),
  } : {
    inProgress: t('pattern.patternsUpdateInProgress'),
    cancel: t('pattern.cancelPatternsUpdate'),
  }));

  /**
   * Mirrors `useAiChat` → `emit('update:pending', …)`: `newPending` is in-flight, second arg is
   * `emptyCurrentFormPatterns` (drives creation vs update copy), third is the chat `stop` handler.
   *
   * @param {boolean} newPending - Socket thinking / request pending.
   * @param {boolean|null} [newCreation=null] - Stored on `creation`; when truthy, `pendingTexts` uses the
   *   “creation” strings, when falsy (including `null`) the “update” strings.
   * @param {function} [newCancel] - Replaces the cancel callback (typically `stop` from the chat).
   */
  const updatePending = (newPending, newCreation = null) => {
    pending.value = newPending;
    creation.value = newCreation;
  };

  return {
    pending,
    pendingTexts,

    updatePending,
  };
};

/**
 * Short-lived `jsonString` ref for pushing JSON into the LLM chat from the pattern editor (cleared next tick).
 *
 * @returns {{
 *   jsonString: import('vue').Ref<string>,
 *   updateJsonString: function(string),
 * }}
 */
export const useAiChatJsonString = () => {
  const jsonString = ref('');

  /**
   * Sets `jsonString` then clears it on the next tick so `useAiChatJsonString` in `ai-chat.js` fires once per push.
   *
   * @param {string} newJsonString - Pattern JSON text to surface to the chat watcher.
   */
  const updateJsonString = (newJsonString) => {
    jsonString.value = newJsonString;

    nextTick(() => jsonString.value = '');
  };

  return {
    jsonString,
    updateJsonString,
  };
};

/**
 * Keeps `formPatterns` (normalized map for the sidebar) in sync with the modal form and applies LLM/restored
 * patterns back into the form with validation cleanup.
 *
 * @param {Object} options
 * @param {import('vue').Ref<Object|Array>} options.form - Whole form ref, or scenario `actions` array ref; when
 *   the model has a `patterns` object, pass the full form ref so updates merge into `form.patterns`.
 * @param {import('vue').Ref<string>|string|undefined} options.field - When set, `patternToForm` is used on bulk
 *   updates; otherwise `filterPatternsToForm`.
 * @param {import('vue').Ref<string>|string|undefined} options.context - Passed through to `aiChatFormToPatterns`.
 * @returns {{
 *   patterns: import('vue').Ref<Object>,
 *   updatePatterns: function(Object, string=),
 *   updateFormPatterns: function,
 *   throttledUpdatePatterns: function,
 * }}
 */
export const useAiChatFormPatterns = ({ form, field, context }) => {
  const validator = useValidator();

  const patterns = ref({});

  /**
   * Writes `newPatterns` into the host form: either a single scenario action when `key` matches, or merges into
   * `form.patterns` when present, otherwise spreads pattern fields onto the form root.
   *
   * @param {Object} newPatterns - Normalized pattern fields from the chat.
   * @param {string} [key] - When set, targets the scenario action row with this `key`.
   */
  const updateFormPatterns = (newPatterns, key) => {
    const formRef = form;

    if (key) {
      const index = formRef.value.findIndex(item => item.key === key);
      const formItem = formRef.value[index];

      set(formRef.value, index, {
        ...formItem,
        ...filterPatternsToForm(newPatterns),
      });

      Object.keys(newPatterns).forEach(patternField => validator.errors.clear(`actions.${formItem.key}.${patternField}.json`));

      return;
    }

    const newPatternsForm = unref(field) ? patternToForm(newPatterns) : filterPatternsToForm(newPatterns);

    if (formRef.value.patterns) {
      set(formRef.value, 'patterns', {
        ...formRef.value.patterns,
        ...newPatternsForm,
      });
    } else {
      formRef.value = {
        ...formRef.value,
        ...newPatternsForm,
      };
    }

    Object.keys(newPatterns).forEach(patternField => validator.errors.clear(`${patternField}.json`));
  };

  /**
   * Recomputes `patterns` from the current `form`/`field`/`context` (used directly and via throttle).
   */
  const updatePatterns = () => {
    patterns.value = aiChatFormToPatterns({
      form: unref(form),
      field: unref(field),
      context: unref(context),
    });
  };

  const throttledUpdatePatterns = throttle(updatePatterns, THROTTLE_WAIT_MS);

  return {
    patterns,
    updatePatterns,
    updateFormPatterns,
    throttledUpdatePatterns,
  };
};

/**
 * Lazy-opens the AI sidebar after models load and attaches a deep watcher on `form` to refresh sidebar patterns.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref|import('vue').ComputedRef} [options.form] - Watched while the chat is shown.
 * @param {function} [options.throttledUpdateSidebarConfig] - From `useAiChatFormPatterns` (throttled
 *   `updateFormPatterns`).
 * @returns {{
 *   llms: import('vue').Ref<Array<Object>>,
 *   shown: import('vue').Ref<boolean>,
 *   showChat: function(): Promise,
 *   hideChat: function,
 * }}
 */
export const useAiChatShown = ({ form, throttledUpdatePatterns } = {}) => {
  const { fetchLlms } = useAiChatLlmModel();

  const llms = ref([]);
  const shown = ref(false);

  /**
   * Fetches enabled LLMs; when the list is non-empty, shows the chat and starts a deep `watch` on `form`.
   */
  const showChat = async () => {
    llms.value = await fetchLlms();

    if (!llms.value.length) {
      return;
    }

    shown.value = true;

    watch(form, throttledUpdatePatterns, { deep: true, immediate: true });
  };

  /** Hides the chat UI (does not tear down the `form` watcher started in `showChat`). */
  const hideChat = () => shown.value = false;

  return {
    llms,
    shown,

    showChat,
    hideChat,
  };
};

/**
 * Buffers LLM chat history ids for the current modal session, then bulk-associates them with a persisted rule via
 * `bulkLinkLlmHistory`. Child chat components call `registerChatId` as histories are created; the modal calls
 * `linkChatsWithRuleId` after save with the new or updated rule `_id`.
 *
 * @returns {{
 *   registerChatId: function(string),
 *   linkChatsWithRuleId: function(string),
 * }}
 */
export const useAiChatLinkChats = ({ ruleId } = {}) => {
  const modals = useModals();
  const { bulkLinkLlmHistory } = useLlm();

  const modal = inject('$modal', null);

  /** LLM chat history ids collected through `registerChatId` for a later `linkChatsWithRuleId` bulk update. */
  const chatIds = [];

  /**
   * Appends a chat history id to `chatIds` for a later `linkChatsWithRuleId` call.
   *
   * @param {string} chatId - Server id of the LLM history row to associate with a rule after save.
   */
  const registerChatId = chatId => !chatIds.includes(chatIds) && chatIds.push(chatId);

  /**
   * Bulk-links every id in `chatIds` to the given rule id via `bulkLinkLlmHistory` (e.g. after the modal rule is
   * persisted).
   *
   * @param {string} ruleIdForLink - Rule `_id` written as `rule` on each `{ id, rule }` payload entry.
   * @returns {*} Return value from `bulkLinkLlmHistory({ data })`.
   */
  const linkChatsWithRuleIds = (rulesIdsForLink) => {
    const unwrappedRuleId = unref(ruleId);
    const rulesIds = rulesIdsForLink.filter(ruleIdForLink => ruleIdForLink !== unwrappedRuleId);

    const data = rulesIds.reduce((acc, ruleIdForLink) => {
      chatIds.forEach((chatId) => {
        acc.push({
          _id: chatId,
          rule: ruleIdForLink,
        });
      });

      return acc;
    }, []);

    if (!data.length) {
      return Promise.resolve();
    }

    return bulkLinkLlmHistory({ data });
  };

  /**
   * `afterSubmit` already set on the host modal before this hook wraps it (optional).
   * Called first from the replacement handler (`newAfterSubmit`) so existing save callbacks keep running.
   *
   * @type {Function|undefined}
   */
  const previousAfterSubmit = modal.config.afterSubmit;

  /**
   * Wraps the modal `afterSubmit` handler: runs `previousAfterSubmit` when defined, then links buffered LLM chat
   * ids to the new rule id(s) via `linkChatsWithRuleIds`. Registered on mount by merging into `modal.config`.
   *
   * @param {*} createdItems - Value passed by the modal after save (single id or list of ids / items, depending on
   *   the host). Normalized to an array for `linkChatsWithRuleIds`.
   * @returns {*} Result of `previousAfterSubmit(createdItems)` when that function exists; otherwise `createdItems`.
   */
  const newAfterSubmit = async (createdItems) => {
    const result = await previousAfterSubmit?.(createdItems) ?? createdItems;

    const ids = (isArray(createdItems) ? createdItems : [createdItems])
      .map(item => item._id ?? item.id)
      .filter(Boolean);

    linkChatsWithRuleIds(ids);

    return result;
  };

  onMounted(() => {
    modals.updateModalConfig({ id: modal.id, config: { ...modal.config, afterSubmit: newAfterSubmit } });
  });

  return {
    registerChatId,
  };
};

/**
 * Wires the modal-hosted AI chat: sidebar `options`, pending overlay, pattern sync, `provide('$aiChat')`, and a
 * delayed first fetch of models after mount. `hideChat` runs on unmount (hides the chat flag only).
 *
 * @param {Object} params
 * @param {{ id: string, dialogProps: Object }} params.modal - Passed to `useAiChatMinimized` for padding updates.
 * @param {import('vue').Ref|import('vue').ComputedRef} params.form - Full host form when it owns a `patterns`
 *   object; scenario modals pass the `actions` array ref instead.
 * @param {import('vue').Ref<string>|string|undefined} [params.ruleId] - Optional rule id for the LLM socket.
 * @param {import('vue').Ref<string>|import('vue').ComputedRef<string>|string|undefined} [params.context] - LLM
 *   socket context (e.g. `LLM_SOCKET_CONTEXTS.scenario` or `${LLM_SOCKET_CONTEXTS.widgetFilter}_${type}`).
 * @param {import('vue').Ref<string>|string|undefined} [params.field] - Optional field key for `useAiChatFormPatterns`.
 * @returns {{
 *   shown: import('vue').Ref<boolean>,
 *   options: import('vue').ComputedRef<{ bind: Object, on: Object }>,
 * }}
 */
export const useAiChatForm = ({
  modal,
  form,
  ruleId,
  context,
  field,
}) => {
  const { patternsItems } = useAiChatPatternsItems({ form, context });
  const { initialMinimized, updateMinimized } = useAiChatMinimized({ modal });
  const { pending, pendingTexts, updatePending } = useAiChatPending();
  const { jsonString, updateJsonString } = useAiChatJsonString();
  const {
    patterns,
    updateFormPatterns,
    throttledUpdatePatterns,
  } = useAiChatFormPatterns({ form, field, context });
  const { llms, shown, showChat, hideChat } = useAiChatShown({ form, throttledUpdatePatterns });
  const { registerChatId, linkChatsWithRuleId } = useAiChatLinkChats({ ruleId });

  const options = computed(() => ({
    bind: {
      initialMinimized: initialMinimized.value,
      patterns: patterns.value,
      jsonString: jsonString.value,
      context: unref(context),
      ruleId: unref(ruleId),
      llms: llms.value,
      patternsItems: patternsItems.value,
      pending: pending.value,
      pendingTexts: pendingTexts.value,
    },
    on: {
      'update:minimized': updateMinimized,
      'update:patterns': updateFormPatterns,
      'update:pending': updatePending,
    },
  }));

  provide('$aiChat', {
    llms,
    updateJsonString,
    registerChatId,
  });

  /**
   * Wait for the vuetify modal opening animation to complete before fetching the LLMs.
   */
  onMounted(() => setTimeout(showChat, VUETIFY_ANIMATION_DELAY));
  onBeforeUnmount(hideChat);

  return {
    shown,
    options,
    linkChatsWithRuleId,
  };
};
