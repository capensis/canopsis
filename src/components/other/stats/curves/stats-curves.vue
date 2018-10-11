<script>
import { Line } from 'vue-chartjs';

export default {
  extends: Line,
  props: {
    labels: {
      type: Array,
    },
    datasets: {
      type: Array,
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
    chartData: {
      handler(value, oldValue) {
        if (value !== oldValue) {
          this.renderChart(value);
        }
      },
    },
  },
  mounted() {
    this.renderChart(this.chartData);
  },
};
</script>
