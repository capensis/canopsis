<template lang="pug">
  v-layout(row)
    v-textarea(
      v-field="output",
      v-validate="required ? 'required' : ''",
      :label="label || $t('common.output')",
      :error-messages="errors.collect(name)",
      :name="name"
    )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'output',
    event: 'input',
  },
  props: {
    output: {
      type: String,
      required: true,
    },
    required: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'output',
    },
  },
};
</script>
