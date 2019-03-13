<template lang="pug">
  div.stats-wrapper
    stats-curves(v-if="!pending", :labels="labels", :datasets="datasets.data")
    v-layout(v-else, justify-center)
      v-progress-circular(
      indeterminate,
      color="primary",
      )
</template>

<script>
import moment from 'moment';

import { parseStringToDateInterval } from '@/helpers/date-intervals';

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
        const data = Object.keys(this.statsValues).reduce((acc, stat) => {
          const values = this.statsValues[stat].sum.map(value => value.value);
          acc.push({
            data: values,
            label: stat,
            borderColor: this.widget.parameters.statsColors ? this.widget.parameters.statsColors[stat] : '#DDDDDD',
            backgroundColor: 'transparent',
          });
          return acc;
        }, []);

        return {
          data,
        };
      }

      return {};
    },
  },
  methods: {
    // Determine if tstart and tstop are valid Dates or Dynamic Date strings (Ex: 'now')
    dateParse(date, type) {
      if (!moment(date).isValid()) {
        try {
          return parseStringToDateInterval(date, type);
        } catch (err) {
          // TODO: DISPLAY AN ALERT TO THE USER
          console.warn(err);
          return err;
        }
      } else {
        return moment(date);
      }
    },

    async fetchList() {
      this.pending = true;
      const params = {};
      const { dateInterval, mfilter, stats } = this.getQuery();
      const { periodValue } = dateInterval;
      let { periodUnit, tstart, tstop } = dateInterval;

      tstart = this.dateParse(tstart, 'start');
      tstop = this.dateParse(tstop, 'stop');


      if (periodUnit === 'm') {
        periodUnit = periodUnit.toUpperCase();
        // If period unit is 'month', we need to put the dates at the first day of the month, at 00:00 UTC
        const monthlyRoundedTstart = moment.tz(tstart, moment.tz.guess()).startOf('month');
        // Add the difference between the local date, and the UTC one.
        tstart = monthlyRoundedTstart.add(monthlyRoundedTstart.utcOffset(), 'm');
        const monthlyRoundedTstop = moment.tz(tstop, moment.tz.guess()).startOf('month');
        // Add the difference between the local date, and the UTC one.
        tstop = monthlyRoundedTstop.add(monthlyRoundedTstop.utcOffset(), 'm');
      }

      params.duration = `${periodValue}${periodUnit.toLowerCase()}`;
      params.periods = Math.ceil((tstop.diff(tstart, periodUnit) + 1) / periodValue);
      params.stats = stats;
      params.mfilter = mfilter.filter ? JSON.parse(mfilter.filter) : {};
      params.tstop = tstop.startOf('h').unix();

      const { aggregations } = await this.fetchStatsEvolutionWithoutStore({
        params,
        aggregate: ['sum'],
      });

      this.statsValues = aggregations;
      this.pending = false;
    },
  },
};
</script>
