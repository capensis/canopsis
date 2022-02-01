import Chart from 'chart.js/auto';
import 'chartjs-adapter-moment';

/**
 * @param {string} chartId
 * @param {string} chartType
 * @return {Object}
 */
export const generateChart = (chartId, chartType) => ({
  render(h) {
    return h('div', {}, [
      h('div', {
        style: this.styles,
        class: this.cssClasses,
      }, [h('canvas', {
        attrs: {
          id: this.chartId,
          width: this.width,
          height: this.height,
        },
        ref: 'canvas',
      })]),
      this.chartRendered && this.$scopedSlots.actions?.({ chart: this.chart }),
    ]);
  },

  props: {
    chartId: {
      default: chartId,
      type: String,
    },
    width: {
      default: 400,
      type: Number,
    },
    height: {
      default: 400,
      type: Number,
    },
    cssClasses: {
      type: String,
      default: '',
    },
    styles: {
      type: Object,
    },
    plugins: {
      type: Array,
      default: () => [],
    },
  },

  data() {
    return {
      chartPlugins: this.plugins,
      chartRendered: false,
    };
  },

  created() {
    this.chart = null;
  },

  methods: {
    addPlugin(plugin) {
      this.chartPlugins.push(plugin);
    },

    generateLegend() {
      if (!this.chart) {
        return;
      }

      this.chart.generateLegend();
    },

    renderChart(data, options) {
      if (this.chart) {
        this.chart.destroy();
      }

      if (!this.$refs.canvas) {
        return;
      }

      this.chart = new Chart(this.$refs.canvas.getContext('2d'), {
        type: chartType,
        data,
        options,
        plugins: this.chartPlugins,
      });

      this.chartRendered = true;
    },

    updateChart(data, options) {
      this.chart.options = options;
      this.chart.data = data;
      this.chart.stop();
      this.chart.update('none');
    },
  },

  beforeDestroy() {
    this.chart.destroy();
  },
});
