<script setup>
import { useAttrs } from 'vue';

import {
  useVuetifyFieldValidation,
  validationEmits,
  validationProps,
} from './validation';

const props = defineProps(validationProps);
const emit = defineEmits(validationEmits);
const attrs = useAttrs();

const { inputValue, errorMessages, onBlur } = useVuetifyFieldValidation(props, emit, attrs);
</script>

<template>
  <v-textarea
    v-bind="$attrs"
    v-model="inputValue"
    :name="name"
    :error-messages="errorMessages"
    @blur="onBlur"
  >
    <template #label>
      <span v-if="required" class="text-error">
        *
      </span>

      {{ label }}
    </template>

    <template
      v-for="(_, slot) in $slots"
      #[slot]="slotProps"
    >
      <slot
        :name="slot"
        v-bind="slotProps || {}"
      />
    </template>
  </v-textarea>
</template>
