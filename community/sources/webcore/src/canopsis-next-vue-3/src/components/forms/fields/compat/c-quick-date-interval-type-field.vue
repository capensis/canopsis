<template>
  <c-select
    v-bind="$attrs"
    :model-value="range"
    :items="quickRanges"
    :label="label || $t('quickRanges.title')"
    :hide-details="hideDetails"
    :disabled="disabled"
    :return-object="returnObject"
    item-title="text"
    item-value="value"
    @update:model-value="updateModel"
  />
</template>

<script>
import { isObject } from 'lodash'

import { QUICK_RANGES } from '@/constants'

import { findQuickRangeValue } from '@/helpers/date/date-intervals'

export default {
  inheritAttrs: false,
  emits: ['input', 'update:value'],
  props: {
    value: {
      type: [String, Number, Object],
      required: false,
    },
    ranges: {
      type: Array,
      required: false,
    },
    label: {
      type: String,
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
      const ranges = this.ranges ?? Object.values(QUICK_RANGES)

      return ranges.map((range) => ({
        ...range,
        text: range.text ?? this.$t(`quickRanges.types.${range.value}`),
      }))
    },

    range() {
      if (!isObject(this.value)) {
        return this.value
      }

      const localRange = findQuickRangeValue(this.value.start, this.value.stop, this.ranges)

      return this.quickRanges.find(({ value }) => value === localRange.value)
    },
  },
  methods: {
    updateModel(value) {
      this.$emit('input', value)
      this.$emit('update:value', value)
    },
  },
}
</script>
