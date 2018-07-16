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
  data() {
    return {
      pagination: {},
    };
  },
  computed: {
    limit() {
      return parseInt(this.$route.query.limit, 10) || PAGINATION_LIMIT;
    },
    first() {
      const { page } = this.$route.query;
      if (page === 1) {
        return 1;
      }
      if (this.$route.query.limit) {
        return page + (this.$route.$query.limit * (page - 1));
      }

      return (page - 1) + (PAGINATION_LIMIT * (page - 1));
    },
    last() {
      if (this.$route.query.page === 1) {
        return this.$route.query.limit || PAGINATION_LIMIT;
      }
      if (this.$route.query.limit) {
        return this.first + this.$route.query.limit;
      }

      return this.first + PAGINATION_LIMIT;
    },
  },
  watch: {
    $route: {
      immediate: true,
      handler() {
        this.fetchList();
      },
    },
    pagination: {
      handler(e) {
        let query = { ...this.$route.query };
        if (e.sortBy) {
          query.sort_key = e.sortBy;
          query.sort_dir = e.descending ? 'DESC' : 'ASC';
        } else {
          query = omit(this.$route.query, ['sort_key', 'sort_dir']);
        }
        this.$router.push({
          query,
        });
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
