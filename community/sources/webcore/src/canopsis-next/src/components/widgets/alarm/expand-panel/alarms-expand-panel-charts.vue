<template lang="pug">
  v-layout(column)
    alarm-metrics-chart(
      v-for="(chart, index) in charts",
      :key="index",
      :widget="chart",
      :alarm="alarm"
    )
</template>

<script>
import AlarmMetricsChart from '@/components/widgets/chart/partials/alarm-metrics-chart.vue';

export default {
  components: { AlarmMetricsChart },
  props: {
    charts: {
      type: Array,
      required: true,
    },
    alarm: {
      type: Object,
      required: true,
    },
  },
};
</script>
