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
  },
  setup(props) {
    const { t } = useI18n();

    const messageHtml = computed(() => t('llm.chat.infoAlert', {
      patternItem: props.patternItem?.text,
    }));

    return {
      messageHtml,
    };
  },
};
</script>
