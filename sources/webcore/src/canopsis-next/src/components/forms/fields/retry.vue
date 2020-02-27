<template lang="pug">
  v-layout
    v-flex(xs4)
      v-text-field(
        v-field.number="value.count",
        v-validate="fieldValidateRules",
        :label="$t('modals.createWebhook.fields.retryCount')",
        :error-messages="errors.collect('retryCount')",
        :min="0",
        name="retryCount",
        type="number"
      )
    v-flex(xs4)
      v-text-field(
        v-field.number="value.delay",
        v-validate="fieldValidateRules",
        :label="$t('modals.createWebhook.fields.retryDelay')",
        :error-messages="errors.collect('retryDelay')",
        :min="0",
        name="retryDelay",
        type="number"
      )
    v-flex(xs4)
      v-select(
        v-field="value.unit",
        v-validate="'required'",
        :items="availableUnits",
        :error-messages="errors.collect('retryUnit')",
        name="retryUnit",
        hide-details
      )
</template>

<script>
import { PERIODIC_REFRESH_UNITS, DEFAULT_RETRY_FIELD } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: Object,
      default: () => ({ ...DEFAULT_RETRY_FIELD }),
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableUnits() {
      return Object.values(PERIODIC_REFRESH_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, 2),
      }));
    },
    fieldValidateRules() {
      return `${this.required ? 'required|' : ''}numeric|min_value:0`;
    },
  },
};
</script>
