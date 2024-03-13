<template>
  <v-layout class="c-splitted-duration-field gap-3" align-end>
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
      return this.maxValue
        ? fromSeconds(this.maxValue, TIME_UNITS.day) > 1
        : true;
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

    valuesByUnit() {
      const { values } = this.availableItems.reduce((acc, { unit }) => {
        const unitValue = this.convertSecondsToRoundedUnitValue(acc.restSeconds, unit);

        acc.values[unit] = unitValue;
        acc.restSeconds -= toSeconds(unitValue, unit);

        return acc;
      }, {
        restSeconds: this.value,
        values: {},
      });

      return values;
    },

    maxValuesByUnit() {
      if (!this.maxValue) {
        return {};
      }

      const { values } = this.availableItems.reduce((acc, { unit }) => {
        acc.values[unit] = this.convertSecondsToRoundedUnitValue(acc.restSeconds, unit);
        acc.restSeconds -= toSeconds(this.valuesByUnit[unit], unit);

        return acc;
      }, {
        restSeconds: this.maxValue,
        values: {},
      });

      return values;
    },
  },
  methods: {
    updateValueByUnit(value, unit) {
      const preparedValue = Math.max(value || 0, 0);
      const newValue = this.value + toSeconds(preparedValue - this.valuesByUnit[unit], unit);

      this.updateModel(this.maxValue ? Math.min(this.maxValue, newValue) : newValue);
    },

    convertSecondsToRoundedUnitValue(value, unit) {
      return Math.floor(fromSeconds(value, unit));
    },
  },
};
</script>

<style lang="scss">
.c-splitted-duration-field {
  > * {
    flex: 1;
  }
}
</style>
