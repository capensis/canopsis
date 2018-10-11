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
    const data = {
      datasets: [{
        data: [15, 20, 26, 33, 34, 35, 26],
      }],
      labels: ['January', 'February', 'March', 'April', 'May', 'June'],
    };
    this.renderChart(data);
  },
};
</script>
