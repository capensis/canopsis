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
import { DATETIME_FORMATS, QUICK_RANGES, KPI_SLI_GRAPH_BAR_PERCENTAGE, KPI_SAMPLINGS } from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { colorToRgba } from '@/helpers/color';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';
import { convertNumberToFixedString } from '@/helpers/string';

const BarChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/bar-chart.vue');

export default {
  components: { BarChart },
  mixins: [entitiesMetricsMixin, localQueryMixin],
  data() {
    return {
      sliMetrics: [],
      query: {
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

    datasets() {
      const { downtime, maintenance, uptime } = this.sliMetrics.reduce((acc, metric) => {
        const x = metric.timestamp * 1000;

        acc.downtime.push({ x, y: this.convertSecondsToPercent(metric.downtime) });
        acc.maintenance.push({ x, y: this.convertSecondsToPercent(metric.maintenance) });
        acc.uptime.push({ x, y: this.convertSecondsToPercent(metric.uptime) });

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
            },
            time: {
              unit: 'day',
              displayFormats: {
                day: DATETIME_FORMATS.short,
              },
            },
          },
          y: {
            stacked: true,
            max: 100,
            beginAtZero: true,
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
              title: () => '',
              label: tooltip => `${tooltip.raw.y}% ${tooltip.dataset.label.toLowerCase()}`,
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
    convertSecondsToPercent(value) {
      /**
       * We have 86400 seconds in day, 86400 / 100 = 864
       */
      return convertNumberToFixedString(value / 864, 2);
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
          sampling: KPI_SAMPLINGS.day,
        },
      });
    },
  },
};
</script>
