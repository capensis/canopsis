<template>
  <widget-settings-item :title="$t('settings.defaultSortColumn')">
    <v-layout>
      <v-select
        v-field="value.column"
        :items="columns"
        :label="columnsLabel"
      />
    </v-layout>
    <v-layout>
      <v-select
        v-field="value.order"
        :items="orders"
      />
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { SORT_ORDERS } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

/**
 * Component to select the default column to sort on settings
 *
 * @prop {Object} [value] - Object containing the default sort column's name and the sort direction
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
