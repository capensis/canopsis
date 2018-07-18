import omit from 'lodash/omit';

import Pagination from '@/components/tables/pagination.vue';
import dateIntervals from '@/helpers/date-intervals';
import { PAGINATION_LIMIT } from '@/config';

/**
 * @mixin Add pagination logic
 */
export default {
  components: {
    Pagination,
  },
  data() {
    const query = {
      page: 1,
      limit: 10,
    };

    if (this.widget.default_sort_column) {
      query.sort_key = `v.${this.widget.default_sort_column.property}`;
      query.sort_dir = this.widget.default_sort_column.direction;
    }

    return {
      selected: [],
      query,
      vDataTablePagination: {
        page: query.page,
        rowsPerPage: query.limit,
        sortBy: query.sort_key,
        descending: query.sort_dir === 'DESC',
      },
    };
  },
  computed: {
    limit() {
      return this.query.limit || PAGINATION_LIMIT;
    },
    page() {
      return this.query.page || 1;
    },
    /**
     * Calculate first item nb to display on pagination, in case it's not given by the backend
     */
    first() {
      return 1 + (this.limit * (this.page - 1));
    },
    /**
     * Calculate last item nb to display on pagination, in case it's not given by the backend
     */
    last() {
      return this.page * this.limit;
    },
  },
  watch: {
    query: {
      immediate: true,
      handler() {
        this.fetchList();
      },
    },
    vDataTablePagination(value, oldValue) {
      if (value.sortBy !== oldValue.sortBy || value.descending !== oldValue.descending) {
        let query = { ...this.query };

        if (value.sortBy) {
          query.sort_key = value.sortBy;
          query.sort_dir = value.descending ? 'DESC' : 'ASC';
        } else {
          query = omit(this.query, ['sort_key', 'sort_dir']);
        }

        this.query = query;
      }
    },
  },
  methods: {
    getQuery() {
      const query = omit(this.query, ['page', 'interval']);
      const { page, interval } = this.query;

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
    fetchList() {
      this.fetchListAction({
        params: this.getQuery(),
      });
    },
  },
};
