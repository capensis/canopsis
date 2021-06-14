import { localQueryMixin } from '@/mixins/query-local/query';

export default {
  mixins: [localQueryMixin],
  methods: {
    customQueryCondition(query, oldQuery) {
      return query.search !== oldQuery.search;
    },

    getQuery({ search } = this.query) {
      const query = {};

      if (search) {
        query.search = search;
      }

      return query;
    },
  },
};
