<script setup>
import { useAttrs } from 'vue';

import { useVuetifyFieldValidation, validationEmits, validationProps } from './validation';

const props = defineProps(validationProps);
const emit = defineEmits(validationEmits);
const attrs = useAttrs();

const { inputName, inputValue, errorMessages, onBlur } = useVuetifyFieldValidation(props, emit, attrs);
</script>

<template>
  <v-radio-group
    v-bind="$attrs"
    v-model="inputValue"
    :color="$attrs.color || 'primary'"
    :name="inputName"
    :error-messages="errorMessages"
    @blur="onBlur"
  >
    <template v-if="label" #label>
      <span v-if="required && label" class="text-error"> * </span>

      {{ label }}
    </template>

    <template v-for="(_, slot) in $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="slotProps || {}" />
    </template>
  </v-radio-group>
</template>
