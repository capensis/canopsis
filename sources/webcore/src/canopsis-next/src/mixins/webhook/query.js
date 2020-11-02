import { PAGINATION_LIMIT } from '@/config';

export default {
  data() {
    return {
      query: {
        page: 1,
        limit: PAGINATION_LIMIT,
        search: '',
      },
    };
  },
  methods: {
    getQuery() {
      const query = {};

      const {
        page,
        limit = PAGINATION_LIMIT,
        search,
      } = this.query;

      query.limit = limit;
      query.skip = ((page - 1) * limit) || 0;

      if (search) {
        query.search = search;
      }

      return query;
    },
  },
};
