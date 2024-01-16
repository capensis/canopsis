<template>
  <c-number-field
    v-field="value"
    :label="label || $t('common.priority')"
    :disabled="disabled"
    :loading="loading"
    :name="name"
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
    loading: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
