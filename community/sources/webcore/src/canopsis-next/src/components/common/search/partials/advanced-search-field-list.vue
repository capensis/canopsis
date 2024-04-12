<template>
  <v-menu
    ref="menuElement"
    v-bind="menuProps"
    :value="value"
    :activator="activator"
    @input="menuInput"
  >
    <v-list>
      <v-list-item v-if="isItemTypeValue">
        <v-list-item-title v-html="listMessage" />
      </v-list-item>
      <v-list-item-group v-else-if="items.length" :value="activeValue">
        <template v-if="isItemTypeField">
          <v-switch
            key="not"
            :input-value="notSwitcherValue"
            :label="$t('advancedSearch.not')"
            class="mx-3"
            color="primary"
            @change="changeNotSwitcher"
          />
          <v-divider key="divider" />
        </template>
        <v-list-item
          v-for="item in items"
          :key="item.value"
          :aria-selected="item.value"
          :value="item.value"
          @click="selectItem(item)"
        >
          <v-list-item-title>{{ item.text }}</v-list-item-title>
        </v-list-item>
      </v-list-item-group>
      <v-list-item v-else>
        <v-list-item-title class="grey--text">
          {{ $t('advancedSearch.noDataList') }}
        </v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { toRef } from 'vue';

import { ADVANCED_SEARCH_ITEM_TYPES } from '@/constants';

import { useAdvancedSearchItemType } from '../hooks/advanced-search';

export default {
  inject: ['$system'],
  props: {
    activator: {
      default: null,
      validator: val => ['string', 'object'].includes(typeof val),
    },
    value: {
      type: Boolean,
      default: false,
    },
    notSwitcherValue: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    menuProps: {
      type: Object,
      default: () => ({
        openOnClick: false,
        disableKeys: true,
        closeOnContentClick: false,
        ignoreClickOutsideOnActivator: true,
        maxHeight: 304,
        nudgeBottom: 1,
        bottom: true,
        offsetY: true,
        transition: false,
      }),
    },
    activeType: {
      type: String,
      default: ADVANCED_SEARCH_ITEM_TYPES.field,
    },
    activeValue: {
      type: String,
      required: false,
    },
    listMessage: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const {
      isItemTypeField,
      isItemTypeValue,
    } = useAdvancedSearchItemType({ type: toRef(props, 'activeType') });

    const menuInput = value => emit('input', value);
    const selectItem = item => emit('select:item', item);
    const changeNotSwitcher = value => emit('change:not', value);

    return {
      isItemTypeField,
      isItemTypeValue,

      menuInput,
      selectItem,
      changeNotSwitcher,
    };
  },
};
</script>
