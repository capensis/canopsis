<script>
import { merge } from 'lodash';

import { chartBackgroundMixin } from '@/mixins/chart/background';

import { Bar } from '@/externals/vue-chart/components';

export default {
  extends: Bar,
  mixins: [chartBackgroundMixin],
  props: {
    ...Bar.props,

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
        datasets: this.datasets,
      };
    },
    mergedOptions() {
      return merge({ responsive: true, maintainAspectRatio: false }, this.options);
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
