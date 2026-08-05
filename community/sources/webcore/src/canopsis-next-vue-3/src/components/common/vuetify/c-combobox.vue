<script setup>
import { computed, useAttrs } from 'vue';

import { useVuetifyFieldValidation, validationEmits, validationProps } from './validation';

const props = defineProps(validationProps);
const emit = defineEmits(validationEmits);
const attrs = useAttrs();

const { inputName, inputValue, errorMessages, onBlur } = useVuetifyFieldValidation(props, emit, attrs);

const comboboxValue = computed({
  get: () => (inputValue.value === '' ? null : inputValue.value),
  set: (value) => {
    inputValue.value = value ?? '';
  },
});
</script>

<template>
  <v-combobox
    v-bind="$attrs"
    v-model="comboboxValue"
    :name="inputName"
    :error-messages="errorMessages"
    class="c-field c-combobox"
    @blur="onBlur"
  >
    <template v-if="label" #label>
      <span v-if="required && label" class="text-error"> * </span>

      {{ label }}
    </template>

    <template v-for="(_, slot) in $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="slotProps || {}" />
    </template>
  </v-combobox>
</template>

<style lang="scss" scoped>
.c-combobox {
  :deep(.v-field),
  :deep(.v-field__field),
  :deep(.v-field__input),
  :deep(input) {
    cursor: text;
  }
}
</style>
