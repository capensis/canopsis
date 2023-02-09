import { isObject, isString } from 'lodash';

/**
 * Convert JSON into JSON with indents
 *
 * @param {string|Object} json
 * @param {number} [indents = 4]
 * @param {string} [defaultValue = '{}']
 * @returns {string}
 */
export const stringifyJson = (json, indents = 4, defaultValue = '{}') => {
  if (!json) {
    return defaultValue;
  }

  if (isObject(json)) {
    return JSON.stringify(json, null, indents);
  }

  return JSON.stringify(JSON.parse(json), null, indents);
};

/**
 * Convert JSON into JSON with indents with error handling
 *
 * @param {string|Object} json
 * @param {number} [indents = 4]
 * @param {string} [defaultValue = '{}']
 * @returns {string}
 */
export const stringifyJsonFilter = (json, indents = 4, defaultValue = '{}') => {
  try {
    return stringifyJson(json, indents, defaultValue);
  } catch (err) {
    console.error(err);

    return defaultValue;
  }
};

/**
 * Json string validation check
 *
 * @param {string} json
 * @returns {boolean}
 */
export const isValidJsonData = (json) => {
  try {
    const data = JSON.parse(json);

    return !isString(data);
  } catch (err) {
    return false;
  }
};
