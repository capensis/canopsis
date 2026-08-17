<template>
  <c-layout gap="5" column>
    <v-progress-linear
      :model-value="uptimePercent"
      :color="uptimeColor"
      :bg-color="downtimeColor"
      height="50"
    />
    <c-layout column>
      <availability-bar-chart-information-row :label="$t('common.uptime')" :color="uptimeColor">
        {{ uptimeValue }}
      </availability-bar-chart-information-row>
      <availability-bar-chart-information-row :label="$t('common.downtime')" :color="downtimeColor">
        {{ downtimeValue }}
      </availability-bar-chart-information-row>
      <v-expand-transition>
        <div v-if="!isPercentType">
          <availability-bar-chart-information-row
            :label="$t('common.totalActiveTime')"
            class="text-medium-emphasis"
          >
            {{ totalActiveTimeDuration }}
          </availability-bar-chart-information-row>
          <availability-bar-chart-information-row
            v-if="availability.total_inactive_time"
            :label="$t('common.inactiveTime')"
            class="text-disabled"
          >
            {{ inactiveTimeDuration }}
          </availability-bar-chart-information-row>
        </div>
      </v-expand-transition>
    </c-layout>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';
import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import { convertDurationToString } from '@/helpers/date/duration';

import AvailabilityBarChartInformationRow from './availability-bar-chart-information-row.vue';

export default {
  components: { AvailabilityBarChartInformationRow },

  props: {
    availability: {
      type: Object,
      required: true,
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
    const uptimePercent = computed(() => Number(props.availability.uptime_share ?? 0));
    const downtimePercent = computed(() => `${props.availability.downtime_share}%`);
    const uptimeDuration = computed(() => convertDurationToString(props.availability.uptime_duration));
    const downtimeDuration = computed(() => convertDurationToString(props.availability.downtime_duration));
    const totalActiveTimeDuration = computed(() =>
      convertDurationToString(props.availability.total_active_time),
    );
    const inactiveTimeDuration = computed(() =>
      convertDurationToString(props.availability.total_inactive_time),
    );
    const uptimeValue = computed(() =>
      isPercentType.value ? `${uptimePercent.value}%` : uptimeDuration.value,
    );
    const downtimeValue = computed(() =>
      isPercentType.value ? downtimePercent.value : downtimeDuration.value,
    );

    return {
      downtimeValue,
      inactiveTimeDuration,
      isPercentType,
      totalActiveTimeDuration,
      uptimePercent,
      uptimeValue,
    };
  },
};
</script>
