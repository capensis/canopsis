<template>
  <c-number-field
    v-field="value"
    :label="label"
    :name="name"
    :max="max"
    :min="min"
    :hide-details="hideDetails"
    prefix="%"
    required
  />
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
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
