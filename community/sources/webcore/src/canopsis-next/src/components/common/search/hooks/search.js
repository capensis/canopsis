import { ref } from 'vue';

import { prepareSearchForSubmit } from '@/helpers/search/search';

/**
 * Custom Vue composition function to manage local search value and handle submission.
 *
 * @param {Object} params - The parameters object.
 * @param {Object} params.value - A ref object containing the initial value of the search input.
 * @param {Object} params.columns - A ref object containing the columns to be used in the search.
 * @param {Function} params.onSubmit - The callback function to execute on search submission.
 * @param {Function} emit - The Vue emit function used to trigger events.
 * @returns {Object} Returns an object containing the local search value ref and functions to submit or clear the
 * search.
 */
export const useSearchLocalValue = ({ value, columns, onSubmit }, emit) => {
  const localValue = ref(value.value);

  const submit = () => {
    emit('submit', prepareSearchForSubmit(localValue.value, columns.value));
    onSubmit(localValue.value);
  };

  const clear = () => {
    localValue.value = '';

    emit('submit', '');
  };

  return {
    localValue,

    submit,
    clear,
  };
};

/**
 * Provides methods to manipulate saved search items by emitting events.
 *
 * @param {Function} emit - The function to call to emit events.
 * @returns {Object} An object containing methods to add, remove, and toggle pinning of search items.
 */
export const useSearchSavedItems = (emit) => {
  const addItem = search => emit('add:item', search);
  const removeItem = search => emit('remove:item', search);
  const togglePinItem = search => emit('toggle-pin:item', search);

  return {
    addItem,
    removeItem,
    togglePinItem,
  };
};
