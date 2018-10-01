<script>
import { Line } from 'vue-chartjs';

export default {
  extends: Line,
  mounted() {
    this.renderChart({
      labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'],
      datasets: [
        {
          label: 'GitHub Commits',
          backgroundColor: '#f87979',
          borderColor: '#f87979',
          fill: false,
          data: [40, 20, 12, 39, 10, 40, 39, 80, 40, 20, 12, 11],
        },
      ],
    }, {
      scales: {
        yAxes: [{
          ticks: {
            callback(value) {
              return `$${value}`;
            },
          },
        }],
      },
    });
  },
};
</script>
