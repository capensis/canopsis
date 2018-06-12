import omit from 'lodash/omit';
import merge from 'lodash/merge';

import { PAGINATION_LIMIT } from '@/config';
import Pagination from '@/components/basic-component/pagination.vue';

export default {
  components: {
    Pagination,
  },
  computed: {
    limit() {
      return parseInt(this.$route.query.limit, 10) || PAGINATION_LIMIT;
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
      const query = omit(this.$route.query, ['page']);
      if (this.queries) {
        merge(query, this.queries);
      }
      return query;
    },
    fetchList() {
      this.fetchListAction({
        params: this.getQuery(),
      });
    },
  },
};
