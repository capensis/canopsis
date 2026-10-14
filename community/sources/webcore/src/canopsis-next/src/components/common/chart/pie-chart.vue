<template>
  <div class="pie-chart">
    <div class="pie-chart__canvas">
      <canvas ref="canvas" />
    </div>
    <slot :chart="chart" name="actions" />
  </div>
</template>

<script>
import {
  ArcElement,
  Chart,
  Legend,
  PieController,
  Tooltip,
} from 'chart.js';
import ChartDataLabels from 'chartjs-plugin-datalabels';
import { merge } from 'lodash';
import { markRaw } from 'vue';

import { chartEmptyPiePlugin } from '@/helpers/charts/empty-pie';

Chart.register(
  ArcElement,
  Legend,
  PieController,
  Tooltip,
);

export default {
  props: {
    labels: {
      type: Array,
      required: false,
    },
    datasets: {
      type: Array,
      default: () => [],
    },
    options: {
      type: Object,
      default: () => ({}),
    },
  },

  data() {
    return {
      chart: null,
    };
  },

  computed: {
    chartData() {
      return {
        labels: this.labels,
        datasets: this.datasets,
      };
    },

    mergedOptions() {
      return merge(
        {
          responsive: true,
          maintainAspectRatio: false,
        },
        this.options,
      );
    },
  },

  watch: {
    chartData: {
      deep: true,
      handler() {
        this.updateChart();
      },
    },

    mergedOptions: {
      deep: true,
      handler() {
        this.updateChart();
      },
    },
  },

  mounted() {
    this.renderChart();
  },

  beforeUnmount() {
    this.chart?.destroy();
    this.chart = null;
  },

  methods: {
    renderChart() {
      this.chart = markRaw(
        new Chart(this.$refs.canvas, {
          type: 'pie',
          data: this.chartData,
          options: this.mergedOptions,
          plugins: [ChartDataLabels, chartEmptyPiePlugin],
        }),
      );
    },

    updateChart() {
      if (!this.chart) {
        return;
      }

      this.chart.data = this.chartData;
      this.chart.options = this.mergedOptions;
      this.chart.update();
    },
  },
};
</script>

<style lang="scss" scoped>
.pie-chart {
  position: relative;
  width: 100%;

  &__canvas {
    position: relative;
    height: 400px;
  }
}
</style>
