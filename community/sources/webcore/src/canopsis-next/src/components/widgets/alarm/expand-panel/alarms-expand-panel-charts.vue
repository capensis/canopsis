<template lang="pug">
  v-layout(column)
    alarm-chart-widget(
      v-for="(chart, index) in charts",
      :key="index",
      :widget="chart",
      :alarm="alarm"
    )
</template>

<script>
import AlarmChartWidget from '@/components/widgets/chart/alarm-chart-widget.vue';

export default {
  components: { AlarmChartWidget },
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
