import {
  ref,
  computed,
  watch,
  unref,
  onMounted,
  onBeforeUnmount,
} from 'vue';
import { pick, throttle, isEqual } from 'lodash';

import { SIDE_BARS, LLM_AI_CHAT_WIDTH } from '@/constants';

import { patternToForm, getChangedPatternsFields } from '@/helpers/entities/pattern/form';
import { filterPatternsToForm } from '@/helpers/entities/filter/form';

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
  const modals = useModals();
  const sidebar = useSidebar();

  const aiChatForm = ref({ ...unref(form) });
  const aiChangedPatternsFields = ref([]);

  const currentChangedPatternsFields = computed(() => {
    const unwrappedField = unref(field);
    const unwrappedForm = unref(form);

    if (unwrappedField) {
      return isEqual(aiChatForm.value, unwrappedForm) ? [] : [unwrappedField];
    }

    const changedPatternsFields = getChangedPatternsFields(
      unwrappedForm,
      aiChatForm.value,
    );

    if (isEqual(aiChangedPatternsFields.value, changedPatternsFields)) {
      return false;
    }

    return changedPatternsFields;
  });

  const throttledUpdateSidebarConfig = throttle(() => {
    const changedPatternsFields = currentChangedPatternsFields.value;

    if (!changedPatternsFields) {
      return;
    }

    const unwrappedForm = unref(form);
    const unwrappedField = unref(field);

    aiChangedPatternsFields.value = changedPatternsFields;

    let patterns = {};

    if (unwrappedField) {
      patterns = changedPatternsFields.length ? { [unwrappedField]: unwrappedForm } : {};
    } else {
      patterns = pick(unwrappedForm, changedPatternsFields);
    }

    sidebar.updateConfig({
      id: modalId,
      config: { patterns },
    });
  }, THROTTLE_WAIT_MS);

  watch(() => unref(form), throttledUpdateSidebarConfig, { deep: true });

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
};
