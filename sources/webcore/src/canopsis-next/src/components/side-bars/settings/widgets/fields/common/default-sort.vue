<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.defaultSortColumn') }}
    v-container
      v-select(
      v-if="!withoutColumns",
      :value="value.column",
      :items="columns",
      :label="$t('settings.columnName')",
      item-text="label",
      item-value="value",
      @change="updateField('column', $event)"
      )
      v-select(
      :value="selectedOrder",
      @input="updateValue",
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
      type: [Object, String],
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
    withoutColumns: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    orders() {
      return Object.values(SORT_ORDERS);
    },
    selectedOrder() {
      if (typeof this.value === 'string') {
        return SORT_ORDERS[this.value.toLowerCase()];
      }

      return this.value.order ? SORT_ORDERS[this.value.order.toLowerCase()] : SORT_ORDERS.desc;
    },
  },
  methods: {
    updateValue(value) {
      if (typeof this.value === 'string') {
        this.$emit('input', value);
      } else {
        this.updateField('order', value);
      }
    },
  },
};
</script>
