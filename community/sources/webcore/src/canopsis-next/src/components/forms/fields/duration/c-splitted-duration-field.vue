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
import { computed } from 'vue';

import { TIME_UNITS } from '@/constants';

import { fromSeconds, toSeconds } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form';

export default {
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
  setup(props, { emit }) {
    const { tc } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const hasDayUnit = computed(
      () => (props.maxValue ? fromSeconds(props.maxValue, TIME_UNITS.day) > 1 : true),
    );

    const availableItems = computed(() => {
      const items = [];

      if (hasDayUnit.value) {
        items.push(
          {
            unit: TIME_UNITS.day,
            label: tc('common.times.day', 2),
          },
        );
      }

      items.push(
        {
          unit: TIME_UNITS.hour,
          label: tc('common.times.hour', 2),
        },
        {
          unit: TIME_UNITS.minute,
          label: tc('common.times.minute', 2),
        },
      );

      return items;
    });

    const convertSecondsToRoundedUnitValue = (value, unit) => Math.floor(fromSeconds(value, unit));

    const valuesByUnit = computed(() => {
      const { values } = availableItems.value.reduce((acc, { unit }) => {
        const unitValue = convertSecondsToRoundedUnitValue(acc.restSeconds, unit);

        acc.values[unit] = unitValue;
        acc.restSeconds -= toSeconds(unitValue, unit);

        return acc;
      }, {
        restSeconds: props.value,
        values: {},
      });

      return values;
    });

    const maxValuesByUnit = computed(() => {
      if (!props.maxValue) {
        return {};
      }

      const { values } = availableItems.value.reduce((acc, { unit }) => {
        acc.values[unit] = convertSecondsToRoundedUnitValue(acc.restSeconds, unit);
        acc.restSeconds -= toSeconds(valuesByUnit.value[unit], unit);

        return acc;
      }, {
        restSeconds: props.maxValue,
        values: {},
      });

      return values;
    });

    const updateValueByUnit = (value, unit) => {
      const preparedValue = Math.max(value || 0, 0);
      const newValue = props.value + toSeconds(preparedValue - valuesByUnit.value[unit], unit);

      updateModel(props.maxValue ? Math.min(props.maxValue, newValue) : newValue);
    };

    return {
      availableItems,
      valuesByUnit,
      maxValuesByUnit,
      updateValueByUnit,
    };
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
