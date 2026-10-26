<template>
  <c-radio-group
    :model-value="modelValue"
    :label="$t('externalData.onMismatch.title')"
    :disabled="disabled"
    column
    hide-details
    mandatory
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio
      v-for="type in types"
      :key="type"
      :label="$t(`externalData.onMismatch.types.${type}`)"
      :value="type"
    />
  </c-radio-group>
</template>

<script>
import { EXTERNAL_DATA_ON_MISMATCH_TYPES } from '@/constants';

export default {
  emits: ['update:modelValue'],
  props: {
    modelValue: {
      type: Number,
      default: EXTERNAL_DATA_ON_MISMATCH_TYPES.throwError,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const types = Object.values(EXTERNAL_DATA_ON_MISMATCH_TYPES);

    return {
      types,
    };
  },
};
</script>
