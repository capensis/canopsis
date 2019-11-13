<template lang="pug">
  div.position-relative
    progress-overlay(:pending="pending")
    alert-overlay(
      :value="hasError",
      :message="serverErrorMessage"
    )
    stats-pareto-chart(:labels="labels", :datasets="datasets", :options="options")
</template>

<script>
import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import widgetStatsWrapperMixin from '@/mixins/widget/stats/stats-wrapper';
import widgetStatsChartWrapperMixin from '@/mixins/widget/stats/stats-chart-wrapper';

import { SORT_ORDERS } from '@/constants';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';
import AlertOverlay from '@/components/layout/alert/alert-overlay.vue';

import StatsParetoChart from './stats-pareto-chart.vue';

export default {
  components: {
    ProgressOverlay,
    AlertOverlay,
    StatsParetoChart,
  },
  mixins: [
    entitiesStatsMixin,
    widgetQueryMixin,
    widgetStatsWrapperMixin,
    widgetStatsChartWrapperMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: true,
      stats: null,
      total: 0,
    };
  },
  computed: {
    statTitle() {
      return this.widget.parameters.stat.title;
    },

    labels() {
      if (this.stats) {
        return this.stats.map(stat => stat.entity.name);
      }

      return [];
    },

    datasets() {
      if (this.stats) {
        const barsData = this.stats
          .map(stat => stat[this.statTitle].value);

        let sum = 0;

        const curveData = this.stats
          .reduce((acc, stat) => {
            sum += (stat[this.statTitle].value / this.total) * 100;
            acc.push(Math.round(sum));
            return acc;
          }, []);

        return [
          {
            label: 'curve',
            type: 'line',
            data: curveData,
            yAxisID: 'y-axis-2',
            backgroundColor: 'transparent',
            borderColor: this.widget.parameters.statsColors.Accumulation || 'rgba(0, 0, 0, 0.1)',
            cubicInterpolationMode: 'monotone',
          },
          {
            label: this.statTitle,
            data: barsData,
            yAxisID: 'y-axis-1',
            backgroundColor: this.widget.parameters.statsColors[this.statTitle] || 'rgba(0, 0, 0, 0.1)',
          },
        ];
      }

      return [];
    },

    options() {
      return {
        scales: {
          xAxes: [{
            ticks: {
              fontSize: 11,
            },
          }],
          yAxes: [
            {
              type: 'linear',
              position: 'left',
              id: 'y-axis-1',
              ticks: {
                suggestedMin: 0,
                fontSize: 11,
              },
              scaleLabel: {
                display: true,
                labelString: this.statTitle,
              },
            },
            {
              type: 'linear',
              position: 'right',
              id: 'y-axis-2',
              ticks: {
                suggestedMin: 0,
                fontSize: 11,
              },
              scaleLabel: {
                display: true,
                labelString: '%',
              },
            },
          ],
        },
      };
    },
  },
  methods: {
    getQuery() {
      const {
        mfilter,
        tstart,
        tstop,
        periodUnit,
        stats = {},
      } = this.getStatsQuery();

      const durationValue = tstop.diff(tstart, periodUnit);

      return {
        stats,
        mfilter,

        duration: `${durationValue}${periodUnit.toLowerCase()}`,
        tstop: tstop.startOf('h').unix(),
        sort_order: SORT_ORDERS.desc.toLowerCase(),
        sort_column: this.statTitle,
      };
    },

    async fetchList() {
      try {
        this.pending = true;
        this.serverErrorMessage = null;

        const { values, aggregations } = await this.fetchStatsListWithoutStore({
          params: this.getQuery(),
        });

        this.total = aggregations[this.statTitle].sum;
        this.stats = values.filter(stat => stat[this.statTitle].value);
      } catch (err) {
        this.serverErrorMessage = err.description || this.$t('errors.statsRequestProblem');
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
