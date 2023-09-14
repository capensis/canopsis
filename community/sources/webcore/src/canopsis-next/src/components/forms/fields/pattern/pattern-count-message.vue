<template lang="pug">
  span.caption.text--secondary(:class="{ 'error--text': error }") {{ message }}
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
