<template>
  <c-select-field
    v-field="value"
    :items="intervals"
    :label="label || $t('common.period')"
    item-text="text"
    item-value="value"
  />
</template>

<script>
import { computed } from 'vue';

import { TIME_UNITS } from '@/constants';

import { getNowIntervalValueForHours } from '@/helpers/date/date-intervals';

import { useI18n } from '@/hooks/i18n';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
  },
  setup() {
    const { t } = useI18n();

    const intervals = computed(() => {
      const result = [];
      const nowHours = (new Date()).getHours();

      for (let i = 0; i < 24; i += 1) {
        const calculatedHours = nowHours - i;
        const hours = calculatedHours >= 0 ? calculatedHours : 24 + calculatedHours;
        const slotValue = getNowIntervalValueForHours(i, TIME_UNITS.hour);
        const isYesterday = nowHours < i;

        result.push({
          value: slotValue,
          text: `${hours}:00 - ${hours}:59${isYesterday ? ` (${t('common.yesterday').toLowerCase()})` : ''}`,
        });
      }

      return result;
    });

    return {
      intervals,
    };
  },
};
</script>
