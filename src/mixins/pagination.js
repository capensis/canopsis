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
      const page = this.$route.query.page || 1;
      const limit = this.$route.query.limit || PAGINATION_LIMIT;

      if (page === 1) {
        return 1;
      }

      return 1 + (limit * (page - 1));
    },
    /**
     * Calculate last item nb to display on pagination, in case it's not given by the backend
     */
    last() {
      const page = this.$route.query.page || 1;
      const limit = this.$route.query.limit || PAGINATION_LIMIT;

      return page * limit;
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
