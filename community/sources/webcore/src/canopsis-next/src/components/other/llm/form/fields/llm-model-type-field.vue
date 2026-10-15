<template>
  <c-select
    :model-value="modelValue"
    :items="typeItems"
    :label="$t('llm.modelType')"
    :name="name"
    :disabled="disabled"
    :required="required"
    item-title="text"
    item-value="value"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { LLM_MODEL_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
    required: {
      type: Boolean,
      default: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    const typeItems = computed(() => Object.values(LLM_MODEL_TYPES).map(value => ({
        value,
        text: t(`llm.modelTypes.${value}`),
      })),
    );

    return {
      typeItems,
    };
  },
};
</script>
