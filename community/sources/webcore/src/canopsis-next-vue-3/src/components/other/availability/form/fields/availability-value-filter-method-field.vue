<template>
  <c-select
    :model-value="modelValue"
    :label="$t('availability.filterByValue')"
    :items="valueFilterMethods"
    item-title="text"
    item-value="value"
    hide-details
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    const valueFilterMethods = computed(() =>
      Object.values(AVAILABILITY_VALUE_FILTER_METHODS).map(valueFilterMethod => ({
        value: valueFilterMethod,
        text: t(`availability.valueFilterMethods.${valueFilterMethod}`),
      })),
    );

    return {
      valueFilterMethods,
    };
  },
};
</script>
