<template>
  <c-alert
    v-if="isVisible"
    class="ma-0"
    type="warning"
    dismissible
    @input="dismiss"
  >
    <div
      v-html="messageHtml"
      class="font-weight-regular"
    />
  </c-alert>
</template>

<script>
import { computed, ref } from 'vue';

import { LLM_AI_CHAT_TOURS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useTourBase } from '@/hooks/tour/tour-base';

export default {
  setup() {
    const { t } = useI18n();
    const { currentUser, updateCurrentUserTours } = useTourBase();

    const dismissedLocally = ref(false);

    const messageHtml = computed(() => t('llm.chat.warningAlert'));

    const isVisible = computed(() => (
      !dismissedLocally.value
      && !currentUser.value?.ui_tours?.[LLM_AI_CHAT_TOURS.hiddenWarningAlert]
    ));

    /**
     * Persists `hiddenWarningAlert` on the current user when the warning is dismissed.
     *
     * @param {boolean} visible - `false` when the user closes the alert (Vuetify dismissible `input`).
     */
    const dismiss = async (visible) => {
      if (visible) {
        return;
      }

      dismissedLocally.value = true;

      updateCurrentUserTours({ data: { [LLM_AI_CHAT_TOURS.hiddenWarningAlert]: true } });
    };

    return {
      messageHtml,
      isVisible,
      dismiss,
    };
  },
};
</script>
