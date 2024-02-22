<template>
  <span
    :class="{ 'error--text': error }"
    class="text-caption text--secondary"
  >{{ message }}</span>
</template>

<script>
export default {
  props: {
    message: {
      type: String,
      required: true,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
