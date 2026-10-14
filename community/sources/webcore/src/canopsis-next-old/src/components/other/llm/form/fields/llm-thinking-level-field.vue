<template>
  <c-select-field
    v-field="value"
    :items="preparedItems"
    :label="$t('llm.thinkingLevel')"
    :name="name"
    :disabled="!preparedItems.length"
    item-text="text"
    item-value="value"
    clearable
  />
</template>

<script>
import { computed, watch } from 'vue';

import { LLM_THINKING_LEVELS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const preparedItems = computed(() => Object.values(props.items).map(value => ({
      value,
      text: t(`llm.thinkingLevels.${value}`),
    })));

    watch(() => props.items, (newItems) => {
      if (newItems.length && !props.value) {
        const middleIndex = Math.ceil(newItems.length / 2);

        updateModel(newItems.includes(LLM_THINKING_LEVELS.medium) ? LLM_THINKING_LEVELS.medium : newItems[middleIndex]);
      } else if (!newItems.includes(props.value)) {
        updateModel(null);
      }
    });

    return {
      preparedItems,
    };
  },
};
</script>
