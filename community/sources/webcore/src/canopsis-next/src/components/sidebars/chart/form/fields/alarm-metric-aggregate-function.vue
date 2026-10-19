<template>
  <widget-settings-item :title="$t('kpi.calculationMethod')">
    <c-radio-group
      :model-value="modelValue"
      :name="name"
      :label="label"
      class="mt-0"
      @update:model-value="$emit('update:modelValue', $event)"
    >
      <v-radio
        v-for="aggregateFunction in availableAggregateFunctions"
        :key="aggregateFunction.value"
        :label="aggregateFunction.label"
        :value="aggregateFunction.value"
        color="primary"
      />
    </c-radio-group>
  </widget-settings-item>
</template>

<script>
import { AGGREGATE_FUNCTIONS } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'aggregate_func',
    },
    aggregateFunctions: {
      type: Array,
      default: () => Object.values(AGGREGATE_FUNCTIONS),
    },
  },

  emits: ['update:modelValue'],

  computed: {
    availableAggregateFunctions() {
      return this.aggregateFunctions.map(value => ({
        value,
        label: this.$t(`kpi.aggregateFunctions.${value}`),
      }));
    },
  },
};
</script>
