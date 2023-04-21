import {
  X_AXES_IDS,
  Y_AXES_IDS,
  DATETIME_FORMATS,
  MAX_METRICS_DISPLAY_COUNT,
  TIME_UNITS,
} from '@/constants';

import { convertDurationToString } from '@/helpers/date/duration';
import {
  getDateLabelBySampling,
  getMaxTimeDurationForMetrics,
  hasHistoryData,
  isRatioMetric,
  isTimeMetric,
} from '@/helpers/metrics';
import {
  convertDateToEndOfUnitTimestamp,
  convertDateToStartOfDayTimestamp,
  convertDateToTimestampByTimezone,
} from '@/helpers/date/date';

export const chartMetricsOptionsMixin = {
  inject: ['$system'],
  computed: {
    hasHistoryData() {
      return hasHistoryData(this.metrics);
    },

    maxMetricsCount() {
      return this.hasHistoryData
        ? MAX_METRICS_DISPLAY_COUNT / 2
        : MAX_METRICS_DISPLAY_COUNT;
    },

    preparedMetrics() {
      return this.metrics?.slice(0, this.maxMetricsCount) ?? [];
    },

    realMetricsCount() {
      return this.hasHistoryData
        ? this.metrics?.length * 2
        : this.metrics?.length;
    },

    tooltipBodyFontSize() {
      if (this.realMetricsCount > 32) {
        return 9;
      }

      if (this.realMetricsCount > 24) {
        return 10;
      }

      return 11;
    },

    maxTimeDuration() {
      return getMaxTimeDurationForMetrics(this.metrics);
    },

    labelsFont() {
      return {
        size: 11,
        family: 'Arial, sans-serif',
      };
    },

    legend() {
      const legend = {
        position({ chart }) {
          return chart.width > 600 ? 'right' : 'top';
        },
        maxWidth: 700,
        labels: {
          font: this.labelsFont,
          boxWidth: 15,
          boxHeight: 15,
        },
      };

      if (this.realMetricsCount > MAX_METRICS_DISPLAY_COUNT) {
        legend.title = {
          text: [
            this.$t('kpi.largeCountOfMetrics'),
            this.$t('kpi.onlyDisplayed', { count: MAX_METRICS_DISPLAY_COUNT }),
          ],
          padding: 10,
          display: true,
          color: 'gray',
          font: {
            ...this.labelsFont,

            style: 'italic',
          },
        };
      }

      return legend;
    },

    xAxes() {
      return {
        [X_AXES_IDS.default]: {
          type: 'time',
          min: convertDateToTimestampByTimezone(
            convertDateToStartOfDayTimestamp(this.interval.from),
            this.$system.timezone,
          ) * 1000,
          max: convertDateToTimestampByTimezone(
            convertDateToEndOfUnitTimestamp(this.interval.to, TIME_UNITS.day),
            this.$system.timezone,
          ) * 1000,
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
