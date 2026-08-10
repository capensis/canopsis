<template>
  <c-layout class="w-100 px-4" align="end">
    <c-search-field
      v-model="localValue"
      :combobox="combobox"
      :items="items"
      :label="label"
      @submit="submit"
      @remove:item="removeItem"
      @toggle-pin:item="togglePinItem"
    >
      <slot />
    </c-search-field>
    <c-action-btn :tooltip="$t('common.search')" icon="search" @click="submit" />
    <c-action-btn :tooltip="$t('common.clearSearch')" icon="clear" @click="clear" />
    <slot />
  </c-layout>
</template>

<script>
import { toRef } from 'vue';

import { useSearchLocalValue, useSearchSavedItems } from './hooks/search';

/**
 * Search component
 */
export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    combobox: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    fields: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue', 'submit', 'add:item', 'remove:item', 'toggle-pin:item'],

  setup(props, { emit }) {
    const { addItem, removeItem, togglePinItem } = useSearchSavedItems(emit);

    const { localValue, submit, clear } = useSearchLocalValue(
      {
        value: toRef(props, 'modelValue'),
        fields: toRef(props, 'fields'),
        onSubmit: addItem,
      },
      emit,
    );

    return {
      localValue,

      submit,
      clear,
      removeItem,
      togglePinItem,
    };
  },
};
</script>
