<template>
  <!-- Keeps default slot flush with tags — line breaks here become visible spaces in the label -->
  <!-- eslint-disable-next-line vue/singleline-html-element-content-newline -->
  <v-label :aria-required="required" :class="labelClasses"><slot /></v-label>
</template>

<script>
import { computed } from 'vue';

export default {
  name: 'CLabel',
  props: {
    required: {
      type: Boolean,
      default: false,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const labelClasses = computed(() => ({
      'error--text': props.error,
      'text-caption display-block': true,
    }));

    return {
      labelClasses,
    };
  },
};
</script>
