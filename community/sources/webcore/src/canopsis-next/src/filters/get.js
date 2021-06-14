import { get, isUndefined } from 'lodash';

/**
 *
 * @param {Object} object - Object to search the property on
 * @param {Array | string} property - Property name
 * @param {Function} [filter] - Filter to apply on the property
 * @param {any} [defaultValue] - Default value for the property
 *
 * @returns {string}
 */
export default function (object, property, filter, defaultValue) {
  let value = get(object, property);

  if (filter) {
    value = filter(value);
  }

  return !isUndefined(value) ? value : defaultValue;
}
