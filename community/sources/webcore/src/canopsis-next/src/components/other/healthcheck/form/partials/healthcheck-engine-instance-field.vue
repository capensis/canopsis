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
            c-number-field(
              v-field="value.minimal",
              :error-messages="getErrorMessages(minimalFieldName)",
              :name="minimalFieldName",
              :label="$t('common.minimal')",
              :disabled="!value.enabled",
              :required="value.enabled",
              :max="+value.optimal",
              :min="0"
            )
        v-layout(row, align-center)
          v-flex(:class="{ 'grey--text': !value.enabled }", xs6) {{ $t('common.optimal') }}
          v-flex(xs6)
            c-number-field.mt-0(
              v-field="value.optimal",
              :error-messages="getErrorMessages(optimalFieldName)",
              :name="optimalFieldName",
              :label="$t('common.optimal')",
              :disabled="!value.enabled",
              :required="value.enabled",
              :min="+value.minimal"
            )
    v-flex(xs9)
      v-messages(:value="errors.collect(name)", color="error")
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
