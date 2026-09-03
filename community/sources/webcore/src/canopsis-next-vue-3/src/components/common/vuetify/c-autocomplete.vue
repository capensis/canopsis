<template>
  <v-autocomplete
    v-bind="$attrs"
    v-model="autocompleteValue"
    :item-props="itemProps ? getItemProps : undefined"
    :name="inputName"
    :error-messages="errorMessages"
    :density="density"
    class="c-field c-autocomplete"
    @blur="onBlur"
  >
    <template v-if="label" #label>
      <span v-if="required && label" class="text-error"> * </span>

      {{ label }}
    </template>

    <template v-for="(_, slot) in $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="getLegacySlotProps({ name: slot, ...(slotProps || {}) })" />
    </template>
  </v-autocomplete>
</template>

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
  props: slotProps.props && slotProps.name === 'item'
    ? {
      ...slotProps.props,
      title: undefined,
    }
    : slotProps.props,
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

<style lang="scss">
.c-autocomplete {
  &.v-autocomplete--chips .v-field.v-field,
  .v-field.v-field:has(.v-autocomplete__selection .v-chip) {
    height: auto;
    min-height: 48px;
  }

  &.v-autocomplete--chips .v-field,
  .v-field:has(.v-autocomplete__selection .v-chip) {
    .v-field__field,
    .v-field__input,
    .v-field__append-inner,
    .v-field__prepend-inner {
      height: auto;
      min-height: 48px;
    }

    .v-field__input {
      align-items: flex-start;
      flex-wrap: wrap;
      padding-top: 18px;
      padding-bottom: 4px;
    }
  }

  .v-autocomplete__selection {
    max-width: 100%;
    min-width: 0;
  }

  .v-chip {
    max-width: 100%;
  }

  .v-chip__content {
    display: block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}
</style>
