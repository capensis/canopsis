import { isEqual, isMatch, omit, pick } from 'lodash';

/**
 * Is equal objects without keys.
 *
 * @param {Object|Array} value
 * @param {Object|Array} other
 * @param {String|Array} [paths]
 * @returns {boolean}
 */
export const isOmitEqual = (value, other, paths) => isEqual(omit(value, paths), omit(other, paths));

/**
 * Is equal objects with special keys.
 *
 * @param {Object|Array} value
 * @param {Object|Array} other
 * @param {String|Array} [paths]
 * @returns {boolean}
 */
export const isPickEqual = (value, other, paths) => isMatch(value, pick(other, paths));

/**
 * Check is some fields in objects equal.
 *
 * @param {Object} value
 * @param {Object} other
 * @param {string[]} keys
 * @returns {boolean}
 */
export const isSeveralEqual = (value, other, keys = []) =>
  keys.every((key) => isEqual(value[key], other[key]));

/**
 * Revert grouped values by key.
 *
 * @param {Object.<string, string[]>} obj
 * @returns {Object.<string, string[]>}
 */
export const revertGroupBy = (obj) =>
  Object.entries(obj).reduce((acc, [id, values]) => {
    values.forEach((value) => {
      if (acc[value]) {
        acc[value].push(id);
      } else {
        acc[value] = [id];
      }
    });

    return acc;
  }, {});
