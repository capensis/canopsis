<template lang="pug">
  v-layout
    v-flex(xs6)
      v-text-field(
        v-field.number="value.interval",
        v-validate="fieldValidateRules",
        :label="label",
        :error-messages="errors.collect('retryDelay')",
        :disabled="disabled",
        :min="0",
        name="popupTimeout",
        type="number"
      )
    v-flex(xs6)
      v-select(
        v-field="value.unit",
        v-validate="unitValidateRules",
        :items="availableUnits",
        :label="$t('parameters.userInterfaceForm.fields.popupTimeoutUnit')",
        :error-messages="errors.collect('popupTimeoutUnit')",
        :disabled="disabled",
        name="popupTimeoutUnit"
      )
</template>

<script>
import { PERIODIC_REFRESH_UNITS, TIME_UNITS } from '@/constants';
import { POPUP_AUTO_CLOSE_DELAY } from '@/config';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: Object,
      default: () => ({
        interval: POPUP_AUTO_CLOSE_DELAY / 1000,
        unit: TIME_UNITS.second,
      }),
    },
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      required: true,
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
      return {
        required: this.required,
        numeric: true,
        min_value: 0,
      };
    },
    unitValidateRules() {
      return { required: this.required };
    },
  },
};
</script>
