<script>
import { Line } from 'vue-chartjs';

export default {
  extends: Line,
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
