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
    chartDatas() {
      return {
        labels: this.labels,
        datasets: this.datasets,
      };
    },
  },
  watch: {
    chartDatas(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.renderChart(newVal);
      }
    },
  },
  mounted() {
    this.renderChart(this.chartDatas);
  },
};
</script>
