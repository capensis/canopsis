<template>
  <v-text-field
    v-bind="$attrs"
    v-model="inputValue"
    :name="inputName"
    :error-messages="errorMessages"
    :density="density"
    :min="min"
    :max="max"
    :class="{ 'c-text-field--hide-label': hideLabel }"
    class="c-field c-text-field"
    @blur="onBlur"
  >
    <template v-if="label && !hideLabel" #label>
      <span v-if="required && label" class="text-error"> * </span>

      {{ label }}
    </template>

    <template v-for="(_, slot) in $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="slotProps || {}" />
    </template>
  </v-text-field>
</template>

<script setup>
import { useAttrs } from 'vue';

import { useVuetifyFieldValidation, validationEmits, validationProps } from './validation';

const props = defineProps({
  ...validationProps,
  hideLabel: {
    type: Boolean,
    default: false,
  },
});
const emit = defineEmits(validationEmits);
const attrs = useAttrs();

defineOptions({
  inheritAttrs: false,
});

const { inputName, inputValue, errorMessages, onBlur } = useVuetifyFieldValidation(
  props,
  emit,
  attrs,
);
defineExpose({
  errorMessages,
});
</script>

<style lang="scss" scoped>
.c-text-field--hide-label {
  :deep(.v-field-label) {
    display: none;
  }

  :deep(.v-field) {
    --v-field-padding-top: 0;
  }
}
</style>
