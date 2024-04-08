import { setField } from '@/helpers/immutable';

import { useModelField } from './useModelField';

/**
 * Provides functionality to manipulate an array model field and emit events to the parent component.
 * It extends `useModelField` by adding methods specific to array manipulation, such as adding, updating,
 * and removing items from the array.
 *
 * @param {Object} props - The props object of the component, expected to contain the model property that is an array.
 * @param {Function} emit - The function to emit events to the parent component.
 * @returns {{
 *   modelEvent: string,
 *   modelProp: string,
 *   updateModel: (value: Array) => Array,
 *   updateField: (fieldName: string, value: *) => Array,
 *   addItemIntoArray: (value: *) => Array,
 *   updateItemInArray: (index: number, value: *) => Array,
 *   updateFieldInArrayItem: (index: number, fieldName: string, value: *) => Array,
 *   removeItemFromArray: (index: number) => Array
 * }}
 */
export const useArrayModelField = (props, emit) => {
  const { modelEvent, modelProp, updateModel, updateField } = useModelField(props, emit);

  /**
   * Emit event to parent with new array with new item
   *
   * @param {*} value
   * @return {Array}
   */
  const addItemIntoArray = value => updateModel([...props[modelProp], value]);

  /**
   * Emit event to parent with new array with updated array item
   *
   * @param {number} index
   * @param {*} value
   * @return {Array}
   */
  const updateItemInArray = (index, value) => {
    const items = [...props[modelProp]];

    items[index] = value;

    return updateModel(items);
  };

  /**
   * Emit event to parent with new array with updated array item with updated field
   *
   * @param {number} index
   * @param {string} fieldName
   * @param {*} value
   * @return {Array}
   */
  const updateFieldInArrayItem = (index, fieldName, value) => updateItemInArray(
    index,
    setField(props[modelProp][index], fieldName, value),
  );

  /**
   * Emit event to parent with new array without array item
   *
   * @param {number} index
   * @return {Array}
   */
  const removeItemFromArray = index => updateModel(props[modelProp].filter((v, i) => i !== index));

  return {
    modelEvent,
    modelProp,
    updateModel,
    updateField,
    addItemIntoArray,
    updateItemInArray,
    updateFieldInArrayItem,
    removeItemFromArray,
  };
};
