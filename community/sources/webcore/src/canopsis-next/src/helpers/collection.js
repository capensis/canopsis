import { isEqual, isMatch, pick, omit } from 'lodash';

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
 * Is equal objects with special keys
 *
 * @todo WORKS INCORRECTLY FOR ARRAYS
 *
 * @param {Object|Array} value
 * @param {Object|Array} other
 * @param {String|Array} [paths]
 * @return {boolean}
 */
export const isPickEqual = (value, other, paths) => isMatch(value, pick(other, paths));

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
 * Filters object entries by applying a filter function and returns a new object with truthy results
 *
 * @param {Object} obj - The object to filter
 * @param {Function} filterFn - Filter function that receives (value, key) and returns filtered value or null/undefined
 * @returns {Object} New object with only truthy filtered values
 */
export const filterObjectEntries = (obj, filterFn) => Object.entries(obj).reduce((acc, [key, value]) => {
  const result = filterFn(value, key);

  if (result) {
    acc[key] = result;
  }

  return acc;
}, {});

/**
 * Recursively collects all leaf node keys from a nested object structure
 *
 * @param {Object} obj - The object to collect keys from
 * @param {string} [childrenKey='children'] - The key name that contains child nodes
 * @returns {Set} Set of all leaf node keys
 */
export const collectLeafKeys = (obj, childrenKey = 'children') => {
  const keys = new Set();

  Object.entries(obj).forEach(([key, node]) => {
    if (node[childrenKey]) {
      collectLeafKeys(node[childrenKey], childrenKey).forEach(id => keys.add(id));
    } else {
      keys.add(key);
    }
  });

  return keys;
};
