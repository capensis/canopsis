<template>
  <v-layout class="gap-3" align-end>
    <c-number-field
      v-for="{ unit, label } of availableItems"
      :key="unit"
      :label="label"
      :value="valuesByUnit[unit]"
      :min="0"
      :max="maxValuesByUnit[unit]"
      hide-details
      @input="updateValueByUnit($event, unit)"
    />
  </v-layout>
</template>

<script>
import { TIME_UNITS } from '@/constants';

import { fromSeconds, toSeconds } from '@/helpers/date/duration';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  props: {
    value: {
      type: Number,
      required: true,
    },
    maxValue: {
      type: Number,
      required: false,
    },
  },
  computed: {
    hasDayUnit() {
      return fromSeconds(this.maxValue, TIME_UNITS.day) > 1;
    },

    availableItems() {
      const items = [];

      if (this.hasDayUnit) {
        items.push(
          {
            unit: TIME_UNITS.day,
            label: this.$tc('common.times.day', 2),
          },
        );
      }

      items.push(
        {
          unit: TIME_UNITS.hour,
          label: this.$tc('common.times.hour', 2),
        },
        {
          unit: TIME_UNITS.minute,
          label: this.$tc('common.times.minute', 2),
        },
      );

      return items;
    },

    maxValuesByUnit() {
      return this.availableItems.reduce((acc, { unit }) => {
        acc[unit] = Math.floor(fromSeconds(this.maxValue, unit));

        return acc;
      }, {});
    },

    valuesByUnit() {
      return this.convertValueToMaxValueByUnit(this.value);
    },
  },
  methods: {
    convertValueToMaxValueByUnit(value = 0) {
      const { values } = this.availableItems.reduce((acc, { unit }) => {
        const unitValue = Math.floor(fromSeconds(acc.restSeconds, unit));

        acc.values[unit] = unitValue;
        acc.restSeconds -= toSeconds(unitValue, unit);

        return acc;
      }, {
        restSeconds: value,
        values: {},
      });

      return values;
    },

    updateValueByUnit(value, unit) {
      const preparedValue = Math.max(value || 0, 0);
      const newValue = this.value + toSeconds(preparedValue - this.valuesByUnit[unit], unit);

      this.updateModel(this.maxValue ? Math.min(this.maxValue, newValue) : newValue);
    },
  },
};
</script>
