<template lang="pug">
  v-radio-group(
    v-field="value",
    hide-details,
    mandatory,
    row,
    color="primary"
  )
    v-radio.radio-operator(
      :label="$t('common.and')",
      :value="$constants.FILTER_MONGO_OPERATORS.and",
      data-test="radioAnd",
      color="primary"
    )
    v-radio.radio-operator(
      :label="$t('common.or')",
      :value="$constants.FILTER_MONGO_OPERATORS.or",
      data-test="radioOr",
      color="primary"
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
      type: String,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
  .radio-operator {
    text-transform: uppercase;
  }
</style>
