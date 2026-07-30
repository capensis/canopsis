<template>
  <c-select
    :model-value="modelValue"
    :items="preparedItems"
    :label="$t('llm.thinkingLevel')"
    :name="name"
    :disabled="disabled || !preparedItems.length"
    item-title="text"
    item-value="value"
    clearable
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed, watch } from 'vue';

import { LLM_THINKING_LEVELS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'thinking_level',
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

  setup(props, { emit }) {
    const { t } = useI18n();

    const preparedItems = computed(() =>
      Object.values(props.items).map((value) => ({
        value,
        text: t(`llm.thinkingLevels.${value}`),
      })),
    );

    watch(
      () => props.items,
      (newItems) => {
        if (newItems.length && !props.modelValue) {
          const middleIndex = Math.ceil(newItems.length / 2);

          emit(
            'update:modelValue',
            newItems.includes(LLM_THINKING_LEVELS.medium)
              ? LLM_THINKING_LEVELS.medium
              : newItems[middleIndex],
          );
        } else if (!newItems.includes(props.modelValue)) {
          emit('update:modelValue', null);
        }
      },
    );

    return {
      preparedItems,
    };
  },
};
</script>
