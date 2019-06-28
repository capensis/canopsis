<script>
import { merge } from 'lodash';
import { Bar } from 'vue-chartjs';

import chartAnnotationMixin from '@/mixins/chart/annotation';

export default {
  extends: Bar,
  mixins: [chartAnnotationMixin],
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
  mounted() {
    this.renderChart(this.chartData, this.mergedOptions);
  },
};
</script>
