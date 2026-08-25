<template>
  <c-form-block-row :label="t('scenario.output')" :depth="depth">
    <c-payload-textarea-field
      :value="output"
      :name="`${name}.output`"
      :label="t('scenario.output')"
      :variables="outputVariables"
      @update:value="output = $event"
    >
      <template #append>
        <c-help-icon :text="t('scenario.outputHelp')" icon="help" left />
      </template>
    </c-payload-textarea-field>
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
    templateVars: {
      type: Object,
      default: () => ({}),
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

    const outputVariables = computed(() => props.templateVars.output ?? []);

    return {
      t,
      output,
      outputVariables,
    };
  },
};
</script>
