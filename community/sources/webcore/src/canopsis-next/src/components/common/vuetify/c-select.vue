<template>
  <v-select
    v-bind="$attrs"
    v-model="selectValue"
    :items="selectItems"
    :item-props="itemProps ? getItemProps : undefined"
    :menu-props="selectMenuProps"
    :name="inputName"
    :error-messages="errorMessages"
    :density="density"
    class="c-field c-select"
    @blur="onSelectBlur"
    @update:menu="onMenuUpdate"
  >
    <template v-if="label" #label>
      <span v-if="required && label" class="text-error"> * </span>

      {{ label }}
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

    <template v-if="hasLimitedSelection" #selection="slotProps">
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
  </v-select>
</template>

<script setup>
import { isString } from 'lodash';
import { computed, nextTick, ref, useAttrs } from 'vue';

import LimitedFieldSelection from './limited-field-selection.vue';
import { limitedSelectionProps, useLimitedSelection } from './limited-selection';
import { useVuetifyFieldValidation, validationEmits, validationProps } from './validation';

const props = defineProps({
  ...validationProps,
  ...limitedSelectionProps,
  items: {
    type: Array,
    default: () => [],
  },
  menuProps: {
    type: Object,
    default: () => ({}),
  },
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
const { hasLimitedSelection, selectedCount, getForwardedSlotNames, getSelectionTitle } = useLimitedSelection(
  props,
  inputValue,
);
const menuShown = ref(false);
const pendingBlurEvent = ref(null);

const normalizeLegacyListItem = (item, index) => {
  if (!item || typeof item !== 'object') {
    return item;
  }

  if (item?.divider) {
    return {
      ...item,
      type: 'divider',
      title: '',
      text: '',
      value: item.value ?? `__divider__${index}`,
    };
  }

  if (item?.header) {
    const title = isString(item.header)
      ? item.header
      : item.title ?? item.text ?? '';

    return {
      ...item,
      type: 'subheader',
      title,
      text: title,
      value: item.value ?? `__subheader__${index}`,
    };
  }

  return item;
};

const hasLegacyHeaders = computed(() => props.items.some(item => !!item?.header),
);
const selectItems = computed(() => props.items.map(normalizeLegacyListItem));
const selectMenuProps = computed(() => {
  if (!hasLegacyHeaders.value) {
    return props.menuProps;
  }

  const contentClass = [
    props.menuProps.contentClass,
    'c-select-menu--with-headers',
  ].filter(Boolean).join(' ');

  return {
    ...props.menuProps,
    contentClass,
  };
});

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
    return props.itemProps(item);
  }

  return props.itemProps;
};
</script>

<style lang="scss">
.c-select-menu--with-headers {
  .v-list-subheader {
    font-weight: 700;
    min-height: 32px;
  }

  .v-list-item {
    padding-left: 32px;
  }
}
</style>
