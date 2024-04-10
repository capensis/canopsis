<template>
  <v-menu
    ref="menu"
    v-bind="menuProps"
    :value="value"
    :activator="activator"
    @input="menuInput"
  >
    <v-list>
      <v-list-item v-if="isActiveTypeValue">
        <v-list-item-title v-html="$t('advancedSearch.valueTypeListMessage')" />
      </v-list-item>
      <v-list-item-group v-else :value="activeItem.value">
        <template v-if="isActiveTypeField">
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
    </v-list>
  </v-menu>
</template>

<script>
import { computed } from 'vue';

import { ADVANCED_SEARCH_ITEM_TYPES } from '@/constants';

import { isFieldType, isValueType } from '@/helpers/search/advanced-search';

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
      }),
    },
    activeType: {
      type: String,
      default: ADVANCED_SEARCH_ITEM_TYPES.field,
    },
    activeItem: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const isActiveTypeValue = computed(() => isValueType(props.activeType));
    const isActiveTypeField = computed(() => isFieldType(props.activeType));

    const menuInput = value => emit('input', value);
    const selectItem = item => emit('select:item', item);
    const changeNotSwitcher = value => emit('change:not', value);

    return {
      isActiveTypeValue,
      isActiveTypeField,

      menuInput,
      selectItem,
      changeNotSwitcher,
    };
  },
};
</script>
