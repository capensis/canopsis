<template>
  <v-select
    v-bind="$attrs"
    :value="range"
    :items="quickRanges"
    :label="label || $t('quickRanges.title')"
    :hide-details="hideDetails"
    :disabled="disabled"
    :return-object="returnObject"
    @input="updateModel"
  />
</template>

<script>
import { isObject } from 'lodash';
import { computed } from 'vue';

import { QUICK_RANGES } from '@/constants';

import { findQuickRangeValue } from '@/helpers/date/date-intervals';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

export default {
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
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const quickRanges = computed(() => {
      const ranges = props.ranges ?? Object.values(QUICK_RANGES);

      return ranges.map(range => ({
        ...range,
        text: t(`quickRanges.types.${range.value}`),
      }));
    });

    const range = computed(() => {
      if (!isObject(props.value)) {
        return props.value;
      }

      const localRange = findQuickRangeValue(props.value.start, props.value.stop, props.ranges);

      return quickRanges.value.find(({ value }) => value === localRange.value);
    });

    return {
      quickRanges,
      range,
      updateModel,
    };
  },
};
</script>
