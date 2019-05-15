<script>
import { Bar } from 'vue-chartjs';

import ChartAnnotationPlugin from 'chartjs-plugin-annotation';

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
        tooltips: {
          mode: 'index',
          intersect: false,
        },
        scales: {
          xAxes: [{
            stacked: true,
          }],
          yAxes: [{
            stacked: true,
          }],
        },
        annotation: {
          annotations: [{
            type: 'line',
            mode: 'horizontal',
            scaleID: 'y-axis-0',
            value: 1,
            borderColor: 'rgb(75, 192, 192)',
            borderWidth: 4,
            label: {
              enabled: true,
              content: 'Test label',
            },
          }],
        },
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
  created() {
    this.addPlugin(ChartAnnotationPlugin);
  },
  mounted() {
    this.renderChart(this.chartData, this.options);
  },
};
</script>
