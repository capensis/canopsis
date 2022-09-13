<template lang="pug">
  v-text-field(
    v-bind="$attrs",
    v-field.number="value",
    v-validate="rules",
    :label="label",
    :error-messages="errors.collect(name)",
    :disabled="disabled",
    :name="name",
    :min="min",
    :max="max",
    type="number"
  )
    template(#append="")
      slot(name="append")
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
      default: 'value',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      default: undefined,
    },
    max: {
      type: Number,
      default: undefined,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
        numeric: true,
        min_value: this.min ?? false,
        max_value: this.max ?? false,
      };
    },
  },
};
</script>
