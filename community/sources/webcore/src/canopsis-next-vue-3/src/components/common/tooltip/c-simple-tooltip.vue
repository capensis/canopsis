<template>
  <v-tooltip
    v-model="opened"
    :disabled="disabled || !content"
    :location="top ? 'top' : 'bottom'"
  >
    <template #activator="{ props }">
      <slot :on="props" :props="props" name="activator" />
    </template>
    <template v-if="html">
      <span v-if="opened" v-html="sanitizedContent" />
    </template>
    <span v-else>{{ content }}</span>
  </v-tooltip>
</template>

<script>
import { computed, ref } from 'vue';

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
    const opened = ref(false);
    const sanitizedContent = computed(() => (opened.value ? sanitizeHtml(props.content) : ''));

    return {
      opened,
      sanitizedContent,
    };
  },
};
</script>
