<template>
  <widget-settings-item :title="$t('settings.defaultSortColumn')">
    <v-layout>
      <v-combobox
        v-field="value.column"
        :items="columns"
        :label="$t('settings.orderBy')"
        :return-object="false"
        item-text="label"
        item-value="value"
      >
        <template #no-data="">
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title v-html="$t('settings.sortColumnNoData')" />
            </v-list-item-content>
          </v-list-item>
        </template>
      </v-combobox>
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
