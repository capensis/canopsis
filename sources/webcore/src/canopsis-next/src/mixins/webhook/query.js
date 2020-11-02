export default {
  data() {
    return {
      query: {
        search: '',
      },
    };
  },
  methods: {
    getQuery() {
      const query = {};

      const { search } = this.query;

      if (search) {
        query.search = search;
      }

      return query;
    },
  },
};
