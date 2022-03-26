import { isEqual, omit } from 'lodash';

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
