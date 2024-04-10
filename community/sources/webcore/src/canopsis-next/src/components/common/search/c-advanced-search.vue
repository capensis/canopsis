<template>
  <v-layout column>
    <v-layout>
      <div
        :class="{ 'v-input--is-focused': isFocused, 'theme--light': !$system.dark, 'theme--dark': $system.dark }"
        class="v-input v-text-field v-text-field--is-booted v-select v-select--chips primary--text"
      >
        <div class="v-input__control">
          <div ref="activator" class="v-input__slot">
            <div class="v-select__slot">
              <div class="v-select__selections">
                <c-advanced-search-chip
                  v-for="(item, index) in internalValue"
                  :key="item.key"
                  :item="item"
                  :active="activeIndex === index"
                  @select:item="clickItem(index)"
                />
                <input
                  v-model="internalSearch"
                  ref="internalSearchElement"
                  type="text"
                  autocomplete="off"
                  @blur="blurInput"
                  @click="clickInput"
                  @keydown="keydownInput"
                >
              </div>
            </div>
          </div>
          <c-advanced-search-list
            ref="list"
            :value="isMenuActive"
            :activator="activator"
            :not-switcher-value="notValues[activeIndex]"
            :active-item="activeItem"
            :active-type="activeType"
            :items="filteredItems"
            @input="changeMenu"
            @select:item="selectItem"
            @change:not="changeNot"
          />
          <v-messages
            :value="errorMessages"
            color="error"
          />
        </div>
      </div>
    </v-layout>
  </v-layout>
</template>

<script>
import { last } from 'lodash';
import { computed, ref, set, nextTick } from 'vue';

import { ADVANCED_SEARCH_CONDITIONS, ADVANCED_SEARCH_ITEM_TYPES, ADVANCED_SEARCH_UNION_FIELDS } from '@/constants';

import {
  getNextAdvancedSearchType,
  parseAdvancedSearch,
  prepareAdvancedSearchConditions,
  prepareAdvancedSearchFields,
} from '@/helpers/search/advanced-search';

import {
  useType,
  useInternalValue,
  useInternalSearch,
  useActiveIndex,
  useMenu,
} from './hooks/advanced-search';
import CAdvancedSearchList from './partials/c-advanced-search-list.vue';
import CAdvancedSearchChip from './partials/c-advanced-search-chip.vue';

const KEY_CODES = {
  enter: 13,
  tab: 9,
  delete: 46,
  esc: 27,
  space: 32,
  up: 38,
  down: 40,
  left: 37,
  right: 39,
  end: 35,
  home: 36,
  del: 46,
  backspace: 8,
  insert: 45,
  pageup: 33,
  pagedown: 34,
  shift: 16,
};

export default {
  inject: ['$system'],
  components: {
    CAdvancedSearchList,
    CAdvancedSearchChip,
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    fields: {
      type: Array,
      default: () => [],
    },
    conditions: {
      type: Array,
      default: () => Object.values(ADVANCED_SEARCH_CONDITIONS),
    },
  },
  setup(props) {
    const list = ref();
    const activator = ref();
    const internalSearchElement = ref();
    const errorMessages = ref([]);

    let parsedValue = [];

    try {
      parsedValue = parseAdvancedSearch(props.value);
    } catch (err) {
      console.error(err);
      errorMessages.value.push('Incorrect search');
    }

    const {
      internalValue,
      addItemIntoInternalValue,
      updateItemInInternalValue,
      removeItemFromInternalValue,
    } = useInternalValue({ value: parsedValue });

    const {
      internalSearch,
      isInternalSearchEmpty,
      clearInternalSearch,
    } = useInternalSearch();

    const {
      activeIndex,
      gotToActiveIndex,
      goToNextActiveIndex,
      goToPrevActiveIndex,
    } = useActiveIndex({ onChange: clearInternalSearch });

    const lastItemType = computed(() => last(internalValue.value)?.type);
    const activeItem = computed(() => internalValue.value?.[activeIndex.value]);
    const activeType = computed(() => activeItem.value?.type ?? getNextAdvancedSearchType(lastItemType.value));

    const { isTypeValue: isActiveTypeValue } = useType(activeType);

    const items = computed(() => {
      switch (activeType.value) {
        case ADVANCED_SEARCH_ITEM_TYPES.field:
          return prepareAdvancedSearchFields(props.fields);

        case ADVANCED_SEARCH_ITEM_TYPES.condition:
          return prepareAdvancedSearchConditions(props.conditions);

        case ADVANCED_SEARCH_ITEM_TYPES.union:
          return ADVANCED_SEARCH_UNION_FIELDS;

        default:
          return [];
      }
    });

    const filteredItems = computed(() => {
      if (!isActiveTypeValue.value || !internalSearch.value) {
        return items.value;
      }

      const lowerCaseSearch = internalSearch.value.toLocaleLowerCase();

      return items.value.filter(item => item.text.toLocaleLowerCase().indexOf(lowerCaseSearch) >= 0);
    });

    const goNext = () => {
      clearInternalSearch();
      goToNextActiveIndex();

      errorMessages.value = [];
    };

    const addInternalSearchIntoInternalValue = () => {
      if (!isActiveTypeValue) {
        return;
      }

      addItemIntoInternalValue({
        value: internalSearch.value,
        type: ADVANCED_SEARCH_ITEM_TYPES.value,
        text: internalSearch.value,
        isValue: true,
      });
      goNext();
    };

    const {
      isMenuActive,
      isFocused,
      changeMenu,
      focusMenu,
    } = useMenu();

    const isNewItemActiveIndex = computed(() => internalValue.value.length === activeIndex.value);

    const notValues = ref({});
    const activeIndexNotValue = computed(() => notValues.value[activeIndex.value] ?? false);

    const selectItem = (item) => {
      if (isNewItemActiveIndex.value) {
        addItemIntoInternalValue({ ...item, not: activeIndexNotValue.value });
      } else {
        updateItemInInternalValue(activeIndex.value, { ...item, not: activeIndexNotValue.value });
      }

      goNext();
      internalSearchElement.value.focus();
    };

    const clickItem = (index = internalValue.value.length) => {
      focusMenu();
      gotToActiveIndex(index);
    };

    const clickInput = () => {
      clickItem();
    };

    const blurInput = () => {
      if (!isInternalSearchEmpty.value) {
        addInternalSearchIntoInternalValue();
      }
    };

    const keydownInput = (event) => {
      const { menu } = list.value?.$refs ?? {};

      if (!menu) {
        return;
      }

      if (event.keyCode === KEY_CODES.backspace && !internalSearch.value) {
        removeItemFromInternalValue(internalValue.value.length - 1);
        goToPrevActiveIndex();
        return;
      }

      if (
        isMenuActive.value
        && [
          KEY_CODES.enter,
          KEY_CODES.esc,
          KEY_CODES.down,
          KEY_CODES.up,
          KEY_CODES.home,
          KEY_CODES.end,
        ].includes(event.keyCode)
      ) {
        nextTick(() => {
          menu.changeListIndex(event);
        });
      }

      if (event.keyCode === KEY_CODES.esc) {
        clearInternalSearch();
        return;
      }

      if (isActiveTypeValue.value && event.keyCode === KEY_CODES.enter) {
        addInternalSearchIntoInternalValue();
      }
    };

    const changeNot = (value) => {
      const item = internalValue.value[activeIndex.value];

      set(notValues.value, activeIndex.value, value);

      if (item) {
        updateItemInInternalValue(activeIndex.value, { ...item, not: value });
      }
    };

    return {
      list,
      activator,
      internalSearchElement,

      activeItem,
      activeType,
      activeIndex,
      internalValue,
      internalSearch,
      notValues,
      isFocused,
      isMenuActive,
      filteredItems,
      errorMessages,

      selectItem,
      clickItem,
      clickInput,
      blurInput,
      changeMenu,
      keydownInput,
      changeNot,
    };
  },
};
</script>
