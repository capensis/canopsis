<template>
  <c-select
    :model-value="value"
    :items="intervals"
    :label="label || t('common.period')"
    :name="name"
    item-title="text"
    item-value="value"
    @update:model-value="updateValue"
  />
</template>

<script>
import { computed } from 'vue';

import { TIME_UNITS } from '@/constants';

import { getNowIntervalValueForHours } from '@/helpers/date/date-intervals';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'interval',
    },
  },

  emits: ['input', 'update:value'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const intervals = computed(() => {
      const result = [];
      const nowHours = (new Date()).getHours();

      for (let index = 0; index < 24; index += 1) {
        const calculatedHours = nowHours - index;
        const hours = calculatedHours >= 0 ? calculatedHours : 24 + calculatedHours;
        const isYesterday = nowHours < index;

        result.push({
          value: getNowIntervalValueForHours(index, TIME_UNITS.hour),
          text: `${hours}:00 - ${hours}:59${
            isYesterday ? ` (${t('common.yesterday').toLowerCase()})` : ''
          }`,
        });
      }

      return result;
    });

    const updateValue = (value) => {
      emit('input', value);
      emit('update:value', value);
    };

    return {
      t,
      intervals,
      updateValue,
    };
  },
};
</script>
