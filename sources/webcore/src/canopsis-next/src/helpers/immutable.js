import { get, clone, setWith, unset, isFunction } from 'lodash';

/**
 * Immutable method for deep updating object field or array item
 *
 * @param {Object|Array} obj - Object will be copied and copy will be updated
 * @param {string|Array} path - Path to field or item
 * @param {*} value - New field or item value or customizer function
 * @return {Object|Array}
 */
export function setIn(obj, path, value) {
  const preparedValue = isFunction(value) ? value(get(obj, path)) : value;

  return setWith(clone(obj), path, preparedValue, clone);
}

/**
 * Immutable method for deep updating object fields
 *
 * @param {Object|Array} obj - Object will be copied and copy will be updated
 * @param {Object} pathsValuesMap - Map for paths and values ex: { 'a.b.c': 'value', 'a.b.y': val => val }
 * @return {Object|Array}
 */
export function setInSeveral(obj, pathsValuesMap) {
  const alreadyClonedPaths = {};
  const clonedObject = clone(obj);

  Object.keys(pathsValuesMap).forEach((path) => {
    const preparedValue = isFunction(pathsValuesMap[path]) ?
      pathsValuesMap[path](get(obj, path)) : pathsValuesMap[path];
    let currentPath = '';

    setWith(clonedObject, path, preparedValue, (customizerValue, key) => {
      currentPath += `.${key}`;

      if (alreadyClonedPaths[currentPath]) {
        return customizerValue;
      }

      alreadyClonedPaths[currentPath] = true;

      return clone(customizerValue);
    });
  });

  return clonedObject;
}

/**
 * Immutable method for deep removing object field
 *
 * @param {Object|Array} obj - Object will be copied and copy will be updated
 * @param {string|Array} path - Path to field or item
 * @return {Object|Array}
 */
export function unsetIn(obj, path) {
  const newObj = setIn(obj, path, get(obj, path));

  unset(newObj, path);

  return newObj;
}

/**
 * Immutable method for deep adding new item into array
 *
 * @param {Object|Array} obj - Object will be copied and copy will be updated
 * @param {string|Array} path - Path to field or item
 * @param {*} value - Value of new array item
 * @return {Object|Array}
 */
export function addIn(obj, path, value) {
  return setIn(obj, path, [...get(obj, path, []), value]);
}

/**
 * Immutable method for deep removing item from array
 *
 * @param {Object|Array} obj - Object will be copied and copy will be updated
 * @param {string|Array} path - Path to field or item
 * @param {number} index - Index of array item
 * @return {Object|Array}
 */
export function removeIn(obj, path, index) {
  return setIn(obj, path, get(obj, path, []).filter((v, i) => i !== index));
}

export default {
  setIn,
  setInSeveral,
  unsetIn,
  addIn,
  removeIn,
};
