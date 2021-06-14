<template lang="pug">
  v-layout(row)
    v-flex.pr-3(xs8)
      v-text-field(
        v-field.number="duration.value",
        v-validate="durationValidateRules",
        :label="label || $t('common.duration')",
        :error-messages="errors.collect(intervalName)",
        :disabled="disabled",
        :name="intervalName",
        :min="min",
        type="number"
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

import { AVAILABLE_TIME_UNITS } from '@/constants';

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
  },
  computed: {
    intervalName() {
      return `${this.name}.seconds`;
    },

    unitFieldName() {
      return `${this.name}.unit`;
    },

    availableUnits() {
      if (this.units) {
        return this.units;
      }

      return Object.values(AVAILABLE_TIME_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.duration.value || 0),
      }));
    },

    isRequired() {
      return this.required || isNumber(this.duration.value) || Boolean(this.duration.unit);
    },

    durationValidateRules() {
      return {
        required: this.isRequired,
        numeric: true,
        min_value: this.min,
      };
    },

    unitValidateRules() {
      return { required: this.isRequired };
    },
  },
};
</script>
