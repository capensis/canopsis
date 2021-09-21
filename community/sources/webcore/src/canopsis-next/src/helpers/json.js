import { isObject } from 'lodash';

/**
 * Convert JSON into JSON with indents
 *
 * @param {string|Object} json
 * @param {number} [indents = 4]
 * @param {string} [defaultValue = '{}']
 * @returns {string}
 */
export const stringifyJson = (json, indents = 4, defaultValue = '{}') => {
  try {
    if (!json) {
      return defaultValue;
    }

    if (isObject(json)) {
      return JSON.stringify(json, null, indents);
    }

    return JSON.stringify(JSON.parse(json), null, indents);
  } catch (err) {
    console.error(err);

    return defaultValue;
  }
};
