<template>
  <c-alert
    class="ma-0"
    type="info"
  >
    <div
      v-html="messageHtml"
      class="font-weight-regular"
    />
  </c-alert>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    patternItem: {
      type: Object,
      default: null,
    },
    context: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const messageHtml = computed(() => `${t(`llm.chat.infoAlerts.${props.context}`, {
      patternItem: props.patternItem?.text,
    })}${t('llm.chat.infoAlertEnding')}`);

    return {
      messageHtml,
    };
  },
};
</script>
