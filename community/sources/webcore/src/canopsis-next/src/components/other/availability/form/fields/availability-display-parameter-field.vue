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

import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
    name: {
      type: String,
      default: 'parameter',
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
        value: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
        text: t('common.uptime'),
      },
      {
        value: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        text: t('common.downtime'),
      },
    ]);

    return {
      types,
    };
  },
};
</script>
