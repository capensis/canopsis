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
        v-validate="unitValidateRules",
        :items="availableUnits",
        :label="$t('modals.createWebhook.fields.retryUnit')",
        :error-messages="errors.collect('retryUnit')",
        name="retryUnit",
        clearable,
        hide-details
      )
</template>

<script>
import { PERIODIC_REFRESH_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isRequired() {
      return this.required || Object.values(this.value).some(Boolean);
    },
    availableUnits() {
      return Object.values(PERIODIC_REFRESH_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, 2),
      }));
    },
    fieldValidateRules() {
      return {
        required: this.isRequired,
        numeric: true,
        min_value: 0,
      };
    },
    unitValidateRules() {
      return { required: this.isRequired };
    },
  },
};
</script>
