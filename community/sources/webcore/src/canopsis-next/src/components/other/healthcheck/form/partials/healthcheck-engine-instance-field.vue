<template lang="pug">
  v-layout(row, wrap)
    v-flex(xs3)
      v-checkbox(
        v-validate,
        v-field="value.enabled",
        :label="label",
        :error-messages="errors.collect(enabledFieldName)",
        :name="enabledFieldName",
        color="primary"
      )
        c-help-icon(v-if="helpText", slot="append", :text="helpText", max-width="300", top)
    v-flex(xs4)
      v-layout(column)
        v-layout(row, align-center)
          v-flex(:class="{ 'grey--text': !value.enabled }", xs6) {{ $t('common.minimal') }}
          v-flex(xs6)
            v-text-field(
              v-field.number="value.minimal",
              v-validate="minimalRules",
              :error-messages="errors.collect(minimalFieldName)",
              :name="minimalFieldName",
              :label="$t('common.minimal')",
              :disabled="!value.enabled",
              :required="value.enabled",
              :min="0",
              type="number"
            )
        v-layout(row, align-center)
          v-flex(:class="{ 'grey--text': !value.enabled }", xs6) {{ $t('common.optimal') }}
          v-flex(xs6)
            v-text-field.mt-0(
              v-field.number="value.optimal",
              v-validate="optimalRules",
              :error-messages="errors.collect(optimalFieldName)",
              :name="optimalFieldName",
              :label="$t('common.optimal')",
              :disabled="!value.enabled",
              :required="value.enabled",
              :min="0",
              type="number"
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
    helpText: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'engine',
    },
  },
  computed: {
    enabledFieldName() {
      return `${this.name}.enabled`;
    },

    minimalFieldName() {
      return `${this.name}.minimal`;
    },

    optimalFieldName() {
      return `${this.name}.optimal`;
    },

    minimalRules() {
      return {
        required: true,
        max_value: this.value.optimal,
      };
    },

    optimalRules() {
      return {
        required: true,
        min_value: this.value.minimal,
      };
    },
  },
};
</script>
