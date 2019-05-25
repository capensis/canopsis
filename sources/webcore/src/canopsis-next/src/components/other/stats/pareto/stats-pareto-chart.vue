<script>
import { merge } from 'lodash';
import { Bar } from 'vue-chartjs';

import ChartAnnotationPlugin from 'chartjs-plugin-annotation';

export default {
  extends: Bar,
  props: {
    ...Bar.props,

    labels: {
      type: Array,
    },
    datasets: {
      type: Array,
    },
    options: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    mergedOptions() {
      return merge({
        responsive: true,
        maintainAspectRatio: false,
        tooltips: {
          mode: 'index',
          intersect: false,
        },
      }, this.options);
    },

    chartData() {
      return {
        labels: this.labels,
        datasets: this.datasets,
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
  created() {
    this.addPlugin(ChartAnnotationPlugin);
  },
  mounted() {
    this.renderChart(this.chartData, this.mergedOptions);
  },
};
</script>
