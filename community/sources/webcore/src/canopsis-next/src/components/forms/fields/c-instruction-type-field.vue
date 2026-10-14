<template>
  <c-radio-group
    :model-value="modelValue"
    :label="$t('remediation.instruction.type')"
    :disabled="disabled"
    class="mt-0"
    hide-details
    mandatory
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio
      v-for="type in types"
      :key="type.value"
      :label="type.label"
      :value="type.value"
      color="primary"
    />
  </c-radio-group>
</template>

<script>
import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

export default {
  props: {
    modelValue: {
      type: Number,
      default: REMEDIATION_INSTRUCTION_TYPES.manual,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:modelValue'],
  computed: {
    types() {
      return Object.values(REMEDIATION_INSTRUCTION_TYPES).map(value => ({
        value,
        label: this.$t(`remediation.instruction.types.${value}`),
      }));
    },
  },
};
</script>
