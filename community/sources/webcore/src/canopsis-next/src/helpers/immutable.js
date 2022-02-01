import { get, clone, setWith, unset, isFunction } from 'lodash';

/**
 * Immutable method for deep updating object field or array item
 *
 * @param {Object|Array} obj - Object will be copied and copy will be updated
 * @param {string|Array} path - Path to field or item
 * @param {*} value - New field or item value or customizer function
 * @return {Object|Array}
 */
export function setField(obj, path, value) {
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
export function setSeveralFields(obj, pathsValuesMap) {
  const alreadyClonedPaths = {};
  const clonedObject = clone(obj);

  Object.keys(pathsValuesMap).forEach((path) => {
    const func = pathsValuesMap[path];

    const preparedValue = isFunction(func)
      ? func(get(obj, path))
      : func;
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
export function unsetField(obj, path) {
  const newObj = setField(obj, path, get(obj, path));

  unset(newObj, path);

  return newObj;
}

/**
 * Immutable method for deep removing object fields by conditions
 *
 * @param {Object|Array} obj - Object will be copied and copy will be updated
 * @param {Object} pathsConditionsMap - Map for paths and conditions ex: { 'a.b.c': v => v < 5, 'a.y': v => !v.length }
 * @return {Object|Array}
 */
export function unsetSeveralFieldsWithConditions(obj, pathsConditionsMap) {
  const pathsValuesMap = Object.keys(pathsConditionsMap).reduce((acc, key) => {
    acc[key] = v => v;

    return acc;
  }, {});

  const newObj = setSeveralFields(obj, pathsValuesMap);

  Object.entries(pathsConditionsMap).forEach(([path, condition]) => {
    const value = get(obj, path);

    if (condition(value)) {
      unset(newObj, path);
    }
  });

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
export function addTo(obj, path, value) {
  return setField(obj, path, [...get(obj, path, []), value]);
}

/**
 * Immutable method for deep removing item from array
 *
 * @param {Object|Array} obj - Object will be copied and copy will be updated
 * @param {string|Array} path - Path to field or item
 * @param {number|Function} index - Index of array item or special condition function
 * @return {Object|Array}
 */
export function removeFrom(obj, path, index) {
  const filterFunc = isFunction(index) ? index : (v, i) => i !== index;

  return setField(obj, path, get(obj, path, []).filter(filterFunc));
}
