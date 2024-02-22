<template>
  <v-fade-transition
    v-if="hasMetrics"
    key="progress"
    mode="out-in"
  >
    <v-progress-linear
      class="progress-linear-absolute--top"
      height="2"
      indeterminate
    />
  </v-fade-transition>
  <v-layout
    v-else
    class="pa-4"
    justify-center
  >
    <v-progress-circular
      color="primary"
      indeterminate
    />
  </v-layout>
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
