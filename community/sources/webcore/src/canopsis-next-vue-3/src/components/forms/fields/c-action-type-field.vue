<template>
  <c-select
    v-bind="$attrs"
    :model-value="currentValue"
    :items="items"
    :label="label || t('common.type')"
    :name="name"
    item-title="text"
    item-value="value"
    @update:model-value="updateModel"
  />
</template>

<script>
import { computed } from 'vue';

import { ACTION_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  inheritAttrs: false,

  props: {
    modelValue: {
      type: String,
      default: undefined,
    },
    value: {
      type: String,
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
    types: {
      type: Array,
      default: () => Object.values(ACTION_TYPES),
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const currentValue = computed(() => props.modelValue ?? props.value);

    const items = computed(() => props.types.map(type => ({
      value: type,
      text: t(`scenario.actions.${type}`),
    })));

    const updateModel = (value) => {
      emit('update:modelValue', value);
    };

    return {
      t,
      currentValue,
      items,
      updateModel,
    };
  },
};
</script>
