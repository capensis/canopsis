<template lang="pug">
  div(v-bind="component.bind")
</template>

<script>
import alarmColumnFiltersMixin from '@/mixins/entities/alarm-column-filters';

export default {
  mixins: [alarmColumnFiltersMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    column: {
      type: Object,
      required: true,
    },
  },
  computed: {
    columnFiltersMap() {
      return this.alarmColumnFilters.reduce((acc, { column, filter, attributes = [] }) => {
        acc[column.replace(/^entity\./, '')] = this.getFilter(filter, attributes);

        return acc;
      }, {});
    },

    columnFilter() {
      return this.columnFiltersMap[this.column.value];
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
          text: this.value,
        },
      };
    },
  },
};
</script>
