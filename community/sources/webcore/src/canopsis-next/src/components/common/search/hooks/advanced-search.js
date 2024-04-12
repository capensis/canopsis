import { last } from 'lodash';
import { computed, ref, set } from 'vue';

import { ADVANCED_SEARCH_ITEM_TYPES, ADVANCED_SEARCH_UNION_FIELDS } from '@/constants';

import { uid } from '@/helpers/uid';
import {
  getNextAdvancedSearchType,
  isFieldType,
  isValueType,
  prepareAdvancedSearchConditions,
  prepareAdvancedSearchFields,
} from '@/helpers/search/advanced-search';

/**
 * Creates reactive computed properties to determine if the provided type is a `value` type or a `field` type.
 *
 * @param {Object} config - Configuration object for the hook.
 * @param {Object} config.type - An object that should contain a `value` property representing the type to check.
 * @returns {Object} An object containing two computed properties:
 * - `isItemTypeValue`: A boolean indicating if the type is a `value` type.
 * - `isItemTypeField`: A boolean indicating if the type is a `field` type.
 */
export const useAdvancedSearchItemType = ({ type }) => {
  const isItemTypeValue = computed(() => isValueType(type.value));
  const isItemTypeField = computed(() => isFieldType(type.value));

  return {
    isItemTypeValue,
    isItemTypeField,
  };
};

/**
 * Provides reactive state and functions for managing an internal search component in Vue.
 *
 * @returns {object} An object containing:
 * - `internalSearch`: a reactive reference to the search input value.
 * - `internalSearchElement`: a reactive reference to the search input element.
 * - `internalSearchElementFocus`: a function to focus the search input element if it exists.
 * - `setInternalSearch`: a function to set the value of the search input.
 * - `clearInternalSearch`: a function to clear the value of the search input.
 */
export const useAdvancedSearchInternalSearch = () => {
  const internalSearch = ref('');
  const internalSearchElement = ref();

  const internalSearchElementFocus = () => internalSearchElement.value?.$el?.focus();
  const setInternalSearch = (value = '') => internalSearch.value = value;
  const clearInternalSearch = () => setInternalSearch('');

  return {
    internalSearch,
    internalSearchElement,

    internalSearchElementFocus,
    setInternalSearch,
    clearInternalSearch,
  };
};

/**
 * Provides reactive state management for navigating and manipulating an active index within an advanced search context.
 *
 * @param {Object} config - Configuration object for the hook.
 * @param {Ref<Array>} config.value - A reactive reference to the array of search items.
 * @param {Function} config.onChange - A callback function that is called whenever the active index changes.
 * Defaults to a no-op function.
 *
 * @returns {Object} An object containing:
 * - `activeIndex`: A ref to the currently active index in the search items array.
 * - `lastItemType`: A computed ref that returns the type of the last item in the search items array.
 * - `activeItem`: A computed ref that returns the current active item based on the active index.
 * - `activeType`: A computed ref that determines the type of the active item, defaulting to the next logical type
 * if the active item is undefined.
 * - `isActiveIndexLast`: A computed ref that checks if the active index is the last index in the array.
 * - `isActiveIndexNew`: A computed ref that checks if the active index points to a new, yet-to-be-added item
 * (one index past the last).
 * - `goToActiveIndex`: A method to set the active index to a specific value and trigger the onChange callback.
 * - `goToNextActiveIndex`: A method to increment the active index.
 * - `goToPrevActiveIndex`: A method to decrement the active index.
 * - `clearActiveIndex`: A method to reset the active index to -1 (none selected).
 */

export const useAdvancedSearchActiveIndex = ({ value, onChange }) => {
  const activeIndex = ref();
  const lastItemType = computed(() => last(value.value)?.type);
  const activeItem = computed(() => value.value?.[activeIndex.value]);
  const activeType = computed(() => {
    if (activeIndex.value === -1) {
      return null;
    }

    return activeItem.value?.type ?? getNextAdvancedSearchType(lastItemType.value);
  });
  const isActiveIndexLast = computed(() => activeIndex.value === value.value.length - 1);
  const isActiveIndexNew = computed(() => activeIndex.value === value.value.length);

  const goToActiveIndex = (index) => {
    if (activeIndex.value === index) {
      return;
    }

    activeIndex.value = index;

    onChange?.();
  };

  const goToNextActiveIndex = () => goToActiveIndex(activeIndex.value + 1);
  const goToPrevActiveIndex = () => goToActiveIndex(activeIndex.value - 1);
  const clearActiveIndex = () => goToActiveIndex(-1);

  return {
    activeIndex,
    lastItemType,
    activeItem,
    activeType,
    isActiveIndexLast,
    isActiveIndexNew,

    goToActiveIndex,
    goToNextActiveIndex,
    goToPrevActiveIndex,
    clearActiveIndex,
  };
};

/**
 * Custom hook for managing an advanced search internal value state.
 *
 * @param {Object} params - The parameters object.
 * @param {Array} params.value - Initial value for the internal state.
 * @returns {Object} An object containing:
 * - `internalValue`: a reactive reference to the internal value array,
 * - `addItemIntoInternalValue`: function to add a new item to the internal value,
 * - `updateItemInInternalValue`: function to update an item in the internal value by index,
 * - `removeItemFromInternalValue`: function to remove an item from the internal value by index.
 */
export const useAdvancedSearchInternalValue = ({ value }) => {
  const internalValue = ref(value);
  const addItemIntoInternalValue = item => internalValue.value.push({ ...item, key: uid() });
  const updateItemInInternalValue = (index, item) => set(internalValue.value, index, item);
  const removeItemFromInternalValue = (index) => {
    internalValue.value = internalValue.value.filter((item, itemIndex) => index !== itemIndex);
  };

  return {
    internalValue,

    addItemIntoInternalValue,
    updateItemInInternalValue,
    removeItemFromInternalValue,
  };
};

/**
 * Provides reactive state and methods to manage the visibility of an advanced search menu.
 *
 * @param {Object} [options={}] - Configuration options for the menu.
 * @param {Function} [options.onChange] - Callback function that is called when the menu's visibility changes.
 * @returns {Object} An object containing:
 * - `isMenuActive` - A Vue ref that holds the boolean state of the menu's visibility.
 * - `changeMenu` - Function to directly set the visibility of the menu.
 * - `focusMenu` - Function to set the menu as active (visible).
 * - `blurMenu` - Function to set the menu as inactive (hidden).
 */
export const useAdvancedSearchMenu = ({ onChange = () => {} } = {}) => {
  const isMenuActive = ref(false);

  const changeMenu = (value) => {
    isMenuActive.value = value;

    onChange(value);
  };

  const focusMenu = () => changeMenu(true);
  const blurMenu = () => changeMenu(false);

  return {
    isMenuActive,

    changeMenu,
    focusMenu,
    blurMenu,
  };
};

/**
 * Custom Vue composition function to manage and filter advanced search items based on the active type and internal
 * search query.
 *
 * @param {Object} params - The parameters object.
 * @param {Ref<Array>} params.fields - A ref to an array of field objects, each containing at least a `value`
 * and `text` property.
 * @param {Ref<Array>} params.conditions - A ref to an array of condition strings.
 * @param {Ref<string>} params.activeType - A ref to the currently active type of item
 * (`field`, `condition`, or `union`).
 * @param {Ref<string>} params.internalSearch - A ref to the internal search string used to filter items.
 * @returns {Object} An object containing reactive properties for field items, condition items, all items based on the
 * active type, and filtered items based on the search query.
 */
export const useAdvancedSearchItems = ({ fields, conditions, activeType, internalSearch }) => {
  const fieldsItems = computed(() => prepareAdvancedSearchFields(fields.value));
  const conditionsItems = computed(() => prepareAdvancedSearchConditions(conditions.value));

  const items = computed(() => {
    switch (activeType.value) {
      case ADVANCED_SEARCH_ITEM_TYPES.field:
        return fieldsItems.value;

      case ADVANCED_SEARCH_ITEM_TYPES.condition:
        return conditionsItems.value;

      case ADVANCED_SEARCH_ITEM_TYPES.union:
        return ADVANCED_SEARCH_UNION_FIELDS;

      default:
        return [];
    }
  });

  const filteredItems = computed(() => {
    if (!internalSearch.value) {
      return items.value;
    }

    const lowerCaseSearch = internalSearch.value.toLocaleLowerCase();

    return items.value.filter(item => item.text.toLocaleLowerCase().indexOf(lowerCaseSearch) >= 0);
  });

  return {
    fieldsItems,
    conditionsItems,
    items,
    filteredItems,
  };
};

/**
 * Provides reactive handling for "not" values in an advanced search context.
 *
 * @param {Object} params - The parameters object.
 * @param {Ref} params.value - A ref object containing the current search values.
 * @param {Ref} params.activeIndex - A ref to the active index in the search values array.
 * @param {Function} params.updateItemInArray - A function to update an item in the search values array.
 *
 * @returns {Object} An object containing:
 * - `notValues`: a ref object mapping indices to their "not" values.
 * - `activeIndexNotValue`: a computed ref that returns the "not" value for the currently active index.
 * - `setNotValue`: a function to set the "not" value for the current active index and update the item in the array.
 */
export const useAdvancedSearchNotValues = ({ value, activeIndex, updateItemInArray }) => {
  const notValues = ref({});
  const activeIndexNotValue = computed(() => notValues.value[activeIndex.value] ?? false);

  const setNotValue = (notValue) => {
    const item = value.value[activeIndex.value];

    set(notValues.value, activeIndex.value, notValue);

    if (item) {
      updateItemInArray(activeIndex.value, {
        ...item,
        not: value,
      });
    }
  };

  return {
    notValues,
    activeIndexNotValue,

    setNotValue,
  };
};
