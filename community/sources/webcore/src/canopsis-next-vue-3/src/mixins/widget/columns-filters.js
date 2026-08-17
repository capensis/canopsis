import { FILTERS } from '@/helpers/filter';

export const widgetColumnsFiltersMixin = {
  props: {
    columnsFilters: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    columnsFiltersMap() {
      return this.columnsFilters.reduce((acc, { column, filter, attributes = [] }) => {
        acc[column] = this.getFilter(filter, attributes);

        return acc;
      }, {});
    },
  },
  methods: {
    getFilter(filter, attributes = []) {
      const filterFunc = FILTERS[filter];

      return value => (filterFunc ? filterFunc(value, ...attributes) : value);
    },
  },
};
