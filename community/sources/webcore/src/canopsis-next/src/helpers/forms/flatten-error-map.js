import { zipObjectDeep } from 'lodash';
import flatten from 'flat';

/**
 * Map flatten errors object.
 *
 * @param {Object} errors
 * @param {function} map
 * @return {Object}
 */
export const flattenErrorMap = (errors = {}, map) => {
  const errorsKeys = Object.keys(errors);
  const errorsValues = Object.values(errors);

  const errorsObject = zipObjectDeep(errorsKeys, errorsValues);

  return flatten(map(errorsObject));
};
