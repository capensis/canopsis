<script>
import { HorizontalBar } from 'vue-chartjs';

export default {
  extends: HorizontalBar,
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
  },
  watch: {
    options(options, oldOptions) {
      if (options !== oldOptions) {
        this.renderChart(this.chartData, this.options);
      }
    },

    chartData(data, oldData) {
      if (data !== oldData) {
        this.renderChart(data, this.options);
      }
    },
  },
  mounted() {
    this.renderChart(this.chartData, this.options);
  },
};
</script>
