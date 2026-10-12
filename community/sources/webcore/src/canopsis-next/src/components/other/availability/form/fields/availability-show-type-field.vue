<template>
  <c-select
    :model-value="modelValue"
    :items="types"
    :name="name"
    :label="label"
    item-title="text"
    item-value="value"
    hide-details
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
    name: {
      type: String,
      default: 'availability',
    },
    label: {
      type: String,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    const types = computed(() => [
      {
        value: AVAILABILITY_SHOW_TYPE.percent,
        text: t('common.percent'),
      },
      {
        value: AVAILABILITY_SHOW_TYPE.duration,
        text: t('common.duration'),
      },
    ]);

    return {
      types,
    };
  },
};
</script>
