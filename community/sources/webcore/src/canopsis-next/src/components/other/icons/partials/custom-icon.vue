<template>
  <span v-html="sanitizedContent" class="custom-icon" />
</template>

<script>
import { computed } from 'vue';

import { sanitizeHtml } from '@/helpers/html';

const CUSTOM_ICON_SANITIZE_OPTIONS = {
  USE_PROFILES: {
    svg: true,
    svgFilters: true,
  },

  ADD_ATTR: ['class', 'style', 'viewBox', 'xmlns'],

  FORBID_TAGS: ['script', 'foreignObject', 'iframe', 'object', 'embed'],
};

export default {
  props: {
    content: {
      type: String,
      default: '',
    },
  },

  setup(props) {
    const sanitizedContent = computed(() => sanitizeHtml(props.content, CUSTOM_ICON_SANITIZE_OPTIONS),
    );

    return {
      sanitizedContent,
    };
  },
};
</script>

<style lang="scss">
.custom-icon {
  display: inline-flex;
  align-items: center;
  max-width: 24px;
  max-height: 24px;
  overflow: hidden;

  svg {
    width: 24px;
    height: 24px;
  }
}
</style>
