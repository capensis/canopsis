<template lang="pug" functional>
  span
    v-icon(v-if="props.value", color="primary") check_circle
    v-icon(v-else, color="error") cancel
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
