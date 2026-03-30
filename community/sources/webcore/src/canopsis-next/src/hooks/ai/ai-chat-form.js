import {
  computed,
  watch,
  ref,
  unref,
  onMounted,
  onBeforeUnmount,
} from 'vue';
import { throttle } from 'lodash';

import { SIDE_BARS, LLM_AI_CHAT_WIDTH, PATTERNS_FIELDS } from '@/constants';

import { patternToForm } from '@/helpers/entities/pattern/form';
import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSidebar } from '@/hooks/sidebar';

const THROTTLE_WAIT_MS = 1000;

/**
 * Opens the AI chat sidebar for a filter (or similar) form, syncs changed pattern blocks to the sidebar
 * config on a throttle, and registers teardown on modal hide / unmount.
 *
 * @param {Object} params
 * @param {import('vue').Ref<Object>} params.form
 * @param {string} params.modalId
 * @param {Function} params.registerOnHide
 * @param {import('vue').Ref|function(): *} [params.ruleId]
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

  const pending = ref(false);
  const creation = ref(false);
  const cancelPending = ref(null);

  const pendingTexts = computed(() => (creation.value ? {
    inProgress: t('pattern.patternsCreationInProgress'),
    cancel: t('pattern.cancelPatternsCreation'),
  } : {
    inProgress: t('pattern.patternsUpdateInProgress'),
    cancel: t('pattern.cancelPatternsUpdate'),
  }));

  const throttledUpdateSidebarConfig = throttle(() => {
    const unwrappedForm = unref(form);
    const unwrappedField = unref(field);

    let patterns = {};

    if (unwrappedField) {
      patterns = unwrappedForm.groups.length ? { [unwrappedField]: unwrappedForm } : {};
    } else {
      patterns = Object.values(PATTERNS_FIELDS).reduce((acc, patternField) => {
        if (unwrappedForm[patternField]?.groups?.length) {
          acc[patternField] = unwrappedForm[patternField];
        }

        return acc;
      }, {});
    }

    sidebar.updateConfig({
      id: modalId,
      config: { patterns },
    });
  }, THROTTLE_WAIT_MS);

  watch(form, throttledUpdateSidebarConfig, { deep: true });

  /**
   * Opens the AI assistant sidebar for this modal instance and passes `socketRoomData`
   * so `ai-chat` can join `SOCKET_ROOMS.llmChat` with widget-filter context and the current filter id.
   */
  const showSidebar = () => sidebar.show({
    id: modalId,
    name: SIDE_BARS.aiChat,
    config: {
      overflowYHidden: true,
      minimizable: true,
      width: LLM_AI_CHAT_WIDTH,
      color: 'primary',
      titleIcon: '$vuetify.icons.ai',
      titleMinimized: 'AI',
      setPatterns: (newPatterns) => {
        const formRef = form;

        // TODO: analyze if we can use patternToForm or filterPatternsToForm directly
        formRef.value = {
          ...formRef.value,
          ...(unref(field) ? patternToForm(newPatterns) : filterPatternsToForm(newPatterns)),
        };
      },
      setPending: (newPending, newCreation = null, newCancelPending = null) => {
        pending.value = newPending;
        creation.value = newCreation;
        cancelPending.value = newCancelPending;
      },
      socketRoomData: {
        context: unref(context),
        rule_id: unref(ruleId),
      },
    },
  });

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
    cancelPending,
  };
};
