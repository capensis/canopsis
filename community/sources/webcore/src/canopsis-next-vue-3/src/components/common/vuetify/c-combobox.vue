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

    <template v-for="(_, slot) in $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="getLegacySlotProps({ name: slot, ...(slotProps || {}) })" />
    </template>
  </v-combobox>
</template>

<script setup>
import { computed, ref, useAttrs } from 'vue';

import VariablesList from '@/components/common/text-editor/variables-list.vue';

import { useVuetifyFieldValidation, validationEmits, validationProps } from './validation';

defineOptions({
  inheritAttrs: false,
});

const props = defineProps({
  ...validationProps,
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
const menuOpened = ref(false);
const pendingBlurEvent = ref(null);

const comboboxValue = computed({
  get: () => (inputValue.value === '' ? null : inputValue.value),
  set: (value) => {
    inputValue.value = value ?? '';
  },
});

const comboboxAttrs = computed(() => {
  if (!props.childrenKey) {
    return attrs;
  }

  return {
    ...attrs,
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
