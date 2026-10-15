<template>
  <v-combobox
    v-bind="comboboxAttrs"
    v-model="comboboxValue"
    :name="inputName"
    :error-messages="errorMessages"
    :density="density"
    class="c-field c-combobox"
    @blur="handleBlur"
    @update:menu="handleMenuUpdate"
    @update:search="onSearchUpdate"
  >
    <template v-if="label" #label>
      <span v-if="required && label" class="text-error"> * </span>

      {{ label }}
    </template>

    <template v-if="childrenKey" #prepend-item>
      <variables-list
        :value="comboboxValue"
        :items="attrs.items"
        :children-key="childrenKey"
        :return-object="Boolean(attrs.returnObject)"
        @update:value="selectChildrenValue"
      />
    </template>

    <template v-if="childrenKey" #no-data />

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
  </v-combobox>
</template>

<script setup>
import { computed, ref, useAttrs } from 'vue';

import VariablesList from '@/components/common/text-editor/variables-list.vue';

import LimitedFieldSelection from './limited-field-selection.vue';
import { limitedSelectionProps, useLimitedSelection } from './limited-selection';
import { useSearchHighlight } from './search-highlight';
import { useVuetifyFieldValidation, validationEmits, validationProps } from './validation';

defineOptions({
  inheritAttrs: false,
});

const props = defineProps({
  ...validationProps,
  ...limitedSelectionProps,
  childrenKey: {
    type: String,
    default: '',
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
const menuOpened = ref(false);
const pendingBlurEvent = ref(null);

const comboboxValue = computed({
  get: () => (inputValue.value === '' ? (attrs.multiple ? [] : null) : inputValue.value),
  set: (value) => {
    inputValue.value = value ?? '';
  },
});

const comboboxAttrs = computed(() => {
  if (!props.childrenKey) {
    return fieldAttrsWithoutSearchListener.value;
  }

  return {
    ...fieldAttrsWithoutSearchListener.value,
    items: [],
    hideNoData: false,
    noFilter: true,
    menu: menuOpened.value,
    menuProps: {
      minWidth: 200,
      closeOnContentClick: false,
      ...attrs.menuProps,
    },
  };
});

const handleBlur = (event) => {
  if (menuOpened.value) {
    pendingBlurEvent.value = event;

    return;
  }

  onBlur(event);
};

const handleMenuUpdate = (opened) => {
  menuOpened.value = opened;

  if (!opened && pendingBlurEvent.value) {
    onBlur(pendingBlurEvent.value);
    pendingBlurEvent.value = null;
  }
};

const selectChildrenValue = (value) => {
  comboboxValue.value = value;
  menuOpened.value = false;
};

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
</script>

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
