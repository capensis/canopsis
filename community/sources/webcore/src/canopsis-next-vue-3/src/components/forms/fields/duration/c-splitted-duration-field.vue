<template>
  <c-layout class="c-splitted-duration-field" gap="3" align="end">
    <c-number-field
      v-for="{ unit, label } of availableItems"
      :key="unit"
      :label="label"
      :value="valuesByUnit[unit]"
      :min="0"
      :max="maxValuesByUnit[unit]"
      hide-details
      @update:value="updateValueByUnit($event, unit)"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { TIME_UNITS } from '@/constants';

import { fromSeconds, toSeconds } from '@/helpers/date/duration';

import { useModelField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Number,
      required: true,
    },
    maxValue: {
      type: Number,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const hasDayUnit = computed(() => props.maxValue ? fromSeconds(props.maxValue, TIME_UNITS.day) > 1 : true,
    );

    const availableItems = computed(() => {
      const items = [];

      if (hasDayUnit.value) {
        items.push({
          unit: TIME_UNITS.day,
          label: t('common.times.day', 2),
        });
      }

      items.push(
        {
          unit: TIME_UNITS.hour,
          label: t('common.times.hour', 2),
        },
        {
          unit: TIME_UNITS.minute,
          label: t('common.times.minute', 2),
        },
      );

      return items;
    });

    const convertSecondsToRoundedUnitValue = (value, unit) => Math.floor(fromSeconds(value, unit));

    const valuesByUnit = computed(() => {
      const { values } = availableItems.value.reduce(
        (acc, { unit }) => {
          const unitValue = convertSecondsToRoundedUnitValue(acc.restSeconds, unit);

          acc.values[unit] = unitValue;
          acc.restSeconds -= toSeconds(unitValue, unit);

          return acc;
        },
        {
          restSeconds: props.modelValue,
          values: {},
        },
      );

      return values;
    });

    const maxValuesByUnit = computed(() => {
      if (!props.maxValue) {
        return {};
      }

      const { values } = availableItems.value.reduce(
        (acc, { unit }) => {
          acc.values[unit] = convertSecondsToRoundedUnitValue(acc.restSeconds, unit);
          acc.restSeconds -= toSeconds(valuesByUnit.value[unit], unit);

          return acc;
        },
        {
          restSeconds: props.maxValue,
          values: {},
        },
      );

      return values;
    });

    const updateValueByUnit = (value, unit) => {
      const preparedValue = Math.max(value || 0, 0);
      const newValue = props.modelValue + toSeconds(preparedValue - valuesByUnit.value[unit], unit);

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
