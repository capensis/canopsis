<template>
  <span>{{ value }}</span>
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE } from '@/constants';

import { convertNumberToRoundedPercentString } from '@/helpers/string';
import { convertDurationToString } from '@/helpers/date/duration';

export default {
  props: {
    availability: {
      type: Object,
      required: true,
    },
    displayParameter: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
    showType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
  },
  setup(props) {
    const isUptimeParameter = computed(() => props.displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime);
    const isPercentType = computed(() => props.showType === AVAILABILITY_SHOW_TYPE.percent);

    const targetValue = computed(() => (
      isUptimeParameter.value ? props.availability.uptime : props.availability.downtime
    ));

    const totalTime = computed(() => props.availability.uptime + props.availability.downtime);
    const percent = computed(() => convertNumberToRoundedPercentString(targetValue.value / totalTime.value));
    const duration = computed(() => convertDurationToString(targetValue.value));
    const value = computed(() => (isPercentType.value ? percent.value : duration.value));

    return {
      value,
    };
  },
};
</script>
