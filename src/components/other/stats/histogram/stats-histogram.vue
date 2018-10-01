<script>
import { Bar } from 'vue-chartjs';

export default {
  extends: Bar,
  mounted() {
    this.renderChart();
  },
};
</script>
