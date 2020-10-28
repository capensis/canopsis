<template lang="pug">
  v-layout(row)
    v-flex(xs8)
      v-text-field(
        v-field.number="value[valueFieldKey]",
        v-validate="'required|numeric|min_value:1'",
        :label="$t('modals.createSnoozeEvent.fields.duration')",
        :error-messages="errors.collect('duration')",
        name="duration",
        type="number"
      )
    v-flex(xs4)
      v-select(
        v-field="value[unitFieldKey]",
        v-validate="'required'",
        :items="availableUnits",
        :error-messages="errors.collect('durationType')",
        name="durationType"
      )
</template>

<script>
import { SNOOZE_DURATION_UNITS, DURATION_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: Object,
      default() {
        return {
          [this.valueFieldKey || 'seconds']: '',
          [this.unitFieldKey || 'unit']: DURATION_UNITS.minute.value,
        };
      },
    },
    valueFieldKey: {
      type: String,
      default: 'seconds',
    },
    unitFieldKey: {
      type: String,
      default: 'unit',
    },
    units: {
      type: Array,
      default: null,
    },
  },
  computed: {
    availableUnits() {
      if (this.units) {
        return this.units;
      }

      return Object.values(SNOOZE_DURATION_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, 2),
      }));
    },
  },
};
</script>
