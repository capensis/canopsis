import { SORT_ORDERS } from '@/constants';

export default ({
  field = 'pagination',
  queryField = 'query',
  defaultSortKey = null,
  mutating = false,
} = {}) => ({
  computed: {
    [field]: {
      get() {
        const descending = this[queryField].sortDir === SORT_ORDERS.desc;
        const sortBy = this[queryField].sortKey || defaultSortKey;

        return { sortBy, descending };
      },
      set(value) {
        const isNotEqualSortBy = value.sortBy !== this[field].sortBy;
        const isNotEqualDescending = value.descending !== this[field].descending;

        if (isNotEqualSortBy || isNotEqualDescending) {
          const newQuery = {
            ...this[queryField],

            sortKey: value.sortBy,
            sortDir: value.descending ? SORT_ORDERS.desc : SORT_ORDERS.asc,
          };

          if (mutating) {
            this[queryField] = newQuery;
          } else {
            this.$emit(`update:${queryField}`, newQuery);
          }
        }
      },
    },
  },
});
