import { isArray, isObject, difference } from 'lodash';

/**
 * Check whether the value is a DOM node.
 *
 * @param {*} value
 * @returns {boolean}
 */
const isNode = value => !!globalThis.Node && value instanceof globalThis.Node;

/**
 * Resolve the parent Vue component instance from Options API or Composition API setup.
 *
 * @param {Object} component Vue component instance.
 * @returns {Object|undefined} Parent component instance, if available.
 */
const getParentComponent = component => component?.$parent ?? component?.$.parent?.proxy;

/**
 * Safely check whether an object has an own property.
 *
 * @param {Object} object
 * @param {string} property
 * @returns {boolean}
 */
const hasOwnProperty = (object, property) => Object.prototype.hasOwnProperty.call(object, property);

/**
 * Check whether a child component is a wrapper whose `_original` reference points to the parent.
 *
 * @param {Object} child
 * @param {Object} parent
 * @returns {boolean}
 */
const isOriginalChild = (child, parent) =>
  hasOwnProperty(child, '_original') && child._original === parent;

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

  if (
    child === parent ||
    isOriginalChild(child, parent) ||
    (isNode(parentElement) && isNode(childElement) && parentElement.contains(childElement))
  ) {
    return true;
  }

  const childParent = getParentComponent(child);

  if (childParent) {
    return isParent(childParent, parent);
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
