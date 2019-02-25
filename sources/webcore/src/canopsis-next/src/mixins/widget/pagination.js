import Pagination from '@/components/tables/pagination.vue';

export default {
  components: {
    Pagination,
  },
  methods: {
    updateQueryPage(page) {
      this.query = { ...this.query, page };
    },
  },
};
