<script>
import { Bar } from 'vue-chartjs';

export default {
  extends: Bar,
  props: {
    widget: {
      type: Object,
    },
  },
  mounted() {
    const data = { ...this.widget.parameters.data };

    this.renderChart(data);
  },
};
</script>
