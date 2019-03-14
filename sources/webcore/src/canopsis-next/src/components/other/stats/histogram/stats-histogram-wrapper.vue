<template lang="pug">
  div.stats-wrapper
    progress-overlay(:pending="pending")
    v-fade-transition
      stats-histogram(:labels="labels", :datasets="datasets")
</template>

<script>
import moment from 'moment';
import { get, isString } from 'lodash';

import { DATETIME_FORMATS, STATS_DURATION_UNITS, STATS_DEFAULT_COLOR } from '@/constants';

import { dateParse } from '@/helpers/date-intervals';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import StatsHistogram from './stats-histogram.vue';

export default {
  components: {
    ProgressOverlay,
    StatsHistogram,
  },
  mixins: [entitiesStatsMixin, widgetQueryMixin, entitiesUserPreferenceMixin],
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
    };
  },
  computed: {
    labels() {
      if (this.stats) {
        const stats = Object.keys(this.stats);

        /**
        'start' correspond to the beginning timestamp.
        It's the same for all stats, that's why we can just take the first.
        We then give it to the date filter, to display it with a date format
        */
        return this.stats[stats[0]].sum.map(value => this.$options.filters.date(value.end, 'long', true));
      }

      return [];
    },
    datasets() {
      if (this.stats) {
        return Object.keys(this.stats).reduce((acc, stat) => {
          const values = this.stats[stat].sum.map(value => value.value);

          acc.push({
            data: values,
            label: stat,
            backgroundColor: get(this.widget.parameters, `statsColors.${stat}`, STATS_DEFAULT_COLOR),
          });
          return acc;
        }, []);
      }

      return [];
    },

  },
  methods: {
    getQuery() {
      const { dateInterval, stats, mfilter } = this.query;
      const { periodValue } = dateInterval;
      let { periodUnit, tstart, tstop } = dateInterval;
      let filter = get(mfilter, 'filter', {});

      if (isString(filter)) {
        filter = JSON.parse(filter);
      }

      tstart = dateParse(tstart, 'start', DATETIME_FORMATS.picker);
      tstop = dateParse(tstop, 'stop', DATETIME_FORMATS.picker);

      if (periodUnit === STATS_DURATION_UNITS.month) {
        periodUnit = periodUnit.toUpperCase();

        /**
         * If period unit is 'month', we need to put the dates at the first day of the month, at 00:00 UTC
         * And add the difference between the local date, and the UTC one.
         */
        tstart = moment.utc(tstart).startOf('month').tz(moment.tz.guess());
        tstop = moment.utc(tstop).startOf('month').tz(moment.tz.guess());
      }

      return {
        stats,

        mfilter: filter,
        duration: `${periodValue}${periodUnit.toLowerCase()}`,
        periods: Math.ceil((tstop.diff(tstart, periodUnit) + 1) / periodValue),
        tstop: tstop.startOf('h').unix(),
      };
    },

    async fetchList() {
      this.pending = true;

      const { aggregations } = await this.fetchStatsEvolutionWithoutStore({
        params: this.getQuery(),
        aggregate: ['sum'],
      });

      this.stats = aggregations;
      this.pending = false;
    },
  },
};
</script>

<style lang="scss" scoped>
  .stats-wrapper {
    position: relative;
  }
</style>
