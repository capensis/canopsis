<template>
  <c-select
    :model-value="modelValue"
    :items="typeItems"
    :label="$t('common.actionType')"
    name="type"
    item-title="title"
    item-value="value"
    required
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: DYNAMIC_INFO_INFORMATION_TYPES.setToInfo,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    const typeItems = computed(() => Object.values(DYNAMIC_INFO_INFORMATION_TYPES).map(type => ({
        title: t(`dynamicInfo.informationTypes.${type}`),
        value: type,
      })),
    );

    return {
      typeItems,
    };
  },
};
</script>
