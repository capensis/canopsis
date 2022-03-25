<template lang="pug">
  v-layout(row)
    v-flex.pr-3(xs8)
      c-number-field(
        v-field="duration.value",
        :label="label || $t('common.duration')",
        :disabled="disabled",
        :name="valueFieldName",
        :required="isRequired",
        :min="min"
      )
    v-flex(xs4)
      v-select(
        v-field="duration.unit",
        v-validate="unitValidateRules",
        :items="availableUnits",
        :disabled="disabled",
        :label="unitsLabel",
        :error-messages="errors.collect(unitFieldName)",
        :name="unitFieldName",
        :clearable="clearable"
      )
</template>

<script>
import { isNumber } from 'lodash';

import { AVAILABLE_TIME_UNITS, SHORT_AVAILABLE_TIME_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'duration',
    event: 'input',
  },
  props: {
    duration: {
      type: Object,
      default: () => ({
        value: 1,
        unit: AVAILABLE_TIME_UNITS.minute.value,
      }),
    },
    label: {
      type: String,
      default: null,
    },
    unitsLabel: {
      type: String,
      default: null,
    },
    units: {
      type: Array,
      default: null,
    },
    name: {
      type: String,
      default: 'duration',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      default: 1,
    },
    long: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    valueFieldName() {
      return `${this.name}.value`;
    },

    unitFieldName() {
      return `${this.name}.unit`;
    },

    availableUnits() {
      if (this.units) {
        return this.units;
      }

      const units = this.long
        ? AVAILABLE_TIME_UNITS
        : SHORT_AVAILABLE_TIME_UNITS;

      return Object.values(units).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.duration.value || 0),
      }));
    },

    isRequired() {
      return this.required || isNumber(this.duration.value) || Boolean(this.duration.unit);
    },

    unitValidateRules() {
      return { required: this.isRequired };
    },
  },
};
</script>
