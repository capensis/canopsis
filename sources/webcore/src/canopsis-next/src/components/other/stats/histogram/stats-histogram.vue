<script>
import { Bar } from 'vue-chartjs';

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
  },
  computed: {
    options() {
      return {
        responsive: true,
        maintainAspectRatio: false,
      };
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
        this.renderChart(value, this.options);
      }
    },
  },
  mounted() {
    this.renderChart(this.chartData, this.options);
  },
};
</script>
