<template>
  <v-radio-group
    :model-value="modelValue"
    class="mt-0"
    hide-details
    mandatory
    inline
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio
      :label="$t('common.role')"
      :value="$constants.REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role"
      color="primary"
    />
    <v-radio
      :label="$t('common.user')"
      :value="$constants.REMEDIATION_INSTRUCTION_APPROVAL_TYPES.user"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
import { REMEDIATION_INSTRUCTION_APPROVAL_TYPES } from '@/constants';

export default {
  props: {
    modelValue: {
      type: Number,
      default: REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role,
    },
  },
  emits: ['update:modelValue'],
};
</script>
