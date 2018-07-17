import getProp from 'lodash/get';

/**
 *
 * @param {Object} [object] - Object to search the property on
 * @param {String} [property] - Property name
 * @param {Function} [filter] - Filter to apply on the property
 *
 * @returns {String}
 */
export default function get(object, property, filter) {
  const value = getProp(object, property);
  if (filter) {
    return filter(value);
  }
  return value;
}
