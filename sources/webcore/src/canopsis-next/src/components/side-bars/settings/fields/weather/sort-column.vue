<template lang="pug">
  v-list-group(data-test="defaultSortColumn")
    v-list-tile(slot="activator") {{ $t('settings.defaultSortColumn') }}
    v-container
      v-layout(data-test="defaultSortColumnOrderByLayout", row)
        v-combobox(
          v-field="value.column",
          :items="columns",
          :label="columnsLabel",
          :return-object="false",
          item-text="label",
          item-value="value"
        )
          template(slot="no-data")
            v-list-tile
              v-list-tile-content
                v-list-tile-title(v-html="$t('settings.sortColumnNoData')")
      v-layout(data-test="defaultSortColumnOrdersLayout", row)
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
 * @prop {Array} [columns] - List of columns suggestions to sort on
 * @prop {String} [columnsLabel] - Setting's title to display on settings panel
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
