<template>
  <v-layout column>
    <entity-chart-widget
      v-for="(chart, index) in filteredCharts"
      :key="index"
      :widget="chart"
      :entity="entity"
      :available-metrics="availableMetrics"
    />
  </v-layout>
</template>

<script>
import EntityChartWidget from '@/components/widgets/chart/entity-chart-widget.vue';

export default {
  components: { EntityChartWidget },
  props: {
    charts: {
      type: Array,
      required: true,
    },
    entity: {
      type: Object,
      required: true,
    },
    availableMetrics: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    filteredCharts() {
      return this.charts.filter(
        ({ parameters }) => parameters.metrics.some(({ metric }) => this.availableMetrics.includes(metric)),
      );
    },
  },
};
</script>
