<template>
  <widget-settings-item :title="$t('settings.defaultSortColumn')">
    <c-layout gap="2" column>
      <c-select
        :model-value="modelValue.column"
        :items="columns"
        :label="columnsLabel"
        :item-title="itemTitle"
        :item-value="itemValue"
        hide-details
        @update:model-value="updateField('column', $event)"
      />
      <c-select
        :model-value="modelValue.order"
        :items="orders"
        hide-details
        @update:model-value="updateField('order', $event)"
      />
    </c-layout>
  </widget-settings-item>
</template>

<script>
import { SORT_ORDERS } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

/**
 * Component to select the default column to sort on settings
 *
 * @prop {Object} [modelValue] - Object containing the default sort column's name and the sort direction
 *
 * @event modelValue#update:modelValue
 */
export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
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
    itemTitle: {
      type: String,
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
  },

  emits: ['update:modelValue'],

  computed: {
    orders() {
      return Object.values(SORT_ORDERS);
    },
  },

  methods: {
    updateField(field, value) {
      this.$emit('update:modelValue', {
        ...this.modelValue,
        [field]: value,
      });
    },
  },
};
</script>
