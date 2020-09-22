<template lang="pug">
  v-rating.layout(
    v-field="value",
    :readonly="readonly",
    backgroundColor="primary",
    medium,
    half-increments
  )
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Number, String],
      default: 0,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
