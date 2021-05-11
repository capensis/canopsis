import { zipObjectDeep } from 'lodash';
import flatten from 'flat';

/**
 * @typedef {Object.<string, string | string[]>} FlattenErrors
 */

/**
 * Map flatten errors object.
 *
 * @param {FlattenErrors} errors
 * @param {function} map
 * @return {FlattenErrors}
 */
export const flattenErrorMap = (errors = {}, map) => {
  const [errorsKeys, errorsValues] = Object.entries(errors);

  const errorsObject = zipObjectDeep(errorsKeys, errorsValues);

  return flatten(map(errorsObject));
};
