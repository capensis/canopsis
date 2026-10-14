<template>
  <div class="line-chart">
    <div class="line-chart__canvas">
      <canvas ref="canvas" />
    </div>
    <slot :chart="chart" name="actions" />
  </div>
</template>

<script>
import {
  CategoryScale,
  Chart,
  Filler,
  Legend,
  LinearScale,
  LineController,
  LineElement,
  PointElement,
  TimeScale,
  Tooltip,
} from 'chart.js';
import zoomPlugin from 'chartjs-plugin-zoom';
import 'chartjs-adapter-moment';
import { merge } from 'lodash';
import { markRaw } from 'vue';

Chart.register(
  CategoryScale,
  Filler,
  Legend,
  LinearScale,
  LineController,
  LineElement,
  PointElement,
  TimeScale,
  Tooltip,
  zoomPlugin,
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
    dark: {
      type: Boolean,
      default: false,
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
          plugins: {
            zoom: {
              pan: {
                enabled: false,
              },
              zoom: {
                wheel: {
                  enabled: false,
                },
                pinch: {
                  enabled: false,
                },
                drag: {
                  enabled: false,
                },
              },
            },
          },
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
          type: 'line',
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
.line-chart {
  position: relative;
  width: 100%;

  &__canvas {
    position: relative;
    height: 360px;
  }
}
</style>
