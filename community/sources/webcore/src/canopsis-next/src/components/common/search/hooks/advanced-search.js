import { computed, ref, unref, set } from 'vue';

import { uid } from '@/helpers/uid';
import { isFieldType, isValueType } from '@/helpers/search/advanced-search';

export const useType = (type) => {
  const isTypeValue = computed(() => isValueType(unref(type)));
  const isTypeField = computed(() => isFieldType(unref(type)));

  return {
    isTypeValue,
    isTypeField,
  };
};

export const useInternalSearch = () => {
  const internalSearch = ref('');
  const isInternalSearchEmpty = computed(() => !internalSearch.value);

  const setInternalSearch = (value = '') => {
    internalSearch.value = value;
  };
  const clearInternalSearch = () => setInternalSearch('');

  return {
    internalSearch,
    isInternalSearchEmpty,

    setInternalSearch,
    clearInternalSearch,
  };
};

export const useActiveIndex = ({ onChange }) => {
  const activeIndex = ref();

  const gotToActiveIndex = (value) => {
    if (value < 0) {
      return;
    }

    if (activeIndex.value === value) {
      return;
    }

    activeIndex.value = value;

    onChange();
  };

  const goToNextActiveIndex = () => gotToActiveIndex(activeIndex.value + 1);

  const goToPrevActiveIndex = () => gotToActiveIndex(activeIndex.value - 1);

  return {
    activeIndex,

    gotToActiveIndex,
    goToNextActiveIndex,
    goToPrevActiveIndex,
  };
};

export const useInternalValue = ({ value }) => {
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

export const useMenu = ({ onChange = () => {} } = {}) => {
  const isMenuActive = ref(false);
  const isFocused = ref(false);

  const changeMenu = (value) => {
    if (isFocused.value !== value) {
      isFocused.value = value;
    }

    isMenuActive.value = value;

    onChange(value);
  };

  const focusMenu = () => changeMenu(true);
  const blurMenu = () => changeMenu(false);

  return {
    isMenuActive,
    isFocused,

    changeMenu,
    focusMenu,
    blurMenu,
  };
};
