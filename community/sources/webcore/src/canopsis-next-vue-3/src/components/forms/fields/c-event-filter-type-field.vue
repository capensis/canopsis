<template>
  <c-select
    v-bind="$attrs"
    :model-value="currentValue"
    :items="eventFilterTypes"
    :label="label || t('common.type')"
    :name="name"
    :autofocus="autofocus"
    :required="required"
    item-title="text"
    item-value="value"
    @update:model-value="updateModel"
  />
</template>

<script>
import { computed } from 'vue';

import { EVENT_FILTER_TYPES } from '@/constants';

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
    autofocus: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: true,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const currentValue = computed(() => props.modelValue ?? props.value);

    const eventFilterTypes = computed(() =>
      Object.values(EVENT_FILTER_TYPES).map(type => ({
        value: type,
        text: t(`eventFilter.types.${type}`),
      })),
    );

    const updateModel = value => emit('update:modelValue', value);

    return {
      t,
      currentValue,
      eventFilterTypes,
      updateModel,
    };
  },
};
</script>
