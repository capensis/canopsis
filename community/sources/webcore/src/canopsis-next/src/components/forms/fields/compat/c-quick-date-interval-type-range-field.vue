<template>
  <v-row density="compact">
    <v-col cols="6">
      <c-quick-date-interval-type-field
        :value="value.from"
        :ranges="fromQuickRanges"
        :item-disabled="itemFromDisabled"
        :label="$t('common.from')"
        :hide-details="hideDetails"
        :disabled="disabled"
        @update:value="updateField('from', $event)"
      />
    </v-col>
    <v-col cols="6">
      <c-quick-date-interval-type-field
        :value="value.to"
        :ranges="toQuickRanges"
        :item-disabled="itemToDisabled"
        :label="$t('common.to')"
        :hide-details="hideDetails"
        :disabled="disabled"
        @update:value="updateField('to', $event)"
      />
    </v-col>
  </v-row>
</template>

<script>
import { toRef } from 'vue';

import { useModelField } from '@/hooks/form/model-field';
import { useQuickDateIntervalRange } from '@/hooks/form/quick-date-interval-range';

export default {
  emits: ['update:value'],
  props: {
    value: {
      type: Object,
      required: true,
    },
    fromRanges: {
      type: Array,
      required: false,
    },
    toRanges: {
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
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const {
      fromQuickRanges,
      toQuickRanges,
      itemFromDisabled,
      itemToDisabled,
    } = useQuickDateIntervalRange({
      fromRanges: toRef(props, 'fromRanges'),
      toRanges: toRef(props, 'toRanges'),
      value: toRef(props, 'value'),
    });

    return {
      fromQuickRanges,
      toQuickRanges,
      itemFromDisabled,
      itemToDisabled,
      updateField,
    };
  },
};
</script>
