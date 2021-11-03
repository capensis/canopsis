<template lang="pug">
  div(@touchend.stop)
    v-layout.ml-4.mb-4(align-center)
      c-quick-date-interval-field(
        :interval="pagination.interval",
        @input="updateInterval"
      )
    div
      bar-chart(:datasets="datasets", :options="sliChartOptions")
</template>

<script>
import { COLORS } from '@/config';
import {
  DATETIME_FORMATS,
  QUICK_RANGES,
  KPI_SLI_GRAPH_BAR_PERCENTAGE,
  SAMPLINGS,
  TIME_UNITS,
  KPI_SLI_GRAPH_DATA_TYPE,
} from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { colorToRgba } from '@/helpers/color';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';
import { convertNumberToFixedString } from '@/helpers/string';
import { addUnitToDate, convertDateToString, getDaysInMonth } from '@/helpers/date/date';
import { convertDurationToString, fromSeconds, toSeconds } from '@/helpers/date/duration';

const BarChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/bar-chart.vue');

export default {
  components: { BarChart },
  mixins: [entitiesMetricsMixin, localQueryMixin],
  data() {
    return {
      sliMetrics: [],
      query: {
        sampling: SAMPLINGS.day,
        type: KPI_SLI_GRAPH_DATA_TYPE.percent,
        interval: {
          from: QUICK_RANGES.last30Days.start,
          to: QUICK_RANGES.last30Days.stop,
        },
      },
    };
  },
  computed: {
    interval() {
      return {
        from: convertStartDateIntervalToTimestamp(this.pagination.interval.from),
        to: convertStopDateIntervalToTimestamp(this.pagination.interval.to),
      };
    },

    maxValueBySampling() {
      return {
        [SAMPLINGS.hour]: 60,
        [SAMPLINGS.day]: 24,
        [SAMPLINGS.week]: 7,
        [SAMPLINGS.month]: 31,
      }[this.pagination.sampling];
    },

    maxValueByType() {
      if (this.pagination.type === KPI_SLI_GRAPH_DATA_TYPE.percent) {
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
      }[this.pagination.sampling];
    },

    unit() {
      if (this.pagination.type === KPI_SLI_GRAPH_DATA_TYPE.percent) {
        return '%';
      }

      return this.samplingUnit;
    },

    datasets() {
      const { downtime, maintenance, uptime } = this.sliMetrics.reduce((acc, metric) => {
        const x = metric.timestamp * 1000;

        acc.downtime.push({ x, y: this.convertSecondsToValue(metric.downtime) });
        acc.maintenance.push({ x, y: this.convertSecondsToValue(metric.maintenance) });
        acc.uptime.push({ x, y: this.convertSecondsToValue(metric.uptime) });

        return acc;
      }, {
        downtime: [],
        maintenance: [],
        uptime: [],
      });

      return [{
        backgroundColor: colorToRgba(COLORS.kpi.uptime),
        barPercentage: KPI_SLI_GRAPH_BAR_PERCENTAGE,
        fill: true,
        label: this.$t('common.uptime'),
        order: 1,
        data: uptime,
      }, {
        backgroundColor: colorToRgba(COLORS.kpi.downtime),
        barPercentage: KPI_SLI_GRAPH_BAR_PERCENTAGE,
        fill: true,
        label: this.$t('common.downtime'),
        order: 2,
        data: downtime,
      }, {
        backgroundColor: colorToRgba(COLORS.kpi.maintenance),
        barPercentage: KPI_SLI_GRAPH_BAR_PERCENTAGE,
        fill: true,
        label: this.$t('common.maintenance'),
        order: 3,
        data: maintenance,
      }];
    },

    sliChartOptions() {
      return {
        animation: false,
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
  mounted() {
    this.fetchList();
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

      return this.getLabelBySampling(value, this.pagination.sampling).split('\n');
    },

    getChartTooltipTitle(data) {
      const [dataset] = data;
      const { x: timestamp } = dataset.raw;

      return this.getLabelBySampling(timestamp, this.pagination.sampling);
    },

    getChartTooltipLabel(tooltip) {
      const { raw, dataset } = tooltip;

      const label = dataset.label.toLowerCase();

      if (this.pagination.type === KPI_SLI_GRAPH_DATA_TYPE.percent) {
        return `${raw.y}${this.unit} ${label}`;
      }

      const duration = convertDurationToString(raw.y, DATETIME_FORMATS.refreshFieldFormat, this.samplingUnit);

      return `${duration} ${label}`;
    },

    convertSecondsToValue(value) {
      if (this.pagination.type === KPI_SLI_GRAPH_DATA_TYPE.percent) {
        return this.convertSecondsToPercent(value);
      }

      return fromSeconds(value, this.samplingUnit);
    },

    getMaxValueForPercent(value) {
      if (this.pagination.sampling === SAMPLINGS.month) {
        return toSeconds(getDaysInMonth(value), TIME_UNITS.day);
      }

      /**
       * We have 86400 seconds in day, 86400 / 100 = 864
       */
      return toSeconds(this.maxValueBySampling, this.samplingUnit);
    },

    convertSecondsToPercent(value) {
      const maxValue = this.getMaxValueForPercent(value);

      return convertNumberToFixedString((value / maxValue) * 100, 2);
    },

    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
    },

    async fetchList() {
      this.sliMetrics = await this.fetchSliMetricsWithoutStore({
        params: {
          from: this.interval.from,
          to: this.interval.to,
          sampling: this.pagination.sampling,
        },
      });
    },
  },
};
</script>
