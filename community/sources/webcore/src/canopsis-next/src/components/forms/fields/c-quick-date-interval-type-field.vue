<template lang="pug">
  v-select(
    :value="range",
    :items="quickRanges",
    :label="$t('quickRanges.title')",
    :hide-details="hideDetails",
    return-object,
    @input="updateModel($event)"
  )
</template>

<script>
import { QUICK_RANGES } from '@/constants';

import { findQuickRangeValue } from '@/helpers/date/date-intervals';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    customFilter: {
      type: Function,
      default: () => true,
    },
    hideDetails: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    range() {
      const range = findQuickRangeValue(this.value.start, this.value.stop);

      return this.quickRanges.find(({ value }) => value === range.value);
    },

    quickRanges() {
      return Object.values(QUICK_RANGES)
        .filter(this.customFilter)
        .map(range => ({
          ...range,
          text: this.$t(`quickRanges.types.${range.value}`),
        }));
    },
  },
};
</script>
