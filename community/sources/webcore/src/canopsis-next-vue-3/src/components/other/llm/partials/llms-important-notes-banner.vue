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

import { sanitizeHtml } from '@/helpers/html';

import { useI18n } from '@/hooks/i18n';

export default {
  setup() {
    const { t } = useI18n();

    const sanitizedImportantNotesMessage = computed(() => sanitizeHtml(
      t('llm.importantNotesMessage'),
    ));

    return {
      sanitizedImportantNotesMessage,
    };
  },
};
</script>
