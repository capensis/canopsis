import { ASSOCIATIVE_TABLES_NAMES } from '@/constants';

import { entitiesAssociativeTableMixin } from './index';

export const entitiesFilterHintsMixin = {
  mixins: [entitiesAssociativeTableMixin],
  methods: {
    fetchFilterHints() {
      return this.fetchAssociativeTable({ name: ASSOCIATIVE_TABLES_NAMES.filterHints });
    },
  },
};
