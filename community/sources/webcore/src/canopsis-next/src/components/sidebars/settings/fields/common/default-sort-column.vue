<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.defaultSortColumn') }}
    v-container
      v-layout(row)
        v-select(
          v-field="value.column",
          :items="columns",
          :label="columnsLabel"
        )
      v-layout(row)
        v-select(
          v-field="value.order",
          :items="orders"
        )
</template>

<script>
import { SORT_ORDERS } from '@/constants';

/**
 * Component to select the default column to sort on settings
 *
 * @prop {Object} [value] - Object containing the default sort column's name and the sort direction
 *
 * @event value#input
 */
export default {
  props: {
    value: {
      type: Object,
      default: () => ({
        column: '',
        order: SORT_ORDERS.asc,
      }),
    },
    columns: {
      type: Array,
      default: () => [],
    },
    columnsLabel: {
      type: String,
      default: null,
    },
  },
  computed: {
    orders() {
      return Object.values(SORT_ORDERS);
    },
  },
};
</script>
