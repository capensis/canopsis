<template>
  <c-select
    v-bind="$attrs"
    :model-value="range"
    :items="quickRanges"
    :label="label || $t('quickRanges.title')"
    :hide-details="hideDetails"
    :disabled="disabled"
    :return-object="returnObject"
    :item-props="itemDisabled ? getItemProps : undefined"
    item-title="text"
    item-value="value"
    @update:model-value="updateModel"
  />
</template>

<script>
import { isFunction, isObject } from 'lodash';
import { computed } from 'vue';

import { QUICK_RANGES } from '@/constants';

import { findQuickRangeValue } from '@/helpers/date/date-intervals';

import { useI18n } from '@/hooks/i18n';

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
    itemDisabled: {
      type: Function,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const quickRanges = computed(() => {
      const ranges = props.ranges ?? Object.values(QUICK_RANGES);

      return ranges.map(range => ({
        ...range,
        text: range.text ?? t(`quickRanges.types.${range.value}`),
      }));
    });

    const range = computed(() => {
      if (!isObject(props.value)) {
        return props.value || null;
      }

      const localRange = findQuickRangeValue(props.value.start, props.value.stop, props.ranges);

      return quickRanges.value.find(({ value }) => value === localRange.value);
    });

    const getItemProps = (item) => {
      const rawItem = item?.raw ?? item;

      if (!isFunction(props.itemDisabled)) {
        return {};
      }

      return { disabled: props.itemDisabled(rawItem) };
    };

    const updateModel = (value) => {
      emit('input', value);
      emit('update:value', value);
    };

    return {
      quickRanges,
      range,
      getItemProps,
      updateModel,
    };
  },
};
</script>
