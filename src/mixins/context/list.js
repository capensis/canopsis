import omit from 'lodash/omit';
import ContextMixin from '@/mixins/context/index';
import PaginationMixin from '@/mixins/query';

export default {
  mixins: [
    ContextMixin,
    PaginationMixin,
  ],
  methods: {
    getQuery() {
      const query = omit(this.$route.query, ['page', 'sort_dir', 'sort_key']);
      query.limit = this.limit;
      query.start = ((this.$route.query.page - 1) * this.limit) || 0;

      if (this.$route.query.sort_key) {
        query.sort = [{
          property: this.$route.query.sort_key,
          direction: this.$route.query.sort_dir ? this.$route.query.sort_dir : 'ASC',
        }];
      }

      return query;
    },
    fetchList() {
      this.fetchContextEntities({
        params: this.getQuery(),
      });
    },
  },
};
