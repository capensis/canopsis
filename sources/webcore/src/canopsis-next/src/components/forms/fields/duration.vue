<template lang="pug">
  v-layout(row)
    v-flex(xs8)
      v-text-field(
        v-field.number="duration.value",
        v-validate="'required|numeric|min_value:1'",
        :label="label || $t('common.duration')",
        :error-messages="errors.collect(name)",
        :name="name",
        type="number"
      )
    v-flex(xs4)
      v-select(
        v-field="duration.unit",
        v-validate="'required'",
        :items="availableUnits",
        :error-messages="errors.collect(unitFieldName)",
        :name="unitFieldName"
      )
</template>

<script>
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
        value: 0,
        unit: AVAILABLE_TIME_UNITS.minute.value,
      }),
    },
    label: {
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
  },
  computed: {
    unitFieldName() {
      return `${this.name}Unit`;
    },
    availableUnits() {
      if (this.units) {
        return this.units;
      }

      return Object.values(AVAILABLE_TIME_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, 2),
      }));
    },
  },
};
</script>
