<template>
  <v-layout class="mb-3 remediation-instruction-stats-summary-row">
    <v-flex class="remediation-instruction-stats-summary-row--label">
      <span class="font-weight-bold">{{ label }}</span>
    </v-flex>
    <v-flex>
      <slot>{{ value }}</slot>
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    value: {
      type: [String, Number],
      required: false,
    },
  },
};
</script>

<style lang="scss">
.remediation-instruction-stats-summary-row {
  &--label {
    max-width: 200px;
  }
}
</style>
