<script setup>
import { computed, nextTick, ref, useAttrs } from 'vue';

import { useVuetifyFieldValidation, validationEmits, validationProps } from './validation';

const props = defineProps(validationProps);
const emit = defineEmits(validationEmits);
const attrs = useAttrs();

const { inputValue, errorMessages, onBlur } = useVuetifyFieldValidation(props, emit, attrs);
const menuShown = ref(false);
const pendingBlurEvent = ref(null);

// TODO: We add it for label floating solving. Fix in the future
const selectValue = computed({
  get: () => (inputValue.value === '' ? (attrs.multiple ? [] : null) : inputValue.value),
  set: (value) => {
    inputValue.value = value;
  },
});

const onSelectBlur = async (event) => {
  await nextTick();

  if (menuShown.value) {
    pendingBlurEvent.value = event;
    onBlur(event, false);

    return;
  }

  onBlur(event);
};

const onMenuUpdate = (value) => {
  menuShown.value = value;

  if (!value && pendingBlurEvent.value) {
    onBlur(pendingBlurEvent.value);
    pendingBlurEvent.value = null;
  }
};
</script>

<template>
  <v-select
    v-bind="$attrs"
    v-model="selectValue"
    :name="name"
    :error-messages="errorMessages"
    class="c-field c-select"
    @blur="onSelectBlur"
    @update:menu="onMenuUpdate"
  >
    <template v-if="label" #label>
      <span v-if="required && label" class="text-error"> * </span>

      {{ label }}
    </template>

    <template v-for="(_, slot) in $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="slotProps || {}" />
    </template>
  </v-select>
</template>
