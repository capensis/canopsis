<template>
  <c-layout>
    <c-enabled-field
      :model-value="Boolean(modelValue)"
      :label="$t('testSuite.compareWithHistorical')"
      @update:model-value="changeEnabledField"
    />
    <c-radio-group
      v-if="modelValue"
      :model-value="modelValue"
      inline
      @update:model-value="$emit('update:modelValue', $event)"
    >
      <v-radio
        v-for="item in preparedItems"
        :key="item.value"
        :value="item.value"
        :label="item.label"
        :color="item.color"
      />
    </c-radio-group>
  </c-layout>
</template>

<script>
import { TEST_SUITE_HISTORICAL_DATA_MONTHS_DEFAULT_ITEMS } from '@/constants';

export default {
  props: {
    modelValue: {
      type: Number,
      default: 0,
    },
    items: {
      type: Array,
      default: () => TEST_SUITE_HISTORICAL_DATA_MONTHS_DEFAULT_ITEMS,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    preparedItems() {
      return this.items.map((value) => {
        const years = value / 12;

        return {
          value,
          color: 'primary',
          label: value % 12 === 0
            ? `${years} ${this.$tc('common.times.year', years)}`
            : `${value} ${this.$tc('common.times.month', value)}`,
        };
      });
    },
  },
  methods: {
    changeEnabledField(value) {
      this.$emit('update:modelValue', Number(value));
    },
  },
};
</script>
