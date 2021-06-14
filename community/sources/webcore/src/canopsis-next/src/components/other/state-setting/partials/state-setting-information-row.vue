<template lang="pug">
  v-layout.mb-2.state-setting-information-row(row)
    v-flex.state-setting-information-row--label
      span.font-weight-bold {{ label }}
    v-flex
      slot {{ value }}{{ unit }}
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
    unit: {
      type: String,
      required: false,
    },
  },
};
</script>

<style lang="scss">
.state-setting-information-row {
  &--label {
    max-width: 100px;
  }
}
</style>
