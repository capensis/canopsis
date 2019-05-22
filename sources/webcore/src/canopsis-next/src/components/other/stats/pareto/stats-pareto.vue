<template lang="pug">
  div.stats-wrapper
    progress-overlay(:pending="pending")
    stats-pareto-chart(:labels="labels", :datasets="datasets", :options="options")
</template>

<script>
import moment from 'moment-timezone';
import { get, isString, omit } from 'lodash';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';

import { dateParse } from '@/helpers/date-intervals';

import { DATETIME_FORMATS, STATS_DURATION_UNITS } from '@/constants';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import StatsParetoChart from './stats-pareto-chart.vue';

export default {
  components: {
    ProgressOverlay,
    StatsParetoChart,
  },
  mixins: [
    entitiesStatsMixin,
    widgetQueryMixin,
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

    filteredStats() {
      if (this.stats) {
        return this.stats.filter(stat => stat[this.statTitle].value);
      }

      return null;
    },

    labels() {
      if (this.filteredStats) {
        return this.filteredStats.map(stat => stat.entity.name);
      }

      return [];
    },

    datasets() {
      if (this.stats) {
        const barsData = this.filteredStats
          .map(stat => stat[this.statTitle].value);

        let sum = 0;

        const curveData = this.filteredStats
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
      const { annotationLine } = this.widget.parameters;
      const options = {
        scales: {
          xAxes: [{
            scaleLabel: {
              display: true,
              labelString: 'Entities',
            },
          }],
          yAxes: [
            {
              type: 'linear',
              position: 'left',
              id: 'y-axis-1',
              ticks: {
                suggestedMin: 0,
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
              },
              scaleLabel: {
                display: true,
                labelString: '%',
              },
            },
          ],
        },
      };

      if (annotationLine && annotationLine.enabled) {
        options.annotation = {
          annotations: [{
            type: 'line',
            mode: 'horizontal',
            scaleID: 'y-axis-0',
            value: annotationLine.value,
            borderColor: annotationLine.lineColor,
            borderWidth: 2,
            label: {
              enabled: true,
              position: 'left',
              fontSize: 10,
              xPadding: 5,
              yPadding: 5,
              content: annotationLine.label,
              backgroundColor: annotationLine.labelColor,
            },
          }],
        };
      }
      return options;
    },
  },
  methods: {
    getStatsQuery() {
      const { dateInterval, stat, mfilter } = this.query;
      const { periodValue } = dateInterval;
      let { periodUnit, tstart, tstop } = dateInterval;
      let filter = get(mfilter, 'filter', {});

      if (isString(filter)) {
        try {
          filter = JSON.parse(filter);
        } catch (err) {
          filter = {};

          console.error(err);
        }
      }

      tstart = dateParse(tstart, 'start', DATETIME_FORMATS.dateTimePicker);
      tstop = dateParse(tstop, 'stop', DATETIME_FORMATS.dateTimePicker);

      if (periodUnit === STATS_DURATION_UNITS.month) {
        periodUnit = periodUnit.toUpperCase();

        /**
         * If period unit is 'month', we need to put the dates at the first day of the month, at 00:00 UTC
         * And add the difference between the local date, and the UTC one.
         */
        tstart = moment.utc(tstart).startOf('month').tz(moment.tz.guess());
        tstop = moment.utc(tstop).startOf('month').tz(moment.tz.guess());
      }

      const stats = {};

      stats[stat.title] = {
        ...omit(stat, ['title']),
        stat: stat.stat.value,
        aggregate: ['sum'],
      };

      return {
        stats,
        filter,
        tstart,
        tstop,
        periodUnit,
        periodValue,

        mfilter: filter,
        duration: `${periodValue}${periodUnit.toLowerCase()}`,
      };
    },

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
        sort_order: 'desc',
        sort_column: this.statTitle,
      };
    },

    async fetchList() {
      this.pending = true;

      const { values, aggregations } = await this.fetchStatsListWithoutStore({
        params: this.getQuery(),
      });

      this.total = aggregations[this.statTitle].sum;
      this.stats = values;
      this.pending = false;
    },
  },
};
</script>

