<template lang="pug">
  v-list(subheader)
    v-subheader {{ $t('kpi.collectionSetting.basicMetrics') }}
    v-list-tile(v-for="basicMetric in basicMetrics", :key="basicMetric")
      v-list-tile-content {{ basicMetric }}
</template>

<script>
import { ALARM_BASIC_METRIC_PARAMETERS } from '@/constants';

export default {
  computed: {
    basicMetrics() {
      return Object.values(ALARM_BASIC_METRIC_PARAMETERS).map(value => this.$t(`alarmList.metrics.${value}`));
    },
  },
};
</script>
