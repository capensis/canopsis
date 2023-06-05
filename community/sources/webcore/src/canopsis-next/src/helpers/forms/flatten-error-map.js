import { zipObjectDeep } from 'lodash';
import flatten from 'flat';

/**
 * @typedef {Object.<string, string | string[]>} FlattenErrors
 */

/**
 * Map flatten errors object.
 *
 * @param {FlattenErrors} [errors = {}]
 * @param {Function} [map = v => v]
 * @return {FlattenErrors}
 */
export const flattenErrorMap = (errors = {}, map = v => v) => {
  const errorsKeys = Object.keys(errors);
  const errorsValues = Object.values(errors);

  const errorsObject = zipObjectDeep(errorsKeys, errorsValues);

  return flatten(map(errorsObject));
};
