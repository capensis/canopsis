<template>
  <c-radio-group
    :model-value="modelValue"
    :label="label"
    :disabled="disabled"
    class="mt-0"
    color="primary"
    hide-details
    mandatory
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio :label="stopLabel" :value="stopValue" color="primary" />
    <v-radio :label="continueLabel" :value="continueValue" color="primary" />
  </c-radio-group>
</template>

<script>
import { WORKFLOW_TYPES } from '@/constants';

export default {
  props: {
    modelValue: {
      type: Boolean,
      default: WORKFLOW_TYPES.stop,
    },
    label: {
      type: String,
      default: '',
    },
    stopLabel: {
      type: String,
      default: 'Stop',
    },
    continueLabel: {
      type: String,
      default: 'Continue',
    },
    stopValue: {
      type: Boolean,
      default: WORKFLOW_TYPES.stop,
    },
    continueValue: {
      type: Boolean,
      default: WORKFLOW_TYPES.continue,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:modelValue'],
};
</script>
