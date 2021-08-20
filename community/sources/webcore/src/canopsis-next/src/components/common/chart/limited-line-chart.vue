<script>
import { merge } from 'lodash';
import { Line } from 'vue-chartjs';

import { chartAnnotationMixin } from '@/mixins/chart/annotation';
import { chartLimitedSegmentMixin } from '@/mixins/chart/limited-segment';
import { chartZoomMixin } from '@/mixins/chart/zoom';

export default {
  extends: Line,
  mixins: [chartAnnotationMixin, chartLimitedSegmentMixin, chartZoomMixin],
  props: {
    ...Line.props,

    labels: {
      type: Array,
      default: () => [],
    },

    unit: {
      type: String,
      default: '',
    },

    limit: {
      /**
       * @type {Object} ChartLimitOptions
       * @property {number} value
       * @property {string} backgroundColor
       * @property {string} borderColor
       * @property {number} borderWidth
       * @property {Array} borderDash
       */
      type: Object,
      default: () => ({}),
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
          limit: this.limit,
          ...dataset,
        })),
      };
    },
    mergedOptions() {
      return merge({
        responsive: true,
        maintainAspectRatio: false,
        annotation: {
          annotations: this.limit.value
            ? [{
              type: 'line',
              mode: 'horizontal',
              scaleID: 'y-axis-0',
              value: this.limit.value,
              borderColor: this.limit.borderColor || this.limit.backgroundColor,
              borderWidth: this.limit.borderWidth || 2,
              borderDash: this.limit.borderDash,
            }]
            : [],
        },
        legend: {
          display: false,
        },
        hover: {
          mode: 'index',
          intersect: false,
        },
        tooltips: {
          mode: 'index',
          intersect: false,
          displayColors: false,
          callbacks: {
            label: tooltip => (this.unit ? `${tooltip.value} ${this.unit}` : tooltip.value),
          },
        },
        scales: {
          xAxes: [{
            ticks: {
              fontSize: 11,
            },
          }],
          yAxes: [{
            ticks: {
              fontSize: 11,
            },
          }],
        },
      }, this.options);
    },

  },
  watch: {
    chartData(value, oldValue) {
      if (value !== oldValue) {
        /* eslint-disable-next-line */
        // this.$data._chart.destroy();
        this.renderChart(value, this.mergedOptions);
        /* eslint-disable-next-line */
        // this.$data._chart.resetZoom();
      }
    },
  },
  mounted() {
    this.renderChart(this.chartData, this.mergedOptions);
  },
};
</script>
