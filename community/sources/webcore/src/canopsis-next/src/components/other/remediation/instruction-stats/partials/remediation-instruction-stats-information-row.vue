<template lang="pug">
  v-layout.mb-3.remediation-instruction-stats-information-row(row)
    v-flex.remediation-instruction-stats-information-row--label
      span.font-weight-bold {{ label }}
    v-flex
      slot {{ value }}
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
.remediation-instruction-stats-information-row {
  &--label {
    max-width: 200px;
  }
}
</style>
