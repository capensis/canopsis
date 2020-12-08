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
