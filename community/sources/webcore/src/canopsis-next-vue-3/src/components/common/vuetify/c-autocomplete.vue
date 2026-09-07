<template>
  <v-autocomplete
    v-bind="fieldAttrsWithoutSearchListener"
    v-model="autocompleteValue"
    :item-props="itemProps ? getItemProps : undefined"
    :name="inputName"
    :error-messages="errorMessages"
    :density="density"
    class="c-field c-autocomplete"
    @blur="onBlur"
    @update:search="onSearchUpdate"
  >
    <template v-if="label" #label>
      <span v-if="required && label" class="text-error"> * </span>

      {{ label }}
    </template>

    <template v-if="!$slots.item" #item="slotProps">
      <v-list-item v-bind="getHighlightedItemProps(slotProps)">
        <v-list-item-title>
          <v-list-item-mask :text="getHighlightedText(slotProps)" :mask="searchQuery" />
        </v-list-item-title>
      </v-list-item>
    </template>

    <template v-if="hasLimitedSelection" #chip="slotProps">
      <limited-field-selection
        :index="slotProps.index"
        :show-count="showCount"
        :total="selectedCount"
      >
        <slot name="chip" v-bind="getLegacySlotProps({ name: 'chip', ...slotProps })">
          <v-chip v-bind="slotProps.props" size="small">
            {{ getSelectionTitle(slotProps) }}
          </v-chip>
        </slot>
      </limited-field-selection>
    </template>

    <template v-if="hasLimitedSelection && isMultiple" #selection="slotProps">
      <limited-field-selection
        :index="slotProps.index"
        :show-count="showCount"
        :total="selectedCount"
      >
        <slot name="selection" v-bind="getLegacySlotProps({ name: 'selection', ...slotProps })">
          <span>{{ getSelectionTitle(slotProps) }}</span>
        </slot>
      </limited-field-selection>
    </template>

    <template v-for="slot in getForwardedSlotNames()" #[slot]="slotProps">
      <slot :name="slot" v-bind="getLegacySlotProps({ name: slot, ...(slotProps || {}) })" />
    </template>
  </v-autocomplete>
</template>

<script setup>
import { computed, useAttrs } from 'vue';

import LimitedFieldSelection from './limited-field-selection.vue';
import { limitedSelectionProps, useLimitedSelection } from './limited-selection';
import { useSearchHighlight } from './search-highlight';
import { useVuetifyFieldValidation, validationEmits, validationProps } from './validation';

const props = defineProps({
  ...validationProps,
  ...limitedSelectionProps,
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
const {
  hasLimitedSelection,
  isMultiple,
  selectedCount,
  getForwardedSlotNames,
  getSelectionTitle,
} = useLimitedSelection(
  props,
  inputValue,
  { hideSingleSelectionSlot: true },
);
const {
  searchQuery,
  onSearchUpdate,
  fieldAttrsWithoutSearchListener,
  getHighlightedText,
  getHighlightedItemProps,
} = useSearchHighlight(attrs);

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
