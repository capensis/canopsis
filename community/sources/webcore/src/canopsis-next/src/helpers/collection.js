import {
  difference,
  isArray,
  isEqual,
  isObject,
  omit,
} from 'lodash';

/**
 * Is equal objects without keys
 *
 * @param {Object|Array} value
 * @param {Object|Array} other
 * @param {String|Array} [paths]
 * @return {boolean}
 */
export const isOmitEqual = (value, other, paths) => isEqual(omit(value, paths), omit(other, paths));

/**
 * Check is some fields in objects equal
 *
 * @param {Object} value
 * @param {Object} other
 * @param {string[]} keys
 * @return {boolean}
 */
export const isSeveralEqual = (value, other, keys = []) => keys.every(key => isEqual(value[key], other[key]));

/**
 * Revert grouped values by key
 *
 * @example
 *  revertGroupBy({
 *    'key1': ['value1', 'value2'],
 *    'key2': ['value1', 'value2', 'value3'],
 *    'key3': ['value3'],
 *  }) -> {
 *    'value1': ['key1', 'key2'],
 *    'value2': ['key1', 'key2'],
 *    'value3': ['key2', 'key3'],
 *  }
 * @param {Object.<string, string[]>} obj
 * @returns {Object.<string, string[]>}
 */
export const revertGroupBy = obj => Object.entries(obj).reduce((acc, [id, values]) => {
  values.forEach((value) => {
    if (acc[value]) {
      acc[value].push(id);
    } else {
      acc[value] = [id];
    }
  });

  return acc;
}, {});

/**
 * Merge object without change links
 *
 * @param {any} oldData
 * @param {any} newData
 * @return {any}
 */
export const mergeChangedProperties = (oldData, newData) => {
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
      oldData[index] = mergeChangedProperties(oldData[index], newData[index]);
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
    // eslint-disable-next-line no-param-reassign
    delete oldData[removedKey];
  }

  for (const key of newKeys) {
    const oldValue = oldData[key];
    const newValue = newData[key];

    // eslint-disable-next-line no-param-reassign
    oldData[key] = mergeChangedProperties(oldValue, newValue);
  }

  return oldData;
};
