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
  computed: {
    isPercentType() {
      return this.showType === AVAILABILITY_SHOW_TYPE.percent;
    },

    totalTime() {
      return this.uptime + this.downtime;
    },

    uptimePercent() {
      return this.convertValueToPercentString(this.uptime);
    },

    uptimeDuration() {
      return convertDurationToString(this.uptime);
    },

    downtimePercent() {
      return this.convertValueToPercentString(this.downtime);
    },

    downtimeDuration() {
      return convertDurationToString(this.downtime);
    },

    uptimeValue() {
      return this.isPercentType ? this.uptimePercent : this.uptimeDuration;
    },

    downtimeValue() {
      return this.isPercentType ? this.downtimePercent : this.downtimeDuration;
    },

    totalActiveTimeDuration() {
      return convertDurationToString(this.totalTime);
    },

    inactiveTimeDuration() {
      return convertDurationToString(this.inactiveTime);
    },
  },
  methods: {
    convertValueToPercent(value) {
      return (value / this.totalTime) * 100;
    },

    convertValueToPercentString(value) {
      return `${convertNumberToFixedString(this.convertValueToPercent(value), 2)}%`;
    },
  },
};
</script>
