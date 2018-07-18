import omit from 'lodash/omit';

import Pagination from '@/components/tables/pagination.vue';
import dateIntervals from '@/helpers/date-intervals';
import { PAGINATION_LIMIT } from '@/config';

const QUERY_FIELDS_MAP = {
  page: 'p',
  limit: 'l',
  interval: 'i',
  sort_key: 'sk',
  sort_dir: 'sd',
};

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
      queryPrefix: `${Math.random()}`,
    };
  },
  computed: {
    limit() {
      return parseInt(this.getQueryValue('limit'), 10) || PAGINATION_LIMIT;
    },
    /**
     * Calculate first item nb to display on pagination, in case it's not given by the backend
     */
    first() {
      const page = this.getQueryValue('page') || 1;

      return 1 + (this.limit * (page - 1));
    },
    /**
     * Calculate last item nb to display on pagination, in case it's not given by the backend
     */
    last() {
      const page = this.getQueryValue('page') || 1;

      return page * this.limit;
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
        const sortKeyKey = this.getQueryKey('sort_key');
        const sortDirKey = this.getQueryKey('sort_dir');
        let query = { ...this.$route.query };

        if (e.sortBy) {
          query[sortKeyKey] = e.sortBy;
          query[sortDirKey] = e.descending ? 'DESC' : 'ASC';
        } else {
          query = omit(this.$route.query, [sortKeyKey, sortDirKey]);
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
      const interval = this.getQueryValue('interval');
      const page = this.getQueryValue('page');

      if (interval && interval !== 'custom') {
        try {
          const { tstart, tstop } = dateIntervals[interval]();
          query.tstart = tstart;
          query.tstop = tstop;
        } catch (err) {
          console.warn(err);
        }
      }
      query.limit = this.limit;
      query.skip = ((page - 1) * this.limit) || 0;

      return query;
    },
    getQueryValue(key) {
      return this.$route.query[this.getQueryKey(key)];
    },
    getQueryKey(key) {
      return this.queryPrefix + QUERY_FIELDS_MAP[key];
    },
    fetchList() {
      this.fetchListAction({
        params: this.getQuery(),
      });
    },
  },
};
