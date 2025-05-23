<template>
  <v-layout class="gap-2">
    <c-quick-date-interval-type-field
      v-field="value.from"
      :items="quickRanges"
      :item-disabled="itemFromDisabled"
      :label="$t('common.from')"
      :hide-details="hideDetails"
      :disabled="disabled"
    />
    <c-quick-date-interval-type-field
      v-field="value.to"
      :items="quickRanges"
      :item-disabled="itemToDisabled"
      :label="$t('common.to')"
      :hide-details="hideDetails"
      :disabled="disabled"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { QUICK_RANGES_WITHOUT_CUSTOM } from '@/constants';

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
      const ranges = props.ranges ?? Object.values(QUICK_RANGES_WITHOUT_CUSTOM);

      return ranges.map(range => ({
        ...range,
        text: t(`quickRanges.types.${range.value}`),
      }));
    });

    const itemFromDisabled = () => {};

    const itemToDisabled = () => {};

    return {
      quickRanges,
      updateModel,
      itemFromDisabled,
      itemToDisabled,
    };
  },
};
</script>
