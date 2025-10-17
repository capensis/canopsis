<template>
  <span
    :class="{ 'error--text': error }"
    class="text-caption text--secondary"
  >
    <slot />
  </span>
</template>

<script>
export default {
  props: {
    error: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
