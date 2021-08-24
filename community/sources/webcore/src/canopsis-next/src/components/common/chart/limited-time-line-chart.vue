<script>
import { merge } from 'lodash';
import { Line } from '@/externals/vue-chart/components';

import { chartAnnotationMixin } from '@/mixins/chart/annotation';
import { chartLimitedSegmentMixin } from '@/mixins/chart/limited-segment';
import { chartZoomMixin } from '@/mixins/chart/zoom';
import { DATETIME_FORMATS } from '@/constants';

export default {
  extends: Line,
  mixins: [chartAnnotationMixin, chartLimitedSegmentMixin, chartZoomMixin],
  props: {
    ...Line.props,

    labels: {
      type: Array,
      required: false,
    },

    unit: {
      type: String,
      default: '',
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
              autoSkipPadding: 10,
              maxRotation: 0,
              beginAtZero: true,
            },
          },
          y: {
            type: 'linear',
            ticks: {
              min: 0,
              fontSize: 11,
              beginAtZero: true,
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
            callbacks: {
              label: tooltip => (this.unit ? `${tooltip.formattedValue} ${this.unit}` : tooltip.formattedValue),
            },
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
                  font: {
                    size: 14,
                  },
                },
              },
            },
          },
        },
      }, this.options);
    },

  },
  watch: {
    chartData(value, oldValue) {
      if (value !== oldValue) {
        this.updateChart(value, this.mergedOptions);

        this.chart.resetZoom();
      }
    },
  },
  mounted() {
    this.renderChart(this.chartData, this.mergedOptions);
  },
};
</script>
