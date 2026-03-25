<template>
  <c-alert type="info">
    <div
      v-html="sanitizedImportantNotesMessage"
      class="llms-important-notes-banner__content font-weight-regular"
    />
  </c-alert>
</template>

<script>
import { computed } from 'vue';

import { GEMINI_AI_CONSOLE_BILLING_URL } from '@/config';

import { useI18n } from '@/hooks/i18n';

export default {
  setup() {
    const { t } = useI18n();

    const sanitizedImportantNotesMessage = computed(() => t('llm.importantNotesMessage', { geminiConsoleUrl: GEMINI_AI_CONSOLE_BILLING_URL }));

    return {
      sanitizedImportantNotesMessage,
    };
  },
};
</script>

<style lang="scss" scoped>
.llms-important-notes-banner__content ::v-deep a {
  color: var(--v-info-base) !important;
}
</style>
