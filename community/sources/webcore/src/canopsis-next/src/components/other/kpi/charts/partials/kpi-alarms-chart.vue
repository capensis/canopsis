<template lang="pug">
  bar-chart(
    :datasets="datasets",
    :options="alarmsChartOptions",
    :width="width",
    :height="height",
    :dark="$system.dark"
  )
    template(#actions="{ chart }")
      kpi-chart-export-actions.mt-4(:downloading="downloading", :chart="chart", v-on="$listeners")
</template>

<script>
import { debounce } from 'lodash';
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
import { convertDateToStartOfUnitTimestamp, getNowTimestamp } from '@/helpers/date/date';

import BarChart from '@/components/common/chart/bar-chart.vue';

import KpiChartExportActions from './kpi-chart-export-actions.vue';

const Y_AXES_IDS = {
  default: 'y',
  percent: 'yPercent',
  time: 'yTime',
};

export default {
  inject: ['$system'],
  components: { KpiChartExportActions, BarChart },
  props: {
    metrics: {
      type: Array,
      default: () => [],
    },
    interval: {
      type: Object,
      default: () => ({
        from: 0,
        to: getNowTimestamp(),
      }),
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
      default: 440,
      type: Number,
    },
    minDate: {
      type: Number,
      required: false,
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
        label: this.$t(`alarm.metrics.${metric}`),
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
          [Y_AXES_IDS.default]: {
            beginAtZero: true,
            ticks: {
              font: {
                size: 11,
                family: 'Arial, sans-serif',
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
                family: 'Arial, sans-serif',
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
                family: 'Arial, sans-serif',
              },
            },
          },
        },
        interaction: {
          intersect: false,
          mode: 'x',
        },
        plugins: {
          legend: {
            position: 'right',
            align: 'start',
            maxWidth: 600,
            labels: {
              font: {
                size: 11,
                family: 'Arial, sans-serif',
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
          zoom: {
            limits: {
              x: {
                min: this.minDate * 1000,
                max: convertDateToStartOfUnitTimestamp(getNowTimestamp(), TIME_UNITS.hour) * 1000,
              },
            },
            pan: {
              enabled: true,
              mode: 'x',
              threshold: 3,
              onPanComplete: this.updateChartInterval,
            },
            zoom: {
              mode: 'x',
              wheel: {
                enabled: true,
                speed: 0.1,
              },
              pinch: {
                enabled: true,
              },
              drag: {
                enabled: true,
                modifierKey: 'ctrl',
              },
              onZoom: this.updateChartInterval,
            },
          },
        },
      };
    },
  },
  created() {
    this.debouncedUpdateInterval = debounce(this.updateInterval, 300);
  },
  methods: {
    updateChartInterval({ chart }) {
      const { min, max } = chart.scales.x;

      this.debouncedUpdateInterval({
        from: Math.floor(min / 1000),
        to: Math.ceil(max / 1000),
      });
    },

    updateInterval(interval) {
      this.$emit('zoom', interval);
    },

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
