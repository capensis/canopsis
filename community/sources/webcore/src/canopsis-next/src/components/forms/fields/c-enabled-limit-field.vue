<template lang="pug">
  v-layout(row, wrap)
    v-flex(xs5)
      v-checkbox(
        v-validate,
        v-field="value.enabled",
        :label="label",
        :error-messages="errors.collect(enabledFieldName)",
        :name="enabledFieldName",
        color="primary"
      )
        c-help-icon(v-if="helpText", slot="append", :text="helpText", max-width="300", top)
    v-flex(xs2)
      c-number-field(
        v-field="value.limit",
        :label="fieldLabel",
        :name="limitFieldName",
        :disabled="!value.enabled",
        :required="value.enabled",
        :min="min"
      )
    v-flex(xs9)
      div.v-messages.theme--light.error--text
        div.v-messages__wrapper
          div.v-messages__message(v-for="error in errors.collect(name)", :key="error") {{ error }}
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
    fieldLabel: {
      type: String,
      default: '',
    },
    helpText: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'limit',
    },
    min: {
      type: Number,
      required: false,
    },
  },
  computed: {
    enabledFieldName() {
      return `${this.name}.enabled`;
    },

    limitFieldName() {
      return `${this.name}.limit`;
    },
  },
};
</script>
