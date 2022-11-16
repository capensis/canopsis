import Chart from 'chart.js/auto';
import { cloneDeep, merge } from 'lodash';

import 'chartjs-adapter-moment';

const LIGHT_COLORS = {
  color: '#666',
  borderColor: 'rgba(0,0,0,0.1)',
};
const DARK_COLORS = {
  color: 'white',
  borderColor: 'rgba(255,255,255,0.1)',
};

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
    dark: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      chartPlugins: this.plugins,
      chartRendered: false,
    };
  },

  watch: {
    dark: {
      immediate: true,
      handler() {
        if (this.chart) {
          const { data } = this.chart.config;

          this.updateChart(data, this.previousOptions);
        }
      },
    },
  },

  created() {
    this.chart = null;
  },

  methods: {
    getOptions({ scales, ...options }) {
      const color = this.dark ? DARK_COLORS.color : LIGHT_COLORS.color;
      const borderColor = this.dark ? DARK_COLORS.borderColor : LIGHT_COLORS.borderColor;

      if (!scales) {
        return options;
      }

      return {
        scales: Object.entries(scales).reduce((acc, [key, scale]) => {
          acc[key] = merge({
            grid: { color: borderColor, borderColor, tickColor: borderColor },
            ticks: { color },
          }, scale);

          return acc;
        }, {}),
        ...options,
      };
    },

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
        this.chartRendered = false;
        this.chart.destroy();
      }

      if (!this.$refs.canvas) {
        return;
      }

      this.previousOptions = cloneDeep(options);
      this.chart = new Chart(this.$refs.canvas.getContext('2d'), {
        type: chartType,
        data,
        options: this.getOptions(options),
        plugins: this.chartPlugins,
      });

      this.chartRendered = true;
    },

    updateChart(data, options) {
      this.previousOptions = cloneDeep(options);
      this.chart.options = this.getOptions(options);
      this.chart.data = data;
      this.chart.stop();
      this.chart.update('none');
    },
  },

  beforeDestroy() {
    this.chart.destroy();
  },
});
