<template>
  <c-select
    :model-value="modelValue"
    :items="methodTypes"
    :label="label || t('common.method')"
    :name="name"
    item-title="text"
    item-value="value"
    autofocus
    required
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { JUNIT_STATE_SETTING_METHODS } from '@/constants';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: JUNIT_STATE_SETTING_METHODS.worst,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'method',
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    const methodTypes = computed(() =>
      Object.values(JUNIT_STATE_SETTING_METHODS).map(method => ({
        value: method,
        text: t(`stateSetting.junit.methods.${method}`),
      })),
    );

    return {
      t,
      methodTypes,
    };
  },
};
</script>
