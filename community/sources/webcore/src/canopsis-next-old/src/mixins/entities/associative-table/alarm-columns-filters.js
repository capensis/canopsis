import { get } from 'lodash';

import { ASSOCIATIVE_TABLES_NAMES } from '@/constants';

import { entitiesAssociativeTableMixin } from './index';

export const entitiesAlarmColumnsFiltersMixin = {
  mixins: [entitiesAssociativeTableMixin],
  methods: {
    async fetchAlarmColumnsFiltersList() {
      const content = await this.fetchAssociativeTable({
        name: ASSOCIATIVE_TABLES_NAMES.alarmColumnsFilters,
      });

      return get(content, 'filters', []);
    },
  },
};
