<template lang="pug">
  bar-chart(:datasets="datasets", :options="sliChartOptions")
</template>

<script>
import { COLORS } from '@/config';
import {
  DATETIME_FORMATS,
  KPI_SLI_GRAPH_BAR_PERCENTAGE,
  SAMPLINGS,
  TIME_UNITS,
  KPI_SLI_GRAPH_DATA_TYPE,
} from '@/constants';

import { colorToRgba } from '@/helpers/color';

import { addUnitToDate, convertDateToString } from '@/helpers/date/date';
import { convertDurationToString, fromSeconds } from '@/helpers/date/duration';

const BarChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/bar-chart.vue');

export default {
  components: { BarChart },
  props: {
    metrics: {
      type: Array,
      default: () => [],
    },
    dataType: {
      type: String,
      default: KPI_SLI_GRAPH_DATA_TYPE.percent,
    },
    sampling: {
      type: String,
      default: SAMPLINGS.day,
    },
    responsive: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    maxValueBySampling() {
      return {
        [SAMPLINGS.hour]: 60,
        [SAMPLINGS.day]: 24,
        [SAMPLINGS.week]: 7,
        [SAMPLINGS.month]: 31,
      }[this.sampling];
    },

    maxValueByType() {
      if (this.dataType === KPI_SLI_GRAPH_DATA_TYPE.percent) {
        return 100;
      }

      return this.maxValueBySampling;
    },

    samplingUnit() {
      return {
        [SAMPLINGS.hour]: TIME_UNITS.minute,
        [SAMPLINGS.day]: TIME_UNITS.hour,
        [SAMPLINGS.week]: TIME_UNITS.day,
        [SAMPLINGS.month]: TIME_UNITS.day,
      }[this.sampling];
    },

    unit() {
      if (this.dataType === KPI_SLI_GRAPH_DATA_TYPE.percent) {
        return '%';
      }

      return this.samplingUnit;
    },

    datasets() {
      const { downtime, maintenance, uptime } = this.metrics.reduce((acc, metric) => {
        const x = metric.timestamp * 1000;

        acc.downtime.push({ x, y: this.convertValueBySamplingUnit(metric.downtime) });
        acc.maintenance.push({ x, y: this.convertValueBySamplingUnit(metric.maintenance) });
        acc.uptime.push({ x, y: this.convertValueBySamplingUnit(metric.uptime) });

        return acc;
      }, {
        downtime: [],
        maintenance: [],
        uptime: [],
      });

      return [{
        backgroundColor: colorToRgba(COLORS.kpi.uptime),
        barPercentage: KPI_SLI_GRAPH_BAR_PERCENTAGE,
        label: this.$t('common.uptime'),
        order: 1,
        data: uptime,
      }, {
        backgroundColor: colorToRgba(COLORS.kpi.downtime),
        barPercentage: KPI_SLI_GRAPH_BAR_PERCENTAGE,
        label: this.$t('common.downtime'),
        order: 2,
        data: downtime,
      }, {
        backgroundColor: colorToRgba(COLORS.kpi.maintenance),
        barPercentage: KPI_SLI_GRAPH_BAR_PERCENTAGE,
        label: this.$t('common.maintenance'),
        order: 3,
        data: maintenance,
      }];
    },

    sliChartOptions() {
      return {
        animation: false,
        responsive: this.responsive,
        scales: {
          x: {
            type: 'time',
            stacked: true,
            ticks: {
              source: 'data',
              fontSize: 11,
              callback: this.getChartTimeTickLabel,
            },
          },
          y: {
            stacked: true,
            max: this.maxValueByType,
            beginAtZero: true,
            ticks: {
              fontSize: 11,
              callback: this.getChartValueTickLabel,
            },
          },
        },
        interaction: {
          intersect: false,
          mode: 'index',
        },
        plugins: {
          legend: {
            position: 'top',
            align: 'end',
            labels: {
              boxWidth: 20,
              boxHeight: 20,
            },
          },
          tooltip: {
            itemSort: (a, b) => b.dataset.order - a.dataset.order,
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
    getLabelBySampling(value, sampling) {
      switch (sampling) {
        case SAMPLINGS.hour:
          return convertDateToString(value, DATETIME_FORMATS.long);
        case SAMPLINGS.day:
          return convertDateToString(value, DATETIME_FORMATS.short);
        case SAMPLINGS.week:
          return [
            convertDateToString(value, DATETIME_FORMATS.short),
            convertDateToString(addUnitToDate(value, 1, TIME_UNITS.week), DATETIME_FORMATS.short),
          ].join(' - \n');
      }

      return convertDateToString(value, DATETIME_FORMATS.yearWithMonth);
    },

    getChartValueTickLabel(value) {
      return `${value}${this.unit}`;
    },

    getChartTimeTickLabel(_, index, data) {
      const { value } = data[index] ?? {};

      return this.getLabelBySampling(value, this.sampling).split('\n');
    },

    getChartTooltipTitle(data) {
      const [dataset] = data;
      const { x: timestamp } = dataset.raw;

      return this.getLabelBySampling(timestamp, this.sampling);
    },

    getChartTooltipLabel(tooltip) {
      const { raw, dataset } = tooltip;

      const label = dataset.label.toLowerCase();

      if (this.dataType === KPI_SLI_GRAPH_DATA_TYPE.percent) {
        return `${raw.y}${this.unit} ${label}`;
      }

      const duration = convertDurationToString(raw.y, DATETIME_FORMATS.refreshFieldFormat, this.samplingUnit);

      return `${duration} ${label}`;
    },

    convertValueBySamplingUnit(value) {
      if (this.dataType === KPI_SLI_GRAPH_DATA_TYPE.percent) {
        return value;
      }

      return fromSeconds(value, this.samplingUnit);
    },
  },
};
</script>
