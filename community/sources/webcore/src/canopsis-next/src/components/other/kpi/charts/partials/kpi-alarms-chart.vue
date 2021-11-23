<template lang="pug">
  bar-chart(:datasets="datasets", :options="alarmsChartOptions", :width="width", :height="height")
    template(#actions="{ chart }")
      kpi-chart-export-actions.mt-4(:downloading="downloading", :chart="chart", v-on="$listeners")
</template>

<script>
import {
  DATETIME_FORMATS,
  KPI_ALARMS_GRAPH_BAR_PERCENTAGE,
  SAMPLINGS,
  TIME_UNITS,
} from '@/constants';

import {
  convertDurationToMaxUnitDuration,
  convertDurationToString,
  fromSeconds,
} from '@/helpers/date/duration';
import { getDateLabelBySampling, isRatioMetric, isTimeMetric } from '@/helpers/metrics';
import { getMetricColor } from '@/helpers/color';

import BarChart from '@/components/common/chart/bar-chart.vue';

import KpiChartExportActions from './kpi-chart-export-actions.vue';

const Y_AXES_IDS = {
  default: 'y',
  percent: 'yPercent',
  time: 'yTime',
};

export default {
  components: { KpiChartExportActions, BarChart },
  props: {
    metrics: {
      type: Array,
      default: () => [],
    },
    sampling: {
      type: String,
      default: SAMPLINGS.day,
    },
    responsive: {
      type: Boolean,
      default: false,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
    width: {
      default: 600,
      type: Number,
    },
    height: {
      default: 400,
      type: Number,
    },
  },
  computed: {
    maxTimeValue() {
      return Math.max.apply(null, this.metrics.reduce((acc, { title: metric, data }) => {
        if (isTimeMetric(metric)) {
          const maxDatasetValue = Math.max.apply(null, data.map(({ value }) => value));

          acc.push(maxDatasetValue);
        }

        return acc;
      }, []));
    },

    maxTimeDuration() {
      return convertDurationToMaxUnitDuration({
        value: this.maxTimeValue,
        unit: TIME_UNITS.second,
      });
    },

    datasets() {
      return this.metrics.map(({ title: metric, data }) => ({
        metric,
        backgroundColor: getMetricColor(metric),
        barPercentage: KPI_ALARMS_GRAPH_BAR_PERCENTAGE,
        yAxisID: this.getMetricYAxisId(metric),
        label: this.$t(`alarmList.metrics.${metric}`),
        data: data.map(({ timestamp, value }) => ({
          x: timestamp * 1000,
          y: this.convertValueByMetricType(value, metric),
        })),
      }));
    },

    alarmsChartOptions() {
      return {
        animation: false,
        responsive: this.responsive,
        scales: {
          x: {
            type: 'time',
            ticks: {
              source: 'data',
              callback: this.getChartTimeTickLabel,
              font: {
                size: 11,
              },
            },
          },
          [Y_AXES_IDS.default]: {
            beginAtZero: true,
            ticks: {
              font: {
                size: 11,
              },
            },
          },
          [Y_AXES_IDS.percent]: {
            display: 'auto',
            position: 'right',
            beginAtZero: true,
            max: 100,
            ticks: {
              callback: this.getChartYPercentTick,
              font: {
                size: 11,
              },
            },
          },
          [Y_AXES_IDS.time]: {
            display: 'auto',
            position: 'right',
            beginAtZero: true,
            ticks: {
              callback: this.getChartYTimeTick,
              font: {
                size: 11,
              },
            },
          },
        },
        interaction: {
          intersect: false,
          mode: 'x',
        },
        plugins: {
          background: {
            color: 'white',
          },
          legend: {
            position: 'right',
            align: 'start',
            maxWidth: 300,
            labels: {
              font: {
                size: 11,
              },
              boxWidth: 15,
              boxHeight: 15,
            },
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
    getChartTooltipLabel({ raw, dataset }) {
      const value = isTimeMetric(dataset.metric)
        ? convertDurationToString(raw.y, DATETIME_FORMATS.refreshFieldFormat, this.maxTimeDuration.unit)
        : raw.y;

      return this.$t(`kpiMetrics.tooltip.${dataset.metric}`, { value });
    },

    getChartTooltipTitle(data) {
      const [dataset] = data;
      const { x: timestamp } = dataset.raw;

      return getDateLabelBySampling(timestamp, this.sampling);
    },

    convertValueByMetricType(value, metric) {
      if (isTimeMetric(metric)) {
        return fromSeconds(value, this.maxTimeDuration.unit);
      }

      return value;
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
</script>
