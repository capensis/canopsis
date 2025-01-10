<template>
  <v-list class="pa-0">
    <slot name="append" />
    <template v-for="item in items">
      <v-subheader
        v-if="item.header"
        :key="item.header"
      >
        {{ item.header }}
      </v-subheader>
      <v-list-item
        v-else
        :key="item.value"
        :input-value="isActiveItem(item)"
        @click="selectVariable(item)"
      >
        <v-list-item-content>
          <v-list-item-title>
            <v-layout class="gap-4" justify-space-between>
              <span v-if="item.text">{{ item.text }}</span>
              <span
                v-if="showValue"
                class="grey--text lighten-1"
              >
                {{ item.value }}
              </span>
            </v-layout>
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </template>
    <slot v-if="!items.length" name="no-data" />
  </v-list>
</template>
<script>
import { find } from 'lodash';
import { toRef, watch, ref } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';

export default {
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    search: {
      type: String,
      default: '',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    itemText: {
      type: String,
      default: 'text',
    },
    fetch: {
      type: Function,
      default: () => {},
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    showValue: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const subItemsShown = ref(false);
    const subItemsPosition = ref({ x: 0, y: 0 });

    const {
      selectedItems,
      items,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      removeItemFromSelectedItemsByIndex,
      updateSearch,
      wholePending: pending,
    } = useLazySearch({
      value: toRef(props, 'value'),
      addable: true,
      idKey: '_id',
      idParamsKey: 'ids',
      fetchHandler: props.fetch,
      multiple: props.multiple,
    }, emit);

    return {
      subItemsShown: false,
      parentItem: undefined,
      subItemsPosition: {
        x: 0,
        y: 0,
      },
    };

    const isActiveItem = (item = {}) => find(selectedItems, { value: item.value });
    const selectVariable = () => {};

    watch(() => props.search, newSearch => updateSearch(newSearch));

    return {
      selectedItems,
      items,
      pending,

      isActiveItem,
      selectVariable,
    };
  },
};
</script>
