<template lang="pug">
  v-text-field(
    v-field.number="value",
    v-validate="rules",
    :label="label || $t('common.priority')",
    :error-messages="errors.collect(name)",
    :disabled="disabled",
    :name="name",
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
      default: 'priority',
    },
    min: {
      type: Number,
      default: 0,
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
    rules() {
      return {
        required: this.required,
        numeric: true,
        min_value: this.min,
      };
    },
  },
};
</script>
