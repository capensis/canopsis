import { zipObjectDeep } from 'lodash';
import flatten from 'flat';

/**
 * @typedef FlattenErrors
 * @type {Object.<string, string | string[]>}
 */

/**
 * Map flatten errors object.
 *
 * @param {Object} errors
 * @param {function} map
 * @return {FlattenErrors}
 */
export const flattenErrorMap = (errors = {}, map) => {
  const [errorsKeys, errorsValues] = Object.entries(errors);

  const errorsObject = zipObjectDeep(errorsKeys, errorsValues);

  return flatten(map(errorsObject));
};
