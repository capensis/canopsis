<template>
  <v-select
    v-bind="$attrs"
    :value="range"
    :items="quickRanges"
    :label="$t('quickRanges.title')"
    :hide-details="hideDetails"
    :disabled="disabled"
    :return-object="returnObject"
    @input="updateModel($event)"
  />
</template>

<script>
import { isObject } from 'lodash';

import { QUICK_RANGES } from '@/constants';

import { findQuickRangeValue } from '@/helpers/date/date-intervals';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Object],
      required: false,
    },
    ranges: {
      type: Array,
      required: false,
    },
    hideDetails: {
      type: Boolean,
      required: false,
    },
    disabled: {
      type: Boolean,
      required: false,
    },
    returnObject: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    quickRanges() {
      const ranges = this.ranges ?? Object.values(QUICK_RANGES);

      return ranges.map(range => ({
        ...range,
        text: this.$t(`quickRanges.types.${range.value}`),
      }));
    },

    range() {
      if (!isObject(this.value)) {
        return this.value;
      }

      const range = findQuickRangeValue(this.value.start, this.value.stop, this.ranges);

      return this.quickRanges.find(({ value }) => value === range.value);
    },
  },
};
</script>
