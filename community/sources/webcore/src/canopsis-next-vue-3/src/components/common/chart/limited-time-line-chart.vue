<template>
  <div class="limited-time-line-chart">
    <div class="limited-time-line-chart__canvas">
      <canvas ref="canvas" />
    </div>
    <slot :chart="chart" name="actions" />
  </div>
</template>

<script>
import { merge } from 'lodash';
import { markRaw } from 'vue';
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
import 'chartjs-adapter-moment';
import annotationPlugin from 'chartjs-plugin-annotation';
import zoomPlugin from 'chartjs-plugin-zoom';

import { DATETIME_FORMATS } from '@/constants';

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
  annotationPlugin,
  zoomPlugin,
);

const defaultDatasetOptions = {
  backgroundColor: 'green',
  pointBackgroundColor: 'transparent',
  pointBorderColor: 'grey',
  pointBorderWidth: 2,
  radius: 0,
  hoverRadius: 5,
  tension: 0.2,
};

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
        datasets: this.datasets.map(dataset => ({
          ...defaultDatasetOptions,
          ...dataset,
        })),
      };
    },

    mergedOptions() {
      const limit = this.options.plugins?.limit ?? {};

      return merge(
        {
          responsive: true,
          maintainAspectRatio: false,
          hover: {
            mode: 'index',
            intersect: false,
          },
          scales: {
            x: {
              type: 'time',
              beginAtZero: true,
              time: {
                tooltipFormat: DATETIME_FORMATS.longWithDayOfWeek,
                displayFormats: {
                  year: DATETIME_FORMATS.shortWithDayOfWeek,
                  quarter: DATETIME_FORMATS.shortWithDayOfWeek,
                  month: DATETIME_FORMATS.shortWithDayOfWeek,
                  week: DATETIME_FORMATS.shortWithDayOfWeek,
                  day: DATETIME_FORMATS.shortWithDayOfWeek,
                  hour: DATETIME_FORMATS.medium,
                  minute: DATETIME_FORMATS.medium,
                  second: DATETIME_FORMATS.medium,
                },
              },
              ticks: {
                max: Date.now(),
                font: {
                  size: 11,
                },
                autoSkip: true,
                autoSkipPadding: 5,
                maxRotation: 0,
              },
            },
            y: {
              type: 'linear',
              beginAtZero: true,
              ticks: {
                min: 0,
                font: {
                  size: 11,
                },
                callback: value => (value >= 1000 ? `${value / 1000} k` : value),
              },
            },
          },
          plugins: {
            legend: {
              display: false,
            },
            tooltip: {
              mode: 'index',
              intersect: false,
              displayColors: false,
            },
            annotation: {
              annotations: {
                limitLine: {
                  drawTime: 'afterDatasetsDraw',
                  display: !!limit.value,
                  type: 'line',
                  scaleID: 'y',
                  value: limit.value,
                  borderColor: limit.borderColor || limit.backgroundColor,
                  borderWidth: limit.borderWidth || 2,
                  borderDash: limit.borderDash || [],
                  label: {
                    display: false,
                  },
                },
              },
            },
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
.limited-time-line-chart {
  position: relative;
  width: 100%;

  &__canvas {
    position: relative;
    height: 360px;
  }
}
</style>
