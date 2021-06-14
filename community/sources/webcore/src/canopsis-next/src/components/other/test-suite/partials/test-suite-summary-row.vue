<template lang="pug">
  v-layout.test-suite-summary-row.pt-2(column)
    v-layout(row)
      v-flex.test-suite-summary-row--label
        h4.subheading.font-weight-bold {{ label }}
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
.test-suite-summary-row {
  &--label {
    max-width: 160px;
  }
}
</style>
