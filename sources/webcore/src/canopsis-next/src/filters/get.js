import { get } from 'lodash';

/**
 *
 * @param {Object} object - Object to search the property on
 * @param {string} property - Property name
 * @param {Function} [filter] - Filter to apply on the property
 * @param {string} [defaultValue] - Default value for the property
 *
 * @returns {String}
 */
export default function (object, property, filter, defaultValue) {
  let value = get(object, property);

  if (filter) {
    value = filter(value);
  }

  return value || defaultValue;
}
