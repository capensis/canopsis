<template lang="pug">
  div(v-bind="component.bind")
</template>

<script>
import { get } from 'lodash';

import { widgetColumnsFiltersMixin } from '@/mixins/widget/columns-filters';

import EntityColumnEventStatistics from './entity-column-event-statistics.vue';
import EntityColumnPbehaviorInfo from './entity-column-pbehavior-info.vue';

export default {
  components: {
    EntityColumnEventStatistics,
    EntityColumnPbehaviorInfo,
  },
  mixins: [widgetColumnsFiltersMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    column: {
      type: Object,
      required: true,
    },
    columnsFilters: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    columnsFiltersMap() {
      return this.columnsFilters.reduce((acc, { column, filter, attributes = [] }) => {
        acc[column.replace(/^entity\./, '')] = this.getFilter(filter, attributes);

        return acc;
      }, {});
    },

    columnFilter() {
      return this.columnsFiltersMap[this.column.value];
    },

    value() {
      const value = get(this.entity, this.column.value, '');

      return this.columnFilter ? this.columnFilter(value) : value;
    },

    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        enabled: {
          bind: {
            is: 'c-enabled',
            value: this.value,
          },
        },
        idle_since: {
          bind: {
            is: 'c-no-events-icon',
            value: Number(this.value),
            top: true,
          },
        },
        ko_events: {
          bind: {
            is: 'entity-column-event-statistics',
            entity: this.entity,
          },
        },
        ok_events: {
          bind: {
            is: 'entity-column-event-statistics',
            entity: this.entity,
          },
        },
        pbehavior_info: {
          bind: {
            is: 'entity-column-pbehavior-info',
            pbehaviorInfo: this.entity.pbehavior_info,
          },
        },
      };

      if (PROPERTIES_COMPONENTS_MAP[this.column.value]) {
        return PROPERTIES_COMPONENTS_MAP[this.column.value];
      }

      return {
        bind: {
          is: 'c-ellipsis',
          text: String(this.value),
        },
      };
    },
  },
};
</script>
