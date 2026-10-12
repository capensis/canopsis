<script>
import { merge } from 'lodash';

import { Bar } from '@/externals/vue-chart/components';

export default {
  extends: Bar,
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
  },
  computed: {
    chartData() {
      return {
        labels: this.labels,
        datasets: this.datasets,
      };
    },

    mergedOptions() {
      return merge({ indexAxis: 'y' }, this.options);
    },
  },
  watch: {
    mergedOptions(options, oldOptions) {
      if (options !== oldOptions) {
        this.updateChart(this.chartData, this.mergedOptions);
      }
    },

    chartData(data, oldData) {
      if (data !== oldData) {
        this.updateChart(data, this.mergedOptions);
      }
    },
  },
  mounted() {
    this.renderChart(this.chartData, this.mergedOptions);
  },
};
</script>
