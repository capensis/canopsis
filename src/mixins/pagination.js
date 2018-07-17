import omit from 'lodash/omit';

import Pagination from '@/components/tables/pagination.vue';
import dateIntervals from '@/helpers/date-intervals';
import { PAGINATION_LIMIT } from '@/config';

/**
 * @mixin Add pagination logic, and dynamic route
 */
export default {
  components: {
    Pagination,
  },
  computed: {
    limit() {
      return parseInt(this.$route.query.limit, 10) || PAGINATION_LIMIT;
    },
    /**
     * Calculate first item nb to display on pagination, in case it's not given by the backend
     */
    first() {
      const { page } = this.$route.query;
      if (page === 1 || !this.$route.query.page) {
        return 1;
      }
      if (this.$route.query.limit) {
        return 1 + (this.$route.query.limit * (page - 1));
      }

      return 1 + (PAGINATION_LIMIT * (page - 1));
    },
    /**
     * Calculate last item nb to display on pagination, in case it's not given by the backend
     */
    last() {
      let last;

      if (this.$route.query.page === 1 || !this.$route.query.page) {
        last = this.$route.query.limit || PAGINATION_LIMIT;
      } else if (this.$route.query.limit) {
        last = this.$route.query.page * this.$route.query.limit;
      } else {
        last = this.$route.query.page * PAGINATION_LIMIT;
      }

      return last;
    },
  },
  watch: {
    $route: {
      immediate: true,
      handler() {
        this.fetchList();
      },
    },
  },
  methods: {
    getQuery() {
      const query = omit(this.$route.query, ['page', 'interval']);

      if (this.$route.query.interval && this.$route.query.interval !== 'custom') {
        try {
          const { tstart, tstop } = dateIntervals[this.$route.query.interval]();
          query.tstart = tstart;
          query.tstop = tstop;
        } catch (err) {
          console.warn(err);
        }
      }
      query.limit = this.limit;
      query.skip = ((this.$route.query.page - 1) * this.limit) || 0;

      return query;
    },
    fetchList() {
      this.fetchListAction({
        params: this.getQuery(),
      });
    },
  },
};
