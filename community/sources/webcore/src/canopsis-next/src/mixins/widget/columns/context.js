import { CONTEXT_COLUMN_INFOS_PREFIX, CONTEXT_COLUMNS_WITH_SORTABLE } from '@/constants';

import { widgetColumnsMixin } from './common';

export const widgetColumnsContextMixin = {
  mixins: [widgetColumnsMixin],
  methods: {
    mapColumnEntity({ label, value, ...column }) {
      return {
        ...column,

        value,
        text: label,
        sortable: CONTEXT_COLUMNS_WITH_SORTABLE.includes(value) || value.startsWith(CONTEXT_COLUMN_INFOS_PREFIX),
      };
    },
  },
};
