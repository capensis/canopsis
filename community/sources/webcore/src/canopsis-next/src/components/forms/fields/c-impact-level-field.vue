<template lang="pug">
  v-text-field(
    v-field.number="value",
    v-validate="rules",
    :error-messages="errors.collect(name)",
    :label="label || $t('common.impactLevel')",
    :placeholder="placeholder",
    :single-line="singleLine",
    :name="name",
    :max="max",
    :min="min",
    type="number"
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
      default: 0,
    },
    max: {
      type: Number,
      default: 10,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    placeholder() {
      return `${this.min}-${this.max}`;
    },
    rules() {
      return {
        required: this.required,
        numeric: true,
        min_value: this.min,
        max_value: this.max,
      };
    },
  },
};
</script>
