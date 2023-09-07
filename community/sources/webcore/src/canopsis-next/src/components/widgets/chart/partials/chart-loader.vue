<template lang="pug">
  v-fade-transition(v-if="hasMetrics", key="progress", mode="out-in")
    v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
  v-layout.pa-4(v-else, justify-center)
    v-progress-circular(color="primary", indeterminate)
</template>

<script>
export default {
  props: {
    hasMetrics: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
