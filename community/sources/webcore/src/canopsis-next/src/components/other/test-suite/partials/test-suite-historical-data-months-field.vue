<template>
  <v-layout>
    <c-enabled-field
      :value="Boolean(value)"
      :label="$t('testSuite.compareWithHistorical')"
      @input="changeEnabledField"
    />
    <v-radio-group
      v-if="value"
      v-field="value"
      row
    >
      <v-radio
        v-for="item in preparedItems"
        :key="item.value"
        v-bind="item"
      />
    </v-radio-group>
  </v-layout>
</template>

<script>
import { TEST_SUITE_HISTORICAL_DATA_MONTHS_DEFAULT_ITEMS } from '@/constants';

export default {
  props: {
    value: {
      type: Number,
      default: 0,
    },
    items: {
      type: Array,
      default: () => TEST_SUITE_HISTORICAL_DATA_MONTHS_DEFAULT_ITEMS,
    },
  },
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
      this.$emit('input', Number(value));
    },
  },
};
</script>
