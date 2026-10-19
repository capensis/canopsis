<template>
  <v-list subheader>
    <v-subheader>{{ $t('kpi.collectionSetting.basicMetrics') }}</v-subheader>
    <v-list-item
      v-for="basicMetric in basicMetrics"
      :key="basicMetric"
    >
      <v-list-item-content>{{ basicMetric }}</v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script>
import { ALARM_BASIC_METRIC_PARAMETERS } from '@/constants';

export default {
  computed: {
    basicMetrics() {
      return Object.values(ALARM_BASIC_METRIC_PARAMETERS).map(value => this.$t(`alarm.metrics.${value}`));
    },
  },
};
</script>
