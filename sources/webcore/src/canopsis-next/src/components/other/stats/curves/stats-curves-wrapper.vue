<template lang="pug">
  div.stats-wrapper
    stats-curves(v-if="!pending", :labels="labels", :datasets="datasets")
    v-layout(v-else, justify-center)
      v-progress-circular(
      indeterminate,
      color="primary",
      )
</template>

<script>
import { get, isString } from 'lodash';
import moment from 'moment';

import { DATETIME_FORMATS, STATS_DURATION_UNITS, STATS_DEFAULT_COLOR } from '@/constants';
import { dateParse } from '@/helpers/date-intervals';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import StatsCurves from './stats-curves.vue';

export default {
  components: {
    StatsCurves,
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
      statsValues: {},
      pending: true,
    };
  },
  computed: {
    labels() {
      const labels = [];
      if (this.statsValues) {
        const stats = Object.keys(this.statsValues);
        const values = { ...this.statsValues };
        /*
        'start' correspond to the beginning timestamp.
        It's the same for all stats, that's why we can just take the first.
        We then give it to the date filter, to display it with a date format
        */
        values[stats[0]].sum.map(value => labels.push(this.$options.filters.date(value.end, 'long')));
        return labels;
      }
      return labels;
    },
    datasets() {
      if (this.statsValues) {
        return Object.keys(this.statsValues).reduce((acc, stat) => {
          const values = this.statsValues[stat].sum.map(value => value.value);
          acc.push({
            data: values,
            label: stat,
            borderColor: get(this.widget.parameters, `statsColors.${stat}`, STATS_DEFAULT_COLOR),
            backgroundColor: 'transparent',
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

      this.statsValues = aggregations;
      this.pending = false;
    },
  },
};
</script>
