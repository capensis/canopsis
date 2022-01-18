<template lang="pug">
  v-text-field(
    v-field.number="value",
    v-validate="'required|numeric|min_value:0|max_value:100'",
    :label="label",
    :error-messages="errors.collect(name)",
    :name="name",
    :max="max",
    :min="min",
    type="number",
    prefix="%"
  )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Number, String],
      default: 0,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      required: false,
    },
    max: {
      type: Number,
      default: 100,
    },
    min: {
      type: Number,
      default: 0,
    },
  },
};
</script>
