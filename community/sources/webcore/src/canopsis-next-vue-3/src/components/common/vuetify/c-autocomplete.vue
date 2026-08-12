<script setup>
import { computed, useAttrs } from 'vue';

import { useVuetifyFieldValidation, validationEmits, validationProps } from './validation';

const props = defineProps({
  ...validationProps,
  itemProps: {
    type: [Function, Object, String, Array],
    default: undefined,
  },
});
const emit = defineEmits(validationEmits);
const attrs = useAttrs();

const { inputName, inputValue, errorMessages, onBlur } = useVuetifyFieldValidation(
  props,
  emit,
  attrs,
);

const autocompleteValue = computed({
  get: () => (inputValue.value === '' ? (attrs.multiple ? [] : null) : inputValue.value),
  set: (value) => {
    inputValue.value = value;
  },
});

const getLegacySlotProps = slotProps => ({
  ...slotProps,
  item: slotProps?.item?.raw ?? slotProps?.item,
  vuetifyItem: slotProps?.item,
});

const getItemProps = (item) => {
  if (typeof props.itemProps === 'function') {
    return props.itemProps(item?.raw ?? item);
  }

  return props.itemProps;
};
</script>

<template>
  <v-autocomplete
    v-bind="$attrs"
    v-model="autocompleteValue"
    :item-props="itemProps ? getItemProps : undefined"
    :name="inputName"
    :error-messages="errorMessages"
    class="c-field c-autocomplete"
    @blur="onBlur"
  >
    <template v-if="label" #label>
      <span v-if="required && label" class="text-error"> * </span>

      {{ label }}
    </template>

    <template v-for="(_, slot) in $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="getLegacySlotProps(slotProps || {})" />
    </template>
  </v-autocomplete>
</template>
