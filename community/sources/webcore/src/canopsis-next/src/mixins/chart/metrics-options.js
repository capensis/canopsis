import { X_AXES_IDS, Y_AXES_IDS, DATETIME_FORMATS } from '@/constants';

import { convertDurationToString } from '@/helpers/date/duration';
import { getDateLabelBySampling, getMaxTimeDurationForMetrics, isRatioMetric, isTimeMetric } from '@/helpers/metrics';

export const chartMetricsOptionsMixin = {
  computed: {
    maxTimeDuration() {
      return getMaxTimeDurationForMetrics(this.metrics);
    },

    labelsFont() {
      return {
        size: 11,
        family: 'Arial, sans-serif',
      };
    },

    xAxes() {
      return {
        [X_AXES_IDS.default]: {
          type: 'time',
          min: this.interval.from * 1000,
          max: this.interval.to * 1000,
          ticks: {
            min: this.minDate * 1000,
            max: Date.now(),
            source: 'data',
            callback: this.getChartTimeTickLabel,
            font: {
              size: 11,
              family: 'Arial, sans-serif',
            },
          },
        },
      };
    },

    yAxes() {
      return {
        [Y_AXES_IDS.default]: {
          stacked: this.stacked,
          beginAtZero: true,
          ticks: {
            font: this.labelsFont,
          },
        },
        [Y_AXES_IDS.percent]: {
          stacked: this.stacked,
          display: 'auto',
          position: 'right',
          beginAtZero: true,
          max: 100,
          ticks: {
            callback: this.getChartYPercentTick,
            font: this.labelsFont,
          },
        },
        [Y_AXES_IDS.time]: {
          stacked: this.stacked,
          display: 'auto',
          position: 'right',
          beginAtZero: true,
          ticks: {
            callback: this.getChartYTimeTick,
            font: this.labelsFont,
          },
        },
      };
    },
  },
  methods: {
    getChartTooltipLabel({ raw, dataset }) {
      const value = isTimeMetric(dataset.metric)
        ? convertDurationToString(
          raw.y,
          DATETIME_FORMATS.refreshFieldFormat,
          this.maxTimeDuration.unit,
        )
        : raw.y;

      return this.$t(`kpi.metrics.tooltip.${dataset.metric}`, { value });
    },

    getMetricYAxisId(metric) {
      if (isRatioMetric(metric)) {
        return Y_AXES_IDS.percent;
      }

      if (isTimeMetric(metric)) {
        return Y_AXES_IDS.time;
      }

      return Y_AXES_IDS.default;
    },

    getChartTooltipTitle(data) {
      const [dataset] = data;
      const { x: timestamp, originalX: originalTimestamp } = dataset.raw;

      return getDateLabelBySampling(originalTimestamp ?? timestamp, this.sampling);
    },

    getChartYPercentTick(value) {
      return `${value}%`;
    },

    getChartYTimeTick(value) {
      return `${value}${this.maxTimeDuration.unit}`;
    },

    getChartTimeTickLabel(_, index, data) {
      const { value } = data[index] ?? {};

      return getDateLabelBySampling(value, this.sampling);
    },
  },
};
