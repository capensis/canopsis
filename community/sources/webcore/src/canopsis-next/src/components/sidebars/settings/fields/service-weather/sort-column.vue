<template lang="pug">
  widget-settings-item(:title="$t('settings.defaultSortColumn')")
    v-layout(row)
      v-combobox(
        v-field="value.column",
        :items="columns",
        :label="$t('settings.orderBy')",
        :return-object="false",
        item-text="label",
        item-value="value"
      )
        template(#no-data="")
          v-list-tile
            v-list-tile-content
              v-list-tile-title(v-html="$t('settings.sortColumnNoData')")
    v-layout(row)
      v-select(
        v-field="value.order",
        :items="orders"
      )
</template>

<script>
import { SORT_ORDERS } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

/**
 * Component to select the default column to sort on settings
 *
 * @prop {Object} [value] - Object containing the default sort column's name and the sort direction
 * @prop {Array} [columns] - List of columns suggestions to sort on
 *
 * @event value#input
 */
export default {
  components: { WidgetSettingsItem },
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
  },
  computed: {
    orders() {
      return Object.values(SORT_ORDERS);
    },
  },
};
</script>
