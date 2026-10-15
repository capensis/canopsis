<template>
  <v-layout class="gap-2">
    <v-flex xs6>
      <c-quick-date-interval-type-field
        v-field="value.from"
        :ranges="fromQuickRanges"
        :item-disabled="itemFromDisabled"
        :label="$t('common.from')"
        :hide-details="hideDetails"
        :disabled="disabled"
      />
    </v-flex>
    <v-flex xs6>
      <c-quick-date-interval-type-field
        v-field="value.to"
        :ranges="toQuickRanges"
        :item-disabled="itemToDisabled"
        :label="$t('common.to')"
        :hide-details="hideDetails"
        :disabled="disabled"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { toRef } from 'vue';

import { useModelField } from '@/hooks/form/model-field';
import { useQuickDateIntervalRange } from '@/hooks/form/quick-date-interval-range';

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
    const { updateModel } = useModelField(props, emit);

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
      updateModel,
      itemFromDisabled,
      itemToDisabled,
    };
  },
};
</script>
