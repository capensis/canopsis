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
    v-flex(xs4)
      v-layout(column)
        v-layout(row, align-center)
          v-flex(:class="{ 'grey--text': !value.enabled }", xs6) {{ $t('common.minimal') }}
          v-flex(xs6)
            v-text-field(
              v-field.number="value.minimal",
              v-validate="minimalRules",
              :error-messages="getErrorMessages(minimalFieldName)",
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
              :error-messages="getErrorMessages(optimalFieldName)",
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
        required: this.value.enabled,
        max_value: this.value.optimal,
      };
    },

    optimalRules() {
      return {
        required: this.value.enabled,
        min_value: this.value.minimal,
      };
    },
  },
  methods: {
    getErrorMessages(name) {
      return this.errors.collect(name, null, false)
        .map((item) => {
          const messageKey = `healthcheck.validation.${item.rule}`;

          return this.$te(messageKey) ? this.$t(messageKey) : item.msg;
        });
    },
  },
};
</script>
