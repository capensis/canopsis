import defaultItemsPerPageMixin from './default-items-per-page';

export default {
  mixins: [defaultItemsPerPageMixin],

  data() {
    return {
      pagination: {},
    };
  },
  watch: {
    defaultItemsPerPage: {
      immediate: true,
      handler(value) {
        this.pagination = {
          ...this.pagination,

          rowsPerPage: value,
        };
      },
    },
  },
};
