import Pagination from '@/components/tables/pagination.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';

export default {
  components: {
    Pagination,
    RecordsPerPage,
  },
  methods: {
    updateQueryPage(page) {
      this.query = { ...this.query, page };
    },
  },
};
