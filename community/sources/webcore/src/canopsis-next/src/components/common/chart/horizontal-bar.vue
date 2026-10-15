<template>
  <div :style="wrapperStyle" class="horizontal-bar">
    <div class="horizontal-bar__canvas">
      <canvas ref="canvas" />
    </div>
    <slot :chart="chart" name="actions" />
  </div>
</template>

<script>
import {
  BarController,
  BarElement,
  CategoryScale,
  Chart,
  Legend,
  LinearScale,
  Tooltip,
} from 'chart.js';
import { merge } from 'lodash';
import { markRaw } from 'vue';

Chart.register(
  BarController,
  BarElement,
  CategoryScale,
  Legend,
  LinearScale,
  Tooltip,
);

export default {
  props: {
    labels: {
      type: Array,
      default: () => [],
    },
    datasets: {
      type: Array,
      default: () => [],
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    width: {
      type: Number,
      required: false,
    },
    height: {
      type: Number,
      required: false,
    },
  },

  data() {
    return {
      chart: null,
    };
  },

  computed: {
    wrapperStyle() {
      return {
        width: this.width ? `${this.width}px` : '100%',
      };
    },

    canvasStyle() {
      return {
        height: this.height ? `${this.height}px` : '360px',
      };
    },

    chartData() {
      return {
        labels: this.labels,
        datasets: this.datasets,
      };
    },

    mergedOptions() {
      return merge(
        {
          indexAxis: 'y',
          responsive: false,
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
          type: 'bar',
          data: this.chartData,
          options: this.mergedOptions,
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
.horizontal-bar {
  position: relative;

  &__canvas {
    position: relative;
    height: v-bind('canvasStyle.height');
  }
}
</style>
