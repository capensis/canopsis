<template lang="pug">
  bar-chart(:datasets="datasets", :labels="labels", :options="ratingChartOptions")
</template>

<script>
import { COLORS } from '@/config';
import {
  DATETIME_FORMATS,
  KPI_SLI_GRAPH_BAR_PERCENTAGE,
  TIME_UNITS,
  USER_METRIC_PARAMETERS,
} from '@/constants';

import { colorToRgba } from '@/helpers/color';
import {
  convertDurationFormToMaxUnitDurationForm,
  convertDurationToString,
  fromSeconds,
} from '@/helpers/date/duration';
import { isRatioMetric, isTimeMetric } from '@/helpers/metrics';

import BarChart from '@/components/common/chart/bar-chart.vue';

export default {
  components: { BarChart },
  props: {
    metrics: {
      type: Array,
      default: () => [],
    },
    metric: {
      type: String,
      default: USER_METRIC_PARAMETERS.loginsNumber,
    },
    responsive: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    labels() {
      return this.metrics.map(({ label }) => label);
    },

    maxValue() {
      return Math.max.apply(null, this.metrics.map(({ value }) => value));
    },

    maxDuration() {
      return convertDurationFormToMaxUnitDurationForm({
        value: this.maxValue,
        unit: TIME_UNITS.second,
      });
    },

    datasets() {
      return [{
        backgroundColor: colorToRgba(COLORS.kpi.uptime),
        barPercentage: KPI_SLI_GRAPH_BAR_PERCENTAGE,
        data: this.metrics.map(({ value }) => this.convertValueByMetric(value)),
      }];
    },

    isTimeMetric() {
      return isTimeMetric(this.metric);
    },

    isRatioMetric() {
      return isRatioMetric(this.metric);
    },

    maxValueByMetric() {
      if (this.isRatioMetric) {
        return 100;
      }

      return undefined;
    },

    ratingChartOptions() {
      return {
        animation: false,
        responsive: this.responsive,
        scales: {
          x: {
            ticks: {
              font: {
                size: 11,
              },
            },
          },
          y: {
            beginAtZero: true,
            max: this.maxValueByMetric,
            ticks: {
              font: {
                size: 11,
              },
              callback: this.getChartYTick,
            },
          },
        },
        interaction: {
          intersect: false,
          mode: 'index',
        },
        plugins: {
          legend: {
            display: false,
          },
          tooltip: {
            displayColors: false,
            callbacks: {
              label: this.getChartTooltipLabel,
            },
          },
        },
      };
    },
  },
  methods: {
    convertValueByMetric(value) {
      if (!this.isTimeMetric) {
        return value;
      }

      return fromSeconds(value, this.maxDuration.unit);
    },

    getChartYTick(value) {
      if (this.isTimeMetric) {
        return `${value}${this.maxDuration.unit}`;
      }

      if (this.isRatioMetric) {
        return `${value}%`;
      }

      return value;
    },

    getChartTooltipLabel({ raw }) {
      const value = this.isTimeMetric
        ? convertDurationToString(raw, DATETIME_FORMATS.refreshFieldFormat, this.maxDuration.unit)
        : raw;

      return this.$t(`kpiMetrics.tooltip.${this.metric}`, { value });
    },
  },
};
</script>
