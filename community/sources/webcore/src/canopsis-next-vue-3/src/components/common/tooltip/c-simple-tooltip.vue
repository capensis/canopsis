<template>
  <v-tooltip :disabled="disabled || !content" :location="top ? 'top' : 'bottom'">
    <template #activator="{ props }">
      <slot :on="props" :props="props" name="activator" />
    </template>
    <span v-if="html" v-html="sanitizedContent" />
    <span v-else>{{ content }}</span>
  </v-tooltip>
</template>

<script>
import { computed } from 'vue';

import { sanitizeHtml } from '@/helpers/html';

export default {
  props: {
    content: {
      type: String,
      default: '',
    },
    top: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    html: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const sanitizedContent = computed(() => sanitizeHtml(props.content));

    return { sanitizedContent };
  },
};
</script>
