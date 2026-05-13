<template>
  <!-- Keeps default slot flush with tags — line breaks here become visible spaces in the label -->
  <!-- eslint-disable-next-line vue/singleline-html-element-content-newline -->
  <v-label :aria-required="required" class="text-caption display-block"><slot /></v-label>
</template>

<script>
export default {
  name: 'CLabel',
  props: {
    required: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
