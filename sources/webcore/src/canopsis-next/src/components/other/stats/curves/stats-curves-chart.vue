<script>
import { merge } from 'lodash';
import { Line } from 'vue-chartjs';

import ChartAnnotationPlugin from 'chartjs-plugin-annotation';

export default {
  extends: Line,
  props: {
    ...Line.props,

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
      return merge({
        responsive: true,
        maintainAspectRatio: false,
      }, this.options);
    },

  },
  watch: {
    chartData(value, oldValue) {
      if (value !== oldValue) {
        this.renderChart(value, this.mergedOptions);
      }
    },
  },
  created() {
    this.addPlugin(ChartAnnotationPlugin);
  },
  mounted() {
    this.renderChart(this.chartData, this.mergedOptions);
  },
};
</script>
