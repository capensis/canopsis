<script>
import { merge } from 'lodash';

import { Pie } from '@/externals/vue-chart/components';

import { TEST_SUITE_COLORS, TEST_SUITE_STATUSES } from '@/constants';

export default {
  extends: Pie,
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
        plugins: {
          legend: {
            position: 'right',
            labels: {
              boxWidth: 20,
            },
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
