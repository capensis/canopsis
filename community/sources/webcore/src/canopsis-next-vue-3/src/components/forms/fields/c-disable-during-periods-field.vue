<template>
  <c-select
    v-bind="$attrs"
    :model-value="modelValue"
    :items="items"
    :label="label || t('common.disableDuringPeriods')"
    :name="name"
    item-title="text"
    item-value="value"
    multiple
    @update:model-value="updateModel"
  />
</template>

<script>
import { computed } from 'vue';

import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

const PERIOD_TYPES = [
  PBEHAVIOR_TYPE_TYPES.maintenance,
  PBEHAVIOR_TYPE_TYPES.pause,
  PBEHAVIOR_TYPE_TYPES.inactive,
];

export default {
  inheritAttrs: false,

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'disable_during_periods',
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const items = computed(() =>
      PERIOD_TYPES.map(value => ({
        value,
        text: t(`pbehavior.types.${value}`),
      })),
    );

    const updateModel = (value) => {
      emit('update:modelValue', value);
    };

    return {
      t,
      items,
      updateModel,
    };
  },
};
</script>
