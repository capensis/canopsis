<template>
  <c-text-field
    :model-value="modelValue"
    :label="t('remediation.instruction.endpoint')"
    :name="name"
    :disabled="disabled"
    required
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #append>
      <c-help-icon :text="t('remediation.instruction.tooltips.endpoint')" icon="help" left />
    </template>
  </c-text-field>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'endpoint',
    },
  },
  emits: ['update:modelValue'],
  setup() {
    const { t } = useI18n();

    return {
      t,
    };
  },
};
</script>
