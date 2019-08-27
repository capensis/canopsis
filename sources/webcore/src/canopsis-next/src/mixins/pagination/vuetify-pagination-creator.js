import { PAGINATION_LIMIT } from '@/config';

import defaultItemsPerPageMixin from './default-items-per-page';

export default function (itemsKey) {
  return {
    mixins: [defaultItemsPerPageMixin],

    data() {
      return {
        pagination: {
          rowsPerPage: PAGINATION_LIMIT,
          page: 1,
        },
      };
    },
    watch: {
      [itemsKey](value) {
        this.$set(this.pagination, 'totalItems', value.length);
      },

      defaultItemsPerPage: {
        immediate: true,
        handler(value) {
          this.$set(this.pagination, 'rowsPerPage', value);
        },
      },
    },
  };
}
