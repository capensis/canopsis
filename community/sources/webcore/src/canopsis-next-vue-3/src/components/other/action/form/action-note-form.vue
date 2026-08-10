<template>
  <c-form-block-row :label="t('common.output')" :depth="depth">
    <c-description-field v-model="output" :name="`${name}.output`" :label="t('common.output')" />
  </c-form-block-row>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'parameters',
    },
    depth: {
      type: Number,
      default: 0,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const updateField = (field, value) => {
      emit('update:modelValue', {
        ...props.modelValue,
        [field]: value,
      });
    };

    const output = computed({
      get: () => props.modelValue.output ?? '',
      set: value => updateField('output', value),
    });

    return {
      t,
      output,
    };
  },
};
</script>
