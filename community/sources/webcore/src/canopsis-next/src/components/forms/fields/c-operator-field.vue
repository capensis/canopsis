<template lang="pug">
  v-radio-group(
    v-field="value",
    :name="name",
    color="primary",
    hide-details,
    mandatory,
    row
  )
    v-radio.radio-operator(
      :label="$t('common.and')",
      :value="$constants.FILTER_MONGO_OPERATORS.and",
      color="primary"
    )
    v-radio.radio-operator(
      :label="$t('common.or')",
      :value="$constants.FILTER_MONGO_OPERATORS.or",
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
    name: {
      type: String,
      default: 'operator',
    },
  },
};
</script>

<style lang="scss" scoped>
  .radio-operator {
    text-transform: uppercase;
  }
</style>
