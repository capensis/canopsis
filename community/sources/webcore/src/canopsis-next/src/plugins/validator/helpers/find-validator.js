/**
 * Finds validator in all parents
 */
export const findValidator = (node) => {
  const validator = node.$validator ?? node.context?.$validator;

  if (validator || !node.$parent) {
    return validator;
  }

  return findValidator(node.$parent);
};
