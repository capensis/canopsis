import { isEqual } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import {
  convertDataTableSortToQuery,
  convertQuerySortToDataTableSort,
} from '@/helpers/entities/shared/query';

export const widgetOptionsMixin = {
  computed: {
    options: {
      get() {
        const {
          page = 1,
          itemsPerPage = PAGINATION_LIMIT,
          sortBy = [],
          sortDesc = [],
        } = this.query;

        return {
          page,
          itemsPerPage,
          sortBy: convertQuerySortToDataTableSort(sortBy, sortDesc),
        };
      },

      set(newOptions = {}) {
        const {
          page = 1,
          itemsPerPage = PAGINATION_LIMIT,
          sortBy = [],
        } = newOptions;
        const sort = convertDataTableSortToQuery(sortBy);
        const oldItemsPerPage = this.query.itemsPerPage ?? PAGINATION_LIMIT;
        const convertedOptions = {
          page: itemsPerPage <= oldItemsPerPage ? page : 1,
          itemsPerPage,
          sortBy: sort.sortBy,
          sortDesc: sort.sortDesc,
        };
        const oldOptions = {
          page: this.query.page ?? 1,
          itemsPerPage: oldItemsPerPage,
          sortBy: this.query.sortBy ?? [],
          sortDesc: this.query.sortDesc ?? [],
        };

        if (isEqual(convertedOptions, oldOptions)) {
          return;
        }

        if (itemsPerPage !== oldItemsPerPage) {
          this.updateContentInUserPreference?.({ itemsPerPage });
        }

        this.query = {
          ...this.query,
          ...convertedOptions,
        };
      },
    },
  },

  methods: {
    updateOptions(options) {
      this.options = options;
    },
  },
};
