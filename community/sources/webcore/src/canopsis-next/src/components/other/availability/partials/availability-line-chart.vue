<template>
  <line-chart
    :chart-id="chartId"
    :options="chartOptions"
    :datasets="datasets"
    :width="width"
    :height="height"
    :dark="$system.dark"
    class="kpi-widget__chart"
  >
    <template #actions="{ chart }">
      <kpi-chart-export-actions
        :downloading="downloading"
        :chart="chart"
        class="mt-4"
        v-on="$listeners"
      />
    </template>
  </line-chart>
</template>

<script>
import {
  AVAILABILITY_LINE_CHART_X_AXES_IDS,
  SAMPLINGS,
  AVAILABILITY_LINE_CHART_Y_AXES_IDS,
  AVAILABILITY_SHOW_TYPE,
  AVAILABILITY_DISPLAY_PARAMETERS,
  MAX_TIME_VALUE_BY_SAMPLING,
  TIME_UNITS_BY_SAMPLING,
} from '@/constants';
import { COLORS } from '@/config';

import { getDateLabelBySampling } from '@/helpers/entities/metric/list';
import { convertNumberToFixedString } from '@/helpers/string';
import { convertDurationToString, fromSeconds } from '@/helpers/date/duration';

import KpiChartExportActions from '@/components/other/kpi/charts/partials/kpi-chart-export-actions.vue';

const LineChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/line-chart.vue');

export default {
  inject: ['$system'],
  components: { KpiChartExportActions, LineChart },
  props: {
    chartId: {
      type: String,
      required: false,
    },
    availabilities: {
      type: Array,
      default: () => [],
    },
    width: {
      default: 600,
      type: Number,
    },
    height: {
      default: 440,
      type: Number,
    },
    sampling: {
      type: String,
      default: SAMPLINGS.day,
    },
    showType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
    displayParameter: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
    responsive: {
      type: Boolean,
      default: true,
    },
    animation: {
      type: Boolean,
      default: true,
    },
  },
  computed: {
    isUptimeParameter() {
      return this.displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime;
    },

    isPercentType() {
      return this.showType === AVAILABILITY_SHOW_TYPE.percent;
    },

    maxTimeValue() {
      return MAX_TIME_VALUE_BY_SAMPLING[this.sampling];
    },

    samplingUnit() {
      return TIME_UNITS_BY_SAMPLING[this.sampling];
    },

    xAxes() {
      return {
        [AVAILABILITY_LINE_CHART_X_AXES_IDS.default]: {
          type: 'time',
          ticks: {
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
        [AVAILABILITY_LINE_CHART_Y_AXES_IDS.percent]: {
          stacked: this.stacked,
          display: 'auto',
          beginAtZero: true,
          max: 100,
          ticks: {
            callback: this.getChartYPercentTick,
            font: {
              size: 11,
              family: 'Arial, sans-serif',
            },
          },
        },
        [AVAILABILITY_LINE_CHART_Y_AXES_IDS.time]: {
          stacked: this.stacked,
          display: 'auto',
          beginAtZero: true,
          max: this.maxTimeValue,
          ticks: {
            callback: this.getChartYTimeTick,
            font: {
              size: 11,
              family: 'Arial, sans-serif',
            },
          },
        },
      };
    },

    datasets() {
      const color = this.isUptimeParameter ? COLORS.kpi.uptime : COLORS.kpi.downtime;

      return [
        {
          borderColor: color,
          xAxisID: AVAILABILITY_LINE_CHART_X_AXES_IDS.default,
          yAxisID: this.isPercentType
            ? AVAILABILITY_LINE_CHART_Y_AXES_IDS.percent
            : AVAILABILITY_LINE_CHART_Y_AXES_IDS.time,
          pointRadius: 0,
          data: this.availabilities.map(availability => ({
            x: availability.timestamp * 1000,
            y: this.convertValueByShowType(availability),
          })),
        },
      ];
    },

    chartOptions() {
      return {
        responsive: this.responsive,
        animation: this.animation,
        interaction: {
          mode: 'index',
          intersect: false,
        },
        scales: {
          ...this.xAxes,
          ...this.yAxes,
        },
        plugins: {
          legend: {
            display: false,
          },
          tooltip: {
            displayColors: false,
            callbacks: {
              title: this.getChartTooltipTitle,
              label: this.getChartTooltipLabel,
            },
          },
        },
      };
    },
  },
  methods: {
    convertValueByShowType(availability) {
      const { uptime, downtime } = availability;

      const value = this.isUptimeParameter ? uptime : downtime;

      if (this.isPercentType) {
        return (value / (uptime + downtime)) * 100;
      }

      return fromSeconds(value, this.samplingUnit);
    },

    getChartTimeTickLabel(_, index, data) {
      const { value } = data[index] ?? {};

      return getDateLabelBySampling(value, this.sampling);
    },

    getChartYPercentTick(value) {
      return `${value}%`;
    },

    getChartYTimeTick(value) {
      return `${value}${this.samplingUnit}`;
    },

    getChartTooltipTitle(data) {
      const [dataset] = data;
      const { x: timestamp } = dataset.raw;

      return getDateLabelBySampling(timestamp, this.sampling);
    },

    getChartTooltipLabel({ raw }) {
      if (this.isPercentType) {
        return `${convertNumberToFixedString(raw.y, 2)}%`;
      }

      return convertDurationToString(raw.y, undefined, this.samplingUnit);
    },
  },
};
</script>
