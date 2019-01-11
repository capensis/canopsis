import EventBus from '@/event-bus';
import { PAGINATION_LIMIT, PAGINATION_PER_PAGE_VALUES } from '@/config';

export default {
  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      items: [],
      meta: {},
      pending: false,
      pagination: {
        page: 1,
        rowsPerPage: PAGINATION_LIMIT,
      },
    };
  },
  computed: {
    rowsPerPageItems() {
      return PAGINATION_PER_PAGE_VALUES;
    },
  },
  watch: {
    pagination(value, oldValue) {
      if (value.page !== oldValue.page || value.rowsPerPage !== oldValue.rowsPerPage) {
        this.fetchList();
      }
    },
  },
  created() {
    EventBus.$on('filter-editor:results:fetch', this.fetchList);
  },
  beforeDestroy() {
    EventBus.$off('filter-editor:results:fetch', this.fetchList);
  },
};
