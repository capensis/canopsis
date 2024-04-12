<template>
  <v-layout align-end>
    <c-search-field
      v-model="localValue"
      :combobox="combobox"
      :items="items"
      @submit="submit"
      @remove:item="removeItem"
      @toggle-pin:item="togglePinItem"
    >
      <slot />
    </c-search-field>
    <c-action-btn
      :tooltip="$t('common.search')"
      icon="search"
      @click="submit"
    />
    <c-action-btn
      :tooltip="$t('common.clearSearch')"
      icon="clear"
      @click="clear"
    />
    <slot />
  </v-layout>
</template>

<script>
import { toRef } from 'vue';

import { useSearchLocalValue, useSearchSavedItems } from './hooks/search';
import CSearchField from './c-search-field.vue';

/**
 * Search component
 */
export default {
  components: { CSearchField },
  props: {
    value: {
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
    columns: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { addItem, removeItem, togglePinItem } = useSearchSavedItems(emit);

    const {
      localValue,
      submit,
      clear,
    } = useSearchLocalValue({
      value: toRef(props, 'value'),
      columns: toRef(props, 'columns'),
      onSubmit: addItem,
    }, emit);

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
