<script>
import { merge } from 'lodash';
import { Pie } from 'vue-chartjs';

import { TEST_SUITE_COLORS, TEST_SUITE_STATUSES } from '@/constants';

import { chartAnnotationMixin } from '@/mixins/chart/annotation';

export default {
  extends: Pie,
  mixins: [chartAnnotationMixin],
  props: {
    ...Pie.props,

    statuses: {
      type: Object,
      required: true,
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    width: {
      type: Number,
      default: 160,
    },
    height: {
      type: Number,
      default: 160,
    },
  },
  computed: {
    mergedOptions() {
      return merge({
        responsive: true,
        maintainAspectRatio: false,
        legend: {
          position: 'right',
          labels: {
            boxWidth: 20,
          },
        },
      }, this.options);
    },

    chartData() {
      const { labels, dataset } = Object.entries(TEST_SUITE_STATUSES)
        .reduce((acc, [key, status]) => {
          acc.labels.push(this.$t(`testSuite.statuses.${status}`));
          acc.dataset.data.push(this.statuses[key]);
          acc.dataset.backgroundColor.push(TEST_SUITE_COLORS[status]);

          return acc;
        }, {
          labels: [],
          dataset: {
            backgroundColor: [],
            data: [],
          },
        });

      return {
        labels,
        datasets: [dataset],
      };
    },
  },
  watch: {
    chartData(value, oldValue) {
      if (value !== oldValue) {
        this.renderChart(value, this.mergedOptions);
      }
    },
  },
  mounted() {
    this.renderChart(this.chartData, this.mergedOptions);
  },
};
</script>
