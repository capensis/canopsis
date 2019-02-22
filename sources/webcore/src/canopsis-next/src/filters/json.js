import { isObject } from 'lodash';

/**
 * Convert JSON into JSON with indents
 *
 * @param {string|Object} json
 * @param {number} indents
 * @returns {string}
 */
export default function (json, indents = 4) {
  try {
    if (isObject(json)) {
      return JSON.stringify(json, null, indents);
    }

    return JSON.stringify(JSON.parse(json), null, indents);
  } catch (err) {
    console.warn(err);

    return '{}';
  }
}
