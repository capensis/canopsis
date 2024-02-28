<template>
  <v-layout class="gap-5" column>
    <v-progress-linear
      :value="uptimePercent"
      :color="uptimeColor"
      :background-color="downtimeColor"
      height="50"
    />
    <v-layout column>
      <availability-bar-chart-information-row :label="$t('common.uptime')" :color="uptimeColor">
        {{ uptimeValue }}
      </availability-bar-chart-information-row>
      <availability-bar-chart-information-row :label="$t('common.downtime')" :color="downtimeColor">
        {{ downtimeValue }}
      </availability-bar-chart-information-row>
      <template v-if="!isPercentType">
        <availability-bar-chart-information-row :label="$t('common.totalActiveTime')" class="text--secondary">
          {{ totalActiveTimeDuration }}
        </availability-bar-chart-information-row>
        <availability-bar-chart-information-row
          v-if="inactiveTime"
          :label="$t('common.inactiveTime')"
          class="text--disabled"
        >
          {{ inactiveTimeDuration }}
        </availability-bar-chart-information-row>
      </template>
    </v-layout>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';
import { COLORS } from '@/config';

import { convertDurationToString } from '@/helpers/date/duration';
import { convertNumberToFixedString } from '@/helpers/string';

import AvailabilityBarChartInformationRow from './availability-bar-chart-information-row.vue';

export default {
  components: { AvailabilityBarChartInformationRow },
  props: {
    uptime: {
      type: Number,
      required: true,
    },
    downtime: {
      type: Number,
      required: true,
    },
    inactiveTime: {
      type: Number,
      default: 0,
    },
    uptimeColor: {
      type: String,
      default: COLORS.kpi.uptime,
    },
    downtimeColor: {
      type: String,
      default: COLORS.kpi.downtime,
    },
    showType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.duration,
    },
  },
  setup(props) {
    const isPercentType = computed(() => props.showType === AVAILABILITY_SHOW_TYPE.percent);

    const totalTime = computed(() => props.uptime + props.downtime);

    const convertValueToPercent = value => (value / totalTime.value) * 100;
    const convertValueToPercentString = value => `${convertNumberToFixedString(convertValueToPercent(value), 2)}%`;

    const uptimePercent = computed(() => convertValueToPercentString(props.uptime));
    const downtimePercent = computed(() => convertValueToPercentString(props.downtime));

    const uptimeDuration = computed(() => convertDurationToString(props.uptime));
    const downtimeDuration = computed(() => convertDurationToString(props.downtime));
    const totalActiveTimeDuration = computed(() => convertDurationToString(totalTime.value));
    const inactiveTimeDuration = computed(() => convertDurationToString(props.inactiveTime));

    const uptimeValue = computed(() => (isPercentType.value ? uptimePercent.value : uptimeDuration.value));
    const downtimeValue = computed(() => (isPercentType.value ? downtimePercent.value : downtimeDuration.value));

    return {
      isPercentType,
      uptimePercent,
      downtimePercent,
      uptimeValue,
      downtimeValue,
      totalActiveTimeDuration,
      inactiveTimeDuration,
    };
  },
};
</script>
