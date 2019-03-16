import widgetStatsQueryMixin from './stats-query';

export default {
  mixins: [widgetStatsQueryMixin],
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
      return [];
    },
  },
  methods: {
    getQuery() {
      const {
        stats,
        mfilter,
        tstart,
        tstop,
        periodUnit,
        periodValue,
        duration,
      } = this.getStatsQuery();

      return {
        stats,
        duration,
        mfilter,

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
