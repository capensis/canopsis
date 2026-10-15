<template>
  <c-select
    :model-value="modelValue"
    :label="label || t('stateSetting.appliedForEntityType')"
    :items="entityTypeItems"
    :name="name"
    item-title="text"
    item-value="value"
    required
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { STATE_SETTING_ENTITY_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t, te } = useI18n();

    const entityTypeItems = computed(() => STATE_SETTING_ENTITY_TYPES.map(value => ({
      value,
      text: te(`common.${value}`) ? t(`common.${value}`) : value,
    })));

    return {
      t,
      entityTypeItems,
    };
  },
};
</script>
