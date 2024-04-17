import { isArray, isObject, difference } from 'lodash';
import Vue from 'vue';

/**
 * Check if child component has parent component in the $parent
 *
 * @param {VueComponent} child
 * @param {VueComponent} parent
 * @returns {boolean}
 */
export function isParent(child, parent) {
  if (child) {
    if (child === parent || child._original === parent || parent.$el?.contains(child?.$el)) {
      return true;
    }

    if (child.$parent) {
      return isParent(child.$parent, parent);
    }
  }

  return false;
}

/**
 * Merges new data into old data, updating only the properties that have changed.
 * This function is designed to work with Vue's reactivity system, ensuring updates
 * are reactive. It can handle both arrays and objects, recursively merging nested
 * structures. For arrays, it replaces the old array if the length differs, or
 * merges each element. For objects, it updates existing properties, removes
 * properties not present in the new data, and adds any new properties.
 *
 * @param {Object|Array} oldData - The original data to merge into. This data is modified in place.
 * @param {Object|Array} newData - The new data to merge from. This data is not modified.
 * @returns {Object|Array} The updated oldData with changes from newData merged in.
 *
 * @example
 * // For objects
 * const oldObj = { a: 1, b: 2, c: 3 };
 * const newObj = { b: 20, c: 30, d: 40 };
 * mergeReactiveChangedProperties(oldObj, newObj);
 * console.log(oldObj); // Output: { b: 20, c: 30, d: 40 }
 *
 * @example
 * // For arrays
 * const oldArr = [{ a: 1 }, { b: 2 }];
 * const newArr = [{ a: 10 }, { b: 20 }];
 * mergeReactiveChangedProperties(oldArr, newArr);
 * console.log(oldArr); // Output: [{ a: 10 }, { b: 20 }]
 */
export const mergeReactiveChangedProperties = (oldData, newData) => {
  if (isArray(newData) && isArray(oldData)) {
    if (oldData.length !== newData.length) {
      return newData;
    }

    if (oldData.length === 0 && newData.length === 0) {
      return oldData;
    }

    // eslint-disable-next-line no-plusplus
    for (let index = 0; index < newData.length; index++) {
      // eslint-disable-next-line no-param-reassign
      oldData[index] = mergeReactiveChangedProperties(oldData[index], newData[index]);
    }

    return oldData;
  }

  if (!isObject(newData) || !isObject(oldData)) {
    return newData;
  }

  const oldKeys = Object.keys(oldData);
  const newKeys = Object.keys(newData);

  const removedKeys = difference(oldKeys, newKeys);

  for (const removedKey of removedKeys) {
    Vue.delete(oldData, removedKey);

    // eslint-disable-next-line no-param-reassign
    delete oldData[removedKey];
  }

  for (const key of newKeys) {
    const oldValue = oldData[key];
    const newValue = newData[key];

    Vue.set(oldData, key, mergeReactiveChangedProperties(oldValue, newValue));
  }

  return oldData;
};
