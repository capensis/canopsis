<template lang="pug">
  v-text-field(
    v-field="value",
    v-validate="rules",
    v-bind="$attrs",
    :label="label || $t('common.name')",
    :error-messages="errors.collect(name)",
    :name="name"
  )
    template(v-if="tooltip", #append)
      c-help-icon(
        :text="tooltip",
        icon="help",
        left
      )
</template>

<script>
export default {
  inject: ['$validator'],
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'name',
    },
    required: {
      type: Boolean,
      default: false,
    },
    tooltip: {
      type: String,
      required: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
};
</script>
