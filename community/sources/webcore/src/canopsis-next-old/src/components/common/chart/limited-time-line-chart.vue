<script>
import { merge } from 'lodash';

import { DATETIME_FORMATS } from '@/constants';

import { chartAnnotationMixin } from '@/mixins/chart/annotation';
import { chartLimitedSegmentMixin } from '@/mixins/chart/limited-segment';
import { chartBackgroundMixin } from '@/mixins/chart/background';
import { chartZoomMixin } from '@/mixins/chart/zoom';

import { Line } from '@/externals/vue-chart/components';

export default {
  extends: Line,
  mixins: [chartAnnotationMixin, chartLimitedSegmentMixin, chartZoomMixin, chartBackgroundMixin],
  props: {
    ...Line.props,

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
  computed: {
    chartData() {
      return {
        labels: this.labels,
        datasets: this.datasets.map(dataset => ({
          backgroundColor: 'green',
          pointBackgroundColor: 'transparent',
          pointBorderColor: 'grey',
          pointBorderWidth: 2,
          radius: 0,
          hoverRadius: 5,
          tension: 0.2,
          ...dataset,
        })),
      };
    },
    mergedOptions() {
      const { limit = {} } = this.options.plugins;

      return merge({
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
              fontSize: 11,
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
              fontSize: 11,
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
                mode: 'horizontal',
                scaleID: 'y',
                value: limit.value,
                borderColor: limit.borderColor || limit.backgroundColor,
                borderWidth: limit.borderWidth || 2,
                borderDash: limit.borderDash || [],
                label: {
                  enabled: false,
                  font: {},
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
      }, this.options);
    },
  },
  watch: {
    chartData(data, oldData) {
      if (data !== oldData) {
        this.updateChart(data, this.mergedOptions);
      }
    },

    mergedOptions(options, oldOptions) {
      if (options !== oldOptions) {
        this.updateChart(this.chartData, options);
      }
    },
  },
  mounted() {
    this.renderChart(this.chartData, this.mergedOptions);
  },
};
</script>
