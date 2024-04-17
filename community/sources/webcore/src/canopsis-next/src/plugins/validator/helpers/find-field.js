import { findValidator } from './find-validator';

/**
 * Finds the requested field by id from the context object.
 */
export function findField(el, vnode) {
  const validator = findValidator(vnode);

  if (!vnode || !validator) {
    return null;
  }

  // eslint-disable-next-line
  return validator.fields.findById(el._veeValidateId);
}
