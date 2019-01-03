<script>
import { Bar } from 'vue-chartjs';

export default {
  extends: Bar,
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
    chartData(value, oldValue) {
      if (value !== oldValue) {
        this.renderChart(value);
      }
    },
  },
  mounted() {
    this.renderChart(this.chartData);
  },
};
</script>
