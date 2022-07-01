<template lang="pug">
  v-layout(row, justify-between, align-center)
    v-flex.theme--light.v-label {{ label }}
    v-flex(xs2)
      c-number-field.mt-0.pt-0(v-field="value", :min="min", hide-details)
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      default: 0,
    },
    label: {
      type: String,
      required: false,
    },
    min: {
      type: Number,
      default: 0,
    },
  },
};
</script>
