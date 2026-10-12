<template>
  <div :style="{ width: `${width}px`, height: `${height}px` }" class="test-suite-status-pie-chart">
    <pie-chart
      :labels="chartData.labels"
      :datasets="chartData.datasets"
      :options="mergedOptions"
    />
  </div>
</template>

<script>
import { merge } from 'lodash';

import { TEST_SUITE_COLORS, TEST_SUITE_STATUSES } from '@/constants';

import PieChart from '@/components/common/chart/pie-chart.vue';

export default {
  components: { PieChart },

  props: {
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
          datalabels: {
            display: false,
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
};
</script>

<style lang="scss" scoped>
.test-suite-status-pie-chart {
  :deep(.pie-chart__canvas) {
    height: 100%;
  }
}
</style>
