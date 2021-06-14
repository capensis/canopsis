/**
 * Finds the requested field by id from the context object.
 */
export function findField(el, context) {
  if (!context || !context.$validator) {
    return null;
  }

  // eslint-disable-next-line
  return context.$validator.fields.findById(el._veeValidateId);
}
