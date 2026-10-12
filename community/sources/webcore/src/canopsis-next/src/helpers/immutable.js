import { get, clone, setWith, isFunction } from 'lodash';

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
