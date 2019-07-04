<script>
import { merge } from 'lodash';
import { Line } from 'vue-chartjs';

import chartAnnotationMixin from '@/mixins/chart/annotation';

export default {
  extends: Line,
  mixins: [chartAnnotationMixin],
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
  mounted() {
    this.renderChart(this.chartData, this.mergedOptions);
  },
};
</script>
