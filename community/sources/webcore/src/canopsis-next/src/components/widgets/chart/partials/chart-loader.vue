<template>
  <v-progress-linear
    v-if="hasData"
    class="progress-linear-absolute--top"
    height="2"
    indeterminate
  />
  <c-layout
    v-else
    class="pa-4"
    justify="center"
  >
    <v-progress-circular
      color="primary"
      indeterminate
    />
  </c-layout>
</template>

<script>
export default {
  props: {
    hasData: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
