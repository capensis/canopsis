<template>
  <line-chart
    :chart-id="chartId"
    :options="chartOptions"
    :datasets="datasets"
    :width="width"
    :height="height"
    :dark="$system.dark"
  >
    <template #actions="{ chart }">
      <chart-export-actions
        :downloading="downloading"
        :chart="chart"
        class="mt-4"
        v-on="$listeners"
      />
    </template>
  </line-chart>
</template>

<script>
import { computed } from 'vue';

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
import { getAvailabilityFieldByDisplayParameterAndShowType } from '@/helpers/entities/availability/entity';

import ChartExportActions from '@/components/common/chart/chart-export-actions.vue';

const LineChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/line-chart.vue');

export default {
  inject: ['$system'],
  components: { ChartExportActions, LineChart },
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
  setup(props) {
    const isUptimeParameter = computed(() => props.displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime);
    const isPercentType = computed(() => props.showType === AVAILABILITY_SHOW_TYPE.percent);

    const labelsFont = computed(() => ({
      size: 11,
      family: 'Arial, sans-serif',
    }));

    const getChartTimeTickLabel = (_, index, data) => {
      const { value } = data[index] ?? {};

      return getDateLabelBySampling(value, props.sampling);
    };
    const xAxes = computed(() => ({
      [AVAILABILITY_LINE_CHART_X_AXES_IDS.default]: {
        type: 'time',
        ticks: {
          source: 'data',
          callback: getChartTimeTickLabel,
          font: labelsFont.value,
        },
      },
    }));

    const maxTimeValue = computed(() => MAX_TIME_VALUE_BY_SAMPLING[props.sampling]);
    const samplingUnit = computed(() => TIME_UNITS_BY_SAMPLING[props.sampling]);

    const getChartYPercentTick = value => `${value}%`;
    const getChartYTimeTick = value => `${value}${samplingUnit.value}`;

    const yAxes = computed(() => ({
      [AVAILABILITY_LINE_CHART_Y_AXES_IDS.percent]: {
        display: 'auto',
        beginAtZero: true,
        max: 100,
        ticks: {
          callback: getChartYPercentTick,
          font: labelsFont.value,
        },
      },
      [AVAILABILITY_LINE_CHART_Y_AXES_IDS.time]: {
        display: 'auto',
        beginAtZero: true,
        max: maxTimeValue.value,
        ticks: {
          callback: getChartYTimeTick,
          font: labelsFont.value,
        },
      },
    }));

    const convertValueByShowType = (availability) => {
      const valueField = getAvailabilityFieldByDisplayParameterAndShowType(
        props.displayParameter,
        props.showType,
      );
      const value = availability[valueField];

      if (isPercentType.value) {
        return value;
      }

      return fromSeconds(value, props.samplingUnit);
    };
    const datasets = computed(() => {
      const color = isUptimeParameter.value ? COLORS.kpi.uptime : COLORS.kpi.downtime;

      return [
        {
          borderColor: color,
          xAxisID: AVAILABILITY_LINE_CHART_X_AXES_IDS.default,
          yAxisID: isPercentType.value
            ? AVAILABILITY_LINE_CHART_Y_AXES_IDS.percent
            : AVAILABILITY_LINE_CHART_Y_AXES_IDS.time,
          pointRadius: 0,
          data: props.availabilities.map(availability => ({
            x: availability.timestamp * 1000,
            y: convertValueByShowType(availability),
          })),
        },
      ];
    });

    const getChartTooltipTitle = (data) => {
      const [dataset] = data;
      const { x: timestamp } = dataset.raw;

      return getDateLabelBySampling(timestamp, props.sampling);
    };
    const getChartTooltipLabel = ({ raw }) => {
      if (isPercentType.value) {
        return `${convertNumberToFixedString(raw.y, 2)}%`;
      }

      return convertDurationToString(raw.y, undefined, samplingUnit.value);
    };
    const chartOptions = computed(() => ({
      responsive: props.responsive,
      animation: props.animation,
      interaction: {
        mode: 'index',
        intersect: false,
      },
      scales: {
        ...xAxes.value,
        ...yAxes.value,
      },
      plugins: {
        legend: {
          display: false,
        },
        tooltip: {
          displayColors: false,
          callbacks: {
            title: getChartTooltipTitle,
            label: getChartTooltipLabel,
          },
        },
      },
    }));

    return {
      chartOptions,
      datasets,
    };
  },
};
</script>
