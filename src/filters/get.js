import get from 'lodash/get';

/**
 *
 * @param {Object} [object] - Object to search the property on
 * @param {String} [property] - Property name
 * @param {Function} [filter] - Filter to apply on the property
 * @param {string} [defaultValue] - Default value for the property
 *
 * @returns {String}
 */
export default function (object, property, filter, defaultValue) {
  const value = get(object, property);

  if (filter) {
    return filter(value);
  }

  return value || defaultValue;
}
