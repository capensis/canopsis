import getProp from 'lodash/get';

/**
 *
 * @param {Object} object - Object to search the property on
 * @param {String} property - Property name
 * @param {Function} [filter] - Filter to apply on the property
 * @param {Function} [defaultValue] - Default value for the property
 *
 * @returns {String}
 */
export default function get(object, property, filter, defaultValue) {
  let value = getProp(object, property);

  if (filter) {
    value = filter(value);
  }

  return value || defaultValue;
}
