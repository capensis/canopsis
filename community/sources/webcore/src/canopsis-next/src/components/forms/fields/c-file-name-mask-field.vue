<template lang="pug">
  v-layout(row, wrap)
    v-flex(v-show="title", xs12)
      v-layout
        h4.subheading.grey--text.text--darken-2 {{ title }}
        c-help-icon(v-if="helpText", :text="helpText", icon-class="ml-2 storage-help-tooltip", max-width="250", top)
    v-flex(xs12)
      v-layout(column)
        v-text-field(
          v-field="value",
          v-validate="rules",
          :value="value",
          :label="label",
          :placeholder="placeholder",
          :error-messages="errors.collect(name)",
          :disabled="disabled",
          :name="name"
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
      type: String,
      default: '',
    },
    title: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'mask',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    helpText: {
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
