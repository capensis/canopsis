<template lang="pug">
  div(v-bind="component.bind")
</template>

<script>
import { widgetColumnsFiltersMixin } from '@/mixins/widget/columns-filters';

export default {
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
      return this.$options.filters.get(this.entity, this.column.value, this.columnFilter, '');
    },

    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        enabled: {
          bind: {
            is: 'c-enabled',
            value: this.value,
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
