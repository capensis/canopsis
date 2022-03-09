<template lang="pug">
  v-text-field(
    v-field="value.field",
    v-validate="'required'",
    :disabled="disabled",
    :label="label || $t('common.field')",
    :error-messages="errors.collect(name)",
    :name="name",
    hide-details
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
      type: Object,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'extra_infos',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
