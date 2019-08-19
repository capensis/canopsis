<template lang="pug">
  v-list-group(data-test="defaultSortColumn")
    v-list-tile(slot="activator") {{ $t('settings.defaultSortColumn') }}
    v-container(data-test="sortContainer")
      v-select(
      :value="value.column",
      :items="columns",
      :label="columnsLabel",
      item-text="label",
      item-value="value",
      @change="updateField('column', $event)"
      )
      v-select(
      :value="value.order",
      @input="updateField('order', $event)",
      :items="orders"
      )
</template>

<script>
import formMixin from '@/mixins/form';

import { SORT_ORDERS } from '@/constants';

/**
* Component to select the default column to sort on settings
*
* @prop {Object} [value] - Object containing the default sort column's name and the sort direction
*
* @event value#input
*/
export default {
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      default: () => ({
        column: '',
        order: 'ASC',
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
