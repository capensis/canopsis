<template>
  <span
    class="text-caption text--secondary"
    :class="{ 'error--text': error }"
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
