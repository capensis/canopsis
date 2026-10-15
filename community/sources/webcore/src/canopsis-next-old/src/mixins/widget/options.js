import { PAGINATION_LIMIT } from '@/config';

import { convertDataTableOptionsToQuery } from '@/helpers/entities/shared/query';

export const widgetOptionsMixin = {
  computed: {
    options: {
      get() {
        const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [] } = this.query;

        return { page, itemsPerPage, sortBy, sortDesc };
      },

      set(newOptions) {
        const convertedOptions = convertDataTableOptionsToQuery(newOptions, this.options);

        if (convertedOptions === this.options) {
          return;
        }

        const newQuery = { ...this.query, ...convertedOptions };

        if (Object.keys(this.$props).includes('query')) {
          this.$emit('update:query', newQuery);

          return;
        }

        this.query = newQuery;
      },
    },
  },
};
