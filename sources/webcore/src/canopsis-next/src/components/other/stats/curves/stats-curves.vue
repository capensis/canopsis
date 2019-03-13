<script>
import { Line } from 'vue-chartjs';

export default {
  extends: Line,
  props: {
    ...Line.props,
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
    options() {
      return {
        responsive: true,
        maintainAspectRatio: false,
      };
    },

  },
  watch: {
    chartData(value, oldValue) {
      if (value !== oldValue) {
        this.renderChart(value, this.options);
      }
    },
  },
  mounted() {
    this.renderChart(this.chartData, this.options);
  },
};
</script>
