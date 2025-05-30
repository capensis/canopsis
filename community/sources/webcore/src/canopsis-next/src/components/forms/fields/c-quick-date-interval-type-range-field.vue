<template>
  <v-layout class="gap-2">
    <c-quick-date-interval-type-field
      v-field="value.from"
      :ranges="fromQuickRanges"
      :item-disabled="itemFromDisabled"
      :label="$t('common.from')"
      :hide-details="hideDetails"
      :disabled="disabled"
    />
    <c-quick-date-interval-type-field
      v-field="value.to"
      :ranges="toQuickRanges"
      :item-disabled="itemToDisabled"
      :label="$t('common.to')"
      :hide-details="hideDetails"
      :disabled="disabled"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { PATTERN_QUICK_RANGES_WITHOUT_CUSTOM } from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Object],
      required: false,
    },
    fromRanges: {
      type: Array,
      required: false,
    },
    toRanges: {
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

    const fromQuickRanges = computed(() => {
      const ranges = props.fromRanges ?? PATTERN_QUICK_RANGES_WITHOUT_CUSTOM;

      return ranges.map(range => ({
        ...range,
        text: range.text ?? t(`quickRanges.types.${range.value}`),
      }));
    });

    const toQuickRanges = computed(() => {
      const ranges = props.toRanges ?? PATTERN_QUICK_RANGES_WITHOUT_CUSTOM;

      return ranges.map(range => ({
        ...range,
        text: range.text ?? t(`quickRanges.types.${range.value}`),
      }));
    });

    const allQuickRangesIntervals = computed(() => [
      ...fromQuickRanges.value,
      ...toQuickRanges.value,
    ].reduce((acc, range) => {
      acc[range.value] = {
        start: convertStartDateIntervalToTimestamp(range.start),
        stop: convertStopDateIntervalToTimestamp(range.stop),
      };

      return acc;
    }, {}));

    const itemFromDisabled = (item) => {
      if (!props.value.to) {
        return false;
      }
      const { start: itemStart } = allQuickRangesIntervals.value[item.value] ?? {};
      const { start: toStart } = allQuickRangesIntervals.value[props.value.to] ?? {};

      return toStart <= itemStart;
    };

    const itemToDisabled = (item) => {
      if (!props.value.from) {
        return false;
      }

      const { start: fromStart } = allQuickRangesIntervals.value[props.value.from] ?? {};
      const { start: itemStart } = allQuickRangesIntervals.value[item.value] ?? {};

      return fromStart >= itemStart;
    };

    return {
      fromQuickRanges,
      toQuickRanges,
      updateModel,
      itemFromDisabled,
      itemToDisabled,
    };
  },
};
</script>
