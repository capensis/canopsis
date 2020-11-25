<template lang="pug">
  td
    slot(v-if="!pending")
    v-progress-circular.ml-2(v-else, indeterminate, color="primary", size="20")
</template>

<script>
export default {
  props: {
    pending: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
