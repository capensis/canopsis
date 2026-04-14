import { throttle } from 'lodash';
import {
  computed,
  watch,
  ref,
  unref,
  set,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { SIDE_BARS, LLM_AI_CHAT_WIDTH, LLM_SOCKET_CONTEXTS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';
import { aiChatFormToPatterns } from '@/helpers/entities/llm/chat/form';
import { patternToForm } from '@/helpers/entities/pattern/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSidebar } from '@/hooks/sidebar';
import { useValidator } from '@/hooks/validator/validator';

import { useAiChatLlmModel } from './ai-chat';

const THROTTLE_WAIT_MS = 1000;

export const useAiChatPatternsItems = ({ form = [], context } = {}) => {
  const { t } = useI18n();

  const contextToPatternsItems = {
    [LLM_SOCKET_CONTEXTS.scenario]: () => unref(form)?.map?.((formItem, index) => ({
      text: `${index + 1} - ${t(`scenario.actions.${formItem.type}`)}`,
      value: formItem.key,
    })),
  };

  const patternsItems = computed(() => contextToPatternsItems[context]?.(form) ?? []);

  return {
    patternsItems,
  };
};

/**
 * Opens the AI chat sidebar for a filter (or similar) form, syncs changed pattern blocks to the sidebar
 * config on a throttle, and registers teardown on modal hide / unmount.
 *
 * @param {Object} params
 * @param {import('vue').Ref<Object>} params.form
 * @param {string} params.modalId
 * @param {Function} params.registerOnHide
 * @param {import('vue').Ref|function(): *} [params.ruleId]
 * @param {import('vue').Ref|import('vue').ComputedRef|string} [params.context] - LLM socket context
 *   (e.g. `${LLM_SOCKET_CONTEXTS.widgetFilter}_${type}` for widget filters).
 */
export const useAiChatForm = ({
  form,
  modalId,
  ruleId,
  context,
  field,
}) => {
  const { t } = useI18n();
  const modals = useModals();
  const sidebar = useSidebar();
  const validator = useValidator();

  const pending = ref(false);
  const creation = ref(false);

  let cancel = () => {};

  const pendingTexts = computed(() => (creation.value ? {
    inProgress: t('pattern.patternsCreationInProgress'),
    cancel: t('pattern.cancelPatternsCreation'),
  } : {
    inProgress: t('pattern.patternsUpdateInProgress'),
    cancel: t('pattern.cancelPatternsUpdate'),
  }));

  const { patternsItems } = useAiChatPatternsItems({ form, context });

  const throttledUpdateSidebarConfig = throttle(() => {
    const newConfig = {
      patterns: aiChatFormToPatterns({
        form: unref(form),
        field: unref(field),
        context: unref(context),
      }),
    };

    if (patternsItems.value.length) {
      newConfig.patternsItems = patternsItems.value;
    }

    sidebar.updateConfig({
      id: modalId,
      config: newConfig,
    });
  }, THROTTLE_WAIT_MS);

  const { fetchLlms } = useAiChatLlmModel({ sidebar });

  /**
   * Opens the AI assistant sidebar for this modal instance and passes `socketRoomData`
   * so `ai-chat` can join `SOCKET_ROOMS.llmChat` with the given `context` and optional `rule_id`.
   *
   * Wait for the vuetify modal opening animation to complete before fetching the LLMs.
   */
  const showSidebar = async () => {
    const llms = await fetchLlms();

    if (!llms.length) {
      return;
    }

    sidebar.show({
      id: modalId,
      name: SIDE_BARS.aiChat,
      config: {
        llms,
        overflowYHidden: true,
        minimizable: true,
        width: LLM_AI_CHAT_WIDTH,
        color: 'primary',
        titleIcon: '$vuetify.icons.ai',
        titleMinimized: 'AI',
        setPatterns: (newPatterns, key) => {
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

          formRef.value = {
            ...formRef.value,
            ...(unref(field) ? patternToForm(newPatterns) : filterPatternsToForm(newPatterns)),
          };

          Object.keys(newPatterns).forEach(patternField => validator.errors.clear(`${patternField}.json`));
        },
        setPending: (newPending, newCreation = null, newCancel = () => {}) => {
          pending.value = newPending;
          creation.value = newCreation;
          cancel = newCancel;
        },
        socketRoomData: {
          context: unref(context),
          rule: unref(ruleId),
        },
      },
    });

    watch(form, throttledUpdateSidebarConfig, { deep: true, immediate: true });
  };

  /**
   * Closes the AI assistant sidebar that was opened with the same `id` as this modal.
   */
  const hideSidebar = () => sidebar.hide({ id: modalId });

  onMounted(() => {
    showSidebar();
    modals.registerOnHide({ id: modalId, callback: hideSidebar });
  });

  onBeforeUnmount(hideSidebar);

  return {
    pending,
    pendingTexts,
    cancel,
  };
};
