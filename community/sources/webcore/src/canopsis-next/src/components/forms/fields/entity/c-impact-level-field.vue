<template>
  <c-number-field
    v-field.number="value"
    :label="label || $t('common.impactLevel')"
    :placeholder="placeholder"
    :single-line="singleLine"
    :disabled="disabled"
    :name="name"
    :max="max"
    :min="min"
    :required="required"
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
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'impact_level',
    },
    singleLine: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      default: 1,
    },
    max: {
      type: Number,
      default: 10,
    },
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    placeholder() {
      return `${this.min}-${this.max}`;
    },
  },
};
</script>
