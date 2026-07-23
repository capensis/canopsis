import { isArray, isObject, difference } from 'lodash';

/**
 * Check if child component has parent component in the $parent chain.
 *
 * @param {Object|HTMLElement} child
 * @param {Object|HTMLElement} parent
 * @returns {boolean}
 */
export function isParent(child, parent) {
  if (!child || !parent) {
    return false;
  }

  const childElement = child.$el ?? child;
  const parentElement = parent.$el ?? parent;

  if (child === parent || child._original === parent || parentElement?.contains?.(childElement)) {
    return true;
  }

  if (child.$parent) {
    return isParent(child.$parent, parent);
  }

  return false;
}

/**
 * Merge changed properties into reactive object/array.
 *
 * In Vue 3 Vue.set/Vue.delete are not needed.
 *
 * @param {Object|Array} oldData
 * @param {Object|Array} newData
 * @returns {Object|Array}
 */
export const mergeReactiveChangedProperties = (oldData, newData) => {
  if (isArray(newData) && isArray(oldData)) {
    if (oldData.length !== newData.length) {
      return newData;
    }

    for (let index = 0; index < newData.length; index += 1) {
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

  difference(oldKeys, newKeys).forEach((key) => {
    // eslint-disable-next-line no-param-reassign
    delete oldData[key];
  });

  newKeys.forEach((key) => {
    // eslint-disable-next-line no-param-reassign
    oldData[key] = mergeReactiveChangedProperties(oldData[key], newData[key]);
  });

  return oldData;
};
